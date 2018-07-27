package storage

// ErrorHandler handles error
type ErrorHandler func(action string, err error) error

type Engine interface {
	Repository
	// Cache returns CacheProxy
	Cache() CacheProxy
	// Database returns DatabaseProxy
	Database() DatabaseProxy
	// SetErrorHandler sets handler for handling error
	SetErrorHandler(eh ErrorHandler)
	// AddIndex adds an index, should panic if repeated index name on a same table
	AddIndex(index Index)
	// NewSession new a session
	NewSession() Session
}

// engine implements Engine interface
type engine struct {
	name         string
	cache        CacheProxy
	database     DatabaseProxy
	errorHandler ErrorHandler
	indexes      map[string]map[string]Index
}

// NewEngine creates an engine which named name.
// Parameter database MUST be not nil, but cache can be nil.
func NewEngine(name string, database DatabaseProxy, cache CacheProxy) Engine {
	eng := &engine{
		name:     name,
		database: database,
		cache:    cache,
		indexes:  make(map[string]map[string]Index),
	}
	if eng.cache == nil {
		eng.cache = NullCacheProxy
	}
	return eng
}

func (eng *engine) Name() string            { return eng.name }
func (eng *engine) Cache() CacheProxy       { return eng.cache }
func (eng *engine) Database() DatabaseProxy { return eng.database }

// SetErrorHandler sets handler for handling error
func (eng *engine) SetErrorHandler(eh ErrorHandler) {
	eng.errorHandler = eh
}

// AddIndex adds an index
// panic if repeated index name on a same table
func (eng *engine) AddIndex(index Index) {
	tableName := index.TableMeta().Name()
	idx, ok := eng.indexes[tableName]
	if !ok {
		idx = make(map[string]Index)
		eng.indexes[tableName] = idx
	}
	indexName := index.Name()
	if _, exist := idx[indexName]; exist {
		panic("index " + indexName + " existed on table " + tableName)
	}
	idx[indexName] = index
}

func (eng *engine) NewSession() Session {
	return eng.newSession()
}

func (eng *engine) newSession() *session {
	s := &session{
		eng:      eng,
		database: eng.database.NewSession(),
	}
	if eng.cache != nil {
		s.cache = eng.cache.NewSession()
	}
	return s
}

// Insert inserts new records or updates all fields of records
func (eng *engine) Insert(tables ...Table) error {
	s := eng.newSession()
	defer s.Close()
	for _, table := range tables {
		action, err := s.update(table, true)
		if err != nil {
			return s.catch("Insert: "+action, err)
		}
	}
	return nil
}

// Update updates specific fields of record
func (eng *engine) Update(table Table, fields ...string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.update(table, false, fields...)
	if err != nil {
		return s.catch("Update: "+action, err)
	}
	return nil
}

// UpdateByIndexScore updates specific fields of records which satify index score range
func (eng *engine) UpdateByIndexScore(table Table, index Index, minScore, maxScore int64, fields ...string) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.updateByIndexScore(table, index, minScore, maxScore, fields...)
	if err != nil {
		action = "UpdateByIndexScore: table=" + table.TableMeta().Name() + ",index=" + index.Name() + ": " + action
		return nil, s.catch(action, err)
	}
	return result, nil
}

// Find gets records all fields by keys and stores loaded data to container
func (eng *engine) Find(keys KeyList, container TableListContainer, opts ...GetOption) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.find(keys, container, s.applyGetOption(opts), nil)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

// FindFields gets records specific fields by keys and stores loaded data to container
func (eng *engine) FindFields(keys KeyList, container TableListContainer, fields ...string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.find(keys, container, getOptions{}, fields)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

func (eng *engine) FindByIndexScore(index Index, minScore, maxScore int64, container TableListContainer, opts ...GetOption) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.findByIndexScore(index, minScore, maxScore, container, s.applyGetOption(opts), nil)
	if err != nil {
		return s.catch("FindByIndexScore: "+action, err)
	}
	return nil
}

func (eng *engine) FindFieldsByIndexScore(index Index, minScore, maxScore int64, container TableListContainer, fields ...string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.findByIndexScore(index, minScore, maxScore, container, getOptions{}, fields)
	if err != nil {
		return s.catch("FindFieldsByIndexScore: "+action, err)
	}
	return nil
}

// Get gets one record all fields
func (eng *engine) Get(table Table, opts ...GetOption) (bool, error) {
	s := eng.newSession()
	defer s.Close()
	action, ok, err := s.get(table, s.applyGetOption(opts), table.TableMeta().Fields()...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Get gets one record specific fields
func (eng *engine) GetFields(table Table, fields ...string) (bool, error) {
	s := eng.newSession()
	defer s.Close()
	action, ok, err := s.get(table, getOptions{}, fields...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Remove removes one record
func (eng *engine) Remove(table ReadonlyTable) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.remove(table.TableMeta(), table.Key())
	if err != nil {
		return s.catch("Remove: "+action, err)
	}
	return nil
}

// RemoveRecords removes records by keys
func (eng *engine) RemoveRecords(meta TableMeta, keys ...interface{}) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.remove(meta, keys...)
	if err != nil {
		return s.catch("RemoveRecords: "+action, err)
	}
	return nil
}

// Clear removes all records of table
func (eng *engine) Clear(tableName string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.clear(tableName)
	if err != nil {
		return s.catch("Clear "+tableName+": "+action, err)
	}
	return nil
}

// FindView loads view by keys and stores loaded data to container
func (eng *engine) FindView(view View, keys KeyList, container TableListContainer, opts ...GetOption) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.recursivelyLoadView(view, keys, container, s.applyGetOption(opts))
	if err != nil {
		return s.catch("FindView: "+action, err)
	}
	return nil
}

// IndexRank gets rank of table key in index, InvalidRank returned if key not found
func (eng *engine) IndexRank(index Index, key interface{}) (int64, error) {
	s := eng.newSession()
	defer s.Close()
	action, rank, err := s.indexRank(index, key)
	if err != nil {
		return rank, s.catch("IndexRank: "+action, err)
	}
	return rank, nil
}

// IndexScore gets score of table key in index, InvalidScore returned if key not found
func (eng *engine) IndexScore(index Index, key interface{}) (int64, error) {
	s := eng.newSession()
	defer s.Close()
	action, score, err := s.indexScore(index, key)
	if err != nil {
		return score, s.catch("IndexScore: "+action, err)
	}
	return score, nil
}

// IndexRange range index by rank
func (eng *engine) IndexRange(index Index, start, stop int64, opts ...RangeOption) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRange(index, start, stop, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRange: "+action, err)
	}
	return result, nil
}

// IndexRangeByScore range index by score
func (eng *engine) IndexRangeByScore(index Index, min, max int64, opts ...RangeOption) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRangeByScore(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByScore: "+action, err)
	}
	return result, nil
}

// IndexRangeByLex range index by lexicographical
func (eng *engine) IndexRangeByLex(index Index, min, max string, opts ...RangeOption) (RangeLexResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRangeByLex(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByLex: "+action, err)
	}
	return result, nil
}

// AddRecord adds a record with timestamp
func (eng *engine) AddRecord(key string, member interface{}, unixstamp int64) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.addRecord(key, member, unixstamp)
	if err != nil {
		return s.catch("AddRecord: "+action, err)
	}
	return nil
}

func (eng *engine) GetRecordsByTime(key string, startUnixstamp, endUnixstamp int64) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.getRecordsByTime(key, startUnixstamp, endUnixstamp)
	if err != nil {
		return nil, s.catch("GetRecordsByTime: "+action, err)
	}
	return result, nil
}

func (eng *engine) GetRecordsByPage(key string, pageSize int, startRank int64) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.getRecordsByPage(key, pageSize, startRank)
	if err != nil {
		return nil, s.catch("GetRecordsByPage: "+action, err)
	}
	return result, nil
}
