package storage

import (
	"github.com/mkideal/pkg/typeconv"
	"gopkg.in/redis.v5"
)

type Tx interface {
	Begin() error
	Commit() error
	Rollback() error
	Close()
}

type Session interface {
	Tx
	Repository
	Cache() CacheProxySession
	Database() DatabaseProxySession
}

// session implements Session interface
type session struct {
	eng      *engine
	cache    CacheProxySession
	database DatabaseProxySession
}

func (s *session) catch(action string, err error) error {
	if err != nil && s.eng.errorHandler != nil {
		err = s.eng.errorHandler(action, err)
	}
	return err
}

func (s *session) Name() string                   { return s.eng.name }
func (s *session) Cache() CacheProxySession       { return s.cache }
func (s *session) Database() DatabaseProxySession { return s.database }

func (s *session) Begin() error {
	err := s.database.Begin()
	if err == nil && s.cache != nil {
		err = s.cache.Begin()
	}
	return err
}

func (s *session) Commit() error {
	err := s.database.Commit()
	if err == nil && s.cache != nil {
		err = s.cache.Commit()
	}
	return err
}

func (s *session) Rollback() error {
	err := s.database.Rollback()
	if err == nil && s.cache != nil {
		err = s.cache.Rollback()
	}
	return err
}

func (s *session) Close() {
	s.database.Close()
	if s.cache != nil {
		s.cache.Close()
	}
}

// Insert inserts new records or updates all fields of records
func (s *session) Insert(tables ...Table) error {
	for _, table := range tables {
		action, err := s.update(table, true)
		if err != nil {
			return s.catch("Insert: "+action, err)
		}
	}
	return nil
}

// Update updates specific fields of record
func (s *session) Update(table Table, fields ...string) error {
	action, err := s.update(table, false, fields...)
	if err != nil {
		return s.catch("Update: "+action, err)
	}
	return nil
}

// UpdateByIndexScore updates specific fields of records which satify index score range
func (s *session) UpdateByIndexScore(table Table, index Index, minScore, maxScore int64, fields ...string) (RangeResult, error) {
	action, result, err := s.updateByIndexScore(table, index, minScore, maxScore, fields...)
	if err != nil {
		action = "UpdateByIndexScore: table=" + table.TableMeta().Name() + ",index=" + index.Name() + ": " + action
		return nil, s.catch(action, err)
	}
	return result, nil
}

// FindFields gets records all fields by keys and stores loaded data to container
func (s *session) Find(keys KeyList, container TableListContainer, opts ...GetOption) error {
	action, err := s.find(keys, container, s.applyGetOption(opts), nil)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

// FindFields gets records specific fields by keys and stores loaded data to container
func (s *session) FindFields(keys KeyList, container TableListContainer, fields ...string) error {
	action, err := s.find(keys, container, getOptions{}, fields)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

func (s *session) FindByIndexScore(index Index, minScore, maxScore int64, container TableListContainer, opts ...GetOption) error {
	action, err := s.findByIndexScore(index, minScore, maxScore, container, s.applyGetOption(opts), nil)
	if err != nil {
		return s.catch("FindByIndexScore: "+action, err)
	}
	return nil
}

func (s *session) FindFieldsByIndexScore(index Index, minScore, maxScore int64, container TableListContainer, fields ...string) error {
	action, err := s.findByIndexScore(index, minScore, maxScore, container, getOptions{}, fields)
	if err != nil {
		return s.catch("FindFieldsByIndexScore: "+action, err)
	}
	return nil
}

// Get gets one record all fields
func (s *session) Get(table Table, opts ...GetOption) (bool, error) {
	action, ok, err := s.get(table, s.applyGetOption(opts), table.TableMeta().Fields()...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Get gets one record specific fields
func (s *session) GetFields(table Table, fields ...string) (bool, error) {
	action, ok, err := s.get(table, getOptions{}, fields...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Remove removes one record
func (s *session) Remove(table ReadonlyTable) error {
	action, err := s.remove(table.TableMeta(), table.Key())
	if err != nil {
		return s.catch("Remove: "+action, err)
	}
	return nil
}

// RemoveRecords removes records by keys
func (s *session) RemoveRecords(meta TableMeta, keys ...interface{}) error {
	action, err := s.remove(meta, keys...)
	if err != nil {
		return s.catch("RemoveRecords: "+action, err)
	}
	return nil
}

// Clear removes all records of table
func (s *session) Clear(table string) error {
	action, err := s.clear(table)
	if err != nil {
		return s.catch("Clear "+table+": "+action, err)
	}
	return nil
}

// FindView loads view by keys and stores loaded data to container
func (s *session) FindView(view View, keys KeyList, container TableListContainer, opts ...GetOption) error {
	action, err := s.recursivelyLoadView(view, keys, container, s.applyGetOption(opts))
	if err != nil {
		return s.catch("FindView: "+action, err)
	}
	return nil
}

// IndexRank gets rank of table key in index, InvalidRank returned if key not found
func (s *session) IndexRank(index Index, key interface{}) (int64, error) {
	action, rank, err := s.indexRank(index, key)
	if err != nil {
		return rank, s.catch("IndexRank: "+action, err)
	}
	return rank, nil
}

// IndexScore gets score of table key in index, InvalidScore returned if key not found
func (s *session) IndexScore(index Index, key interface{}) (int64, error) {
	action, score, err := s.indexScore(index, key)
	if err != nil {
		return score, s.catch("IndexScore: "+action, err)
	}
	return score, nil
}

// IndexRange range index by rank
func (s *session) IndexRange(index Index, start, stop int64, opts ...RangeOption) (RangeResult, error) {
	action, result, err := s.indexRange(index, start, stop, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRange: "+action, err)
	}
	return result, nil
}

// IndexRangeByScore range index by score
func (s *session) IndexRangeByScore(index Index, min, max int64, opts ...RangeOption) (RangeResult, error) {
	action, result, err := s.indexRangeByScore(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByScore: "+action, err)
	}
	return result, nil
}

// IndexRangeByLex range index by lexicographical order
func (s *session) IndexRangeByLex(index Index, min, max string, opts ...RangeOption) (RangeLexResult, error) {
	action, result, err := s.indexRangeByLex(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByLex: "+action, err)
	}
	return result, nil
}

// AddRecord adds a record with timestamp
func (s *session) AddRecord(key string, member interface{}, unixstamp int64) error {
	action, err := s.addRecord(key, member, unixstamp)
	if err != nil {
		return s.catch("AddRecord: "+action, err)
	}
	return nil
}

func (s *session) GetRecordsByTime(key string, startUnixstamp, endUnixstamp int64) (RangeResult, error) {
	action, result, err := s.getRecordsByTime(key, startUnixstamp, endUnixstamp)
	if err != nil {
		return nil, s.catch("GetRecordsByTime: "+action, err)
	}
	return result, nil
}

func (s *session) GetRecordsByPage(key string, pageSize int, startRank int64) (RangeResult, error) {
	action, result, err := s.getRecordsByPage(key, pageSize, startRank)
	if err != nil {
		return nil, s.catch("GetRecordsByPage: "+action, err)
	}
	return result, nil
}

//----------------
// implementation
//----------------

func (s *session) update(table Table, insert bool, fields ...string) (string, error) {
	// database op
	if insert {
		_, err := s.database.Insert(table)
		if err != nil {
			return action_db_insert, err
		}
	} else {
		if len(fields) == 0 {
			fields = table.TableMeta().Fields()
		}
		_, err := s.database.Update(table, fields...)
		if err != nil {
			return action_db_update, err
		}
	}

	// cache op
	if s.cache == nil {
		return action_null, nil
	}
	return s.updateCache(table, fields...)
}

func (s *session) updateByIndexScore(table Table, index Index, minScore, maxScore int64, fields ...string) (string, RangeResult, error) {
	action, result, err := s.indexRangeByScore(index, minScore, maxScore, rangeOptions{})
	if err != nil {
		return action, nil, err
	}
	originKey := table.Key()
	for i, n := 0, result.Len(); i < n; i++ {
		key := typeconv.ToString(result.Key(i))
		if err = table.SetKey(key); err != nil {
			return "set key " + key + table.TableMeta().Name(), nil, err
		}
		action, err = s.update(table, false, fields...)
		if err != nil {
			return "[key=" + key + "] " + action, nil, err
		}
	}
	table.SetKey(typeconv.ToString(originKey))
	return action_null, result, nil
}

func (s *session) updateCache(table Table, fields ...string) (string, error) {
	var (
		meta = table.TableMeta()
		key  = table.Key()
	)
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	args := make(map[string]string)
	fieldKey := typeconv.ToString(key)
	for _, field := range fields {
		key := JoinField(fieldKey, field)
		value, ok := table.GetField(field)
		if !ok {
			return action_get_field(meta.Name(), field), ErrFieldNotFound
		}
		args[key] = typeconv.ToString(value)
	}
	action, err := s.updateIndex(table, key, fields)
	if err != nil {
		return action, err
	}
	_, err = s.cache.HMSet(JoinKey(s.eng.name, meta.Name()), args)
	return action_cache_hmset, err
}

func (s *session) remove(meta TableMeta, keys ...interface{}) (string, error) {
	// database op
	if _, err := s.database.Remove(meta.Name(), meta.Key(), keys...); err != nil {
		return action_db_update, err
	}

	// cache op
	fields := meta.Fields()
	if s.cache == nil {
		return action_null, nil
	}
	args := make([]string, 0, len(fields)*len(keys))
	for _, key := range keys {
		fieldKey := typeconv.ToString(key)
		for _, field := range fields {
			args = append(args, JoinField(fieldKey, field))
		}
	}
	if action, err := s.removeIndex(meta.Name(), keys...); err != nil {
		return action, err
	}
	_, err := s.cache.HDel(JoinKey(s.eng.name, meta.Name()), args...)
	return action_cache_hdel, err
}

func (s *session) applyGetOption(opts []GetOption) getOptions {
	opt := getOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func (s *session) get(table Table, opt getOptions, fields ...string) (string, bool, error) {
	if s.cache == nil {
		return s.getFromDatabase(table, fields...)
	}
	action, found, err := s.getFromCache(table, fields...)
	if err != nil {
		return action, found, err
	}
	if !found && opt.syncFromDatabase {
		action, found, err = s.getFromDatabase(table, fields...)
		if err != nil {
			return action, found, err
		}
		if found {
			s.updateCache(table, fields...)
			// NOTE: this error sould be ignored
			return action_null, true, nil
		}
	}
	return action, found, err
}

func (s *session) getFromCache(table Table, fields ...string) (string, bool, error) {
	meta := table.TableMeta()
	tableKey := typeconv.ToString(table.Key())
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	fieldSize := len(fields)
	args := make([]string, 0, fieldSize)
	for _, field := range fields {
		args = append(args, JoinField(tableKey, field))
	}
	values, err := s.cache.HMGet(JoinKey(s.eng.name, meta.Name()), args...)
	if err != nil {
		return action_cache_hmget, false, err
	}
	if len(values) != fieldSize {
		return action_cache_hmget, false, ErrUnexpectedLength
	}
	found := false
	for i := 0; i < fieldSize; i++ {
		if values[i] != nil {
			if err := table.SetField(fields[i], typeconv.ToString(values[i])); err != nil {
				return action_set_field(meta.Name(), fields[i]), false, err
			}
			found = true
		}
	}
	return action_null, found, nil
}

func (s *session) getFromDatabase(table Table, fields ...string) (string, bool, error) {
	found, err := s.database.Get(table, fields...)
	if err != nil {
		return action_db_get, false, err
	}
	return action_null, found, nil
}

func (s *session) find(keys KeyList, container TableListContainer, opt getOptions, fields []string) (string, error) {
	meta := container.TableMeta()
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	_, action, err := s.findFields(keys, container, opt, FieldSlice(fields), nil)
	return action, err
}

func (s *session) findByIndexScore(index Index, minScore, maxScore int64, container TableListContainer, opt getOptions, fields []string) (string, error) {
	action, result, err := s.indexRangeByScore(index, minScore, maxScore, rangeOptions{})
	if err != nil {
		return action, err
	}
	return s.find(result, container, opt, fields)
}

func (s *session) clear(tableName string) (string, error) {
	if indexes, ok := s.eng.indexes[tableName]; ok {
		for _, index := range indexes {
			indexKey := JoinIndexKey(s.eng.name, index)
			if _, err := s.cache.Delete(indexKey); err != nil {
				return action_cache_del + ": delete index `" + indexKey + "`", err
			}
		}
	}
	key := JoinKey(s.eng.name, tableName)
	if _, err := s.cache.Delete(key); err != nil {
		return action_cache_del, err
	}
	return action_null, nil
}

func (s *session) findFields(keys KeyList, container TableListContainer, opt getOptions, fields FieldList, refs map[string]View) (map[string]StringKeys, string, error) {
	keySize := keys.Len()
	if keySize == 0 {
		return nil, action_null, nil
	}
	fieldSize := fields.Len()
	args := make([]string, 0, fieldSize*keySize)
	for i := 0; i < keySize; i++ {
		key := typeconv.ToString(keys.Key(i))
		for i := 0; i < fieldSize; i++ {
			args = append(args, JoinField(key, fields.Field(i)))
		}
	}
	tableName := container.TableMeta().Name()
	values, err := s.cache.HMGet(JoinKey(s.eng.name, tableName), args...)
	if err != nil {
		return nil, action_null, err
	}
	length := len(values)
	if length != fieldSize*keySize {
		return nil, action_null, ErrUnexpectedLength
	}
	var keysGroup map[string]StringKeys
	if len(refs) > 0 {
		keysGroup = make(map[string]StringKeys)
		for field := range refs {
			keysGroup[field] = StringKeys(make([]string, keySize))
		}
	}
	for i := 0; i+fieldSize <= length; i += fieldSize {
		index := i / fieldSize
		table, err := container.New(tableName, index, typeconv.ToString(keys.Key(index)))
		if err != nil {
			// NOTE: this error should be ignored
			continue
		}
		found := false
		for j := 0; j < fieldSize; j++ {
			field := fields.Field(j)
			value := values[i+j]
			if value != nil {
				found = true
				if err := table.SetField(field, typeconv.ToString(value)); err != nil {
					return nil, action_set_field(tableName, field), err
				}
			}
		}
		// get table data from database if not found in cache
		// and update to cache
		if !found && opt.syncFromDatabase {
			if action, _, err := s.getFromDatabase(table); err != nil {
				return keysGroup, action, err
			} else if action, err = s.updateCache(table); err != nil {
				return keysGroup, action, err
			}
		}
		// set keysGroup by table data
		for j := 0; j < fieldSize; j++ {
			field := fields.Field(j)
			if ks, ok := keysGroup[field]; ok {
				value, existedField := table.GetField(field)
				if !existedField {
					return keysGroup, action_get_field(table.TableMeta().Name(), field), ErrFieldNotFound
				} else {
					ks[index] = typeconv.ToString(value)
				}
				keysGroup[field] = ks
			}
		}
	}
	return keysGroup, action_null, nil
}

func (s *session) recursivelyLoadView(view View, keys KeyList, container TableListContainer, opt getOptions) (string, error) {
	keysGroup, action, err := s.findFields(keys, container, opt, view.Fields(), view.Refs())
	if err != nil {
		return action, err
	}
	refs := view.Refs()
	if refs == nil {
		return action_null, nil
	}
	if len(keysGroup) != len(refs) {
		return action, ErrUnexpectedLength
	}
	for field, ref := range refs {
		if tmpKeys, ok := keysGroup[field]; ok {
			if action, err := s.recursivelyLoadView(ref, tmpKeys, container, opt); err != nil {
				return action, err
			}
		} else {
			return action, ErrViewRefFieldMissing
		}
	}
	return action_null, nil
}

func (s *session) updateIndex(table ReadonlyTable, key interface{}, updatedFields []string) (action string, err error) {
	if indexes, ok := s.eng.indexes[table.TableMeta().Name()]; ok {
		for _, index := range indexes {
			if err = index.Update(s, table, key, updatedFields); err != nil {
				action = "update index `" + JoinIndexKey(s.eng.name, index) + "`"
				return
			}
		}
	}
	return
}

func (s *session) removeIndex(tableName string, keys ...interface{}) (action string, err error) {
	if indexes, ok := s.eng.indexes[tableName]; ok {
		for _, index := range indexes {
			if err = index.Remove(s, keys...); err != nil {
				action = "remove index `" + JoinIndexKey(s.eng.name, index) + "`"
				return
			}
		}
	}
	return
}

func (s *session) indexRank(index Index, key interface{}) (string, int64, error) {
	rank, err := s.cache.ZRank(JoinIndexKey(s.eng.name, index), typeconv.ToString(key))
	if err != nil {
		if err == ErrNotFound {
			return action_null, InvalidRank, nil
		}
		return action_cache_zrank, InvalidRank, err
	}
	if rank < 0 {
		rank = InvalidRank
	}
	return action_null, rank, nil
}

func (s *session) indexScore(index Index, key interface{}) (string, int64, error) {
	score, err := s.cache.ZScore(JoinIndexKey(s.eng.name, index), typeconv.ToString(key))
	if err != nil {
		if err == ErrNotFound {
			return action_null, InvalidScore, nil
		}
		return action_cache_zscore, score, err
	}
	return action_null, score, nil
}

func (s *session) applyRangeOption(opts []RangeOption) rangeOptions {
	opt := rangeOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func (s *session) indexRange(index Index, start, stop int64, opt rangeOptions) (action string, result RangeResult, err error) {
	key := JoinIndexKey(s.eng.name, index)
	if opt.rev {
		if opt.withScores {
			action = action_cache_zrevrangewithscores
			result, err = s.cache.ZRevRangeWithScores(key, start, stop)
		} else {
			action = action_cache_zrevrange
			result, err = s.cache.ZRevRange(key, start, stop)
		}
	} else {
		if opt.withScores {
			action = action_cache_zrangewithscores
			result, err = s.cache.ZRangeWithScores(key, start, stop)
		} else {
			action = action_cache_zrange
			result, err = s.cache.ZRange(key, start, stop)
		}
	}
	return
}

func (s *session) indexRangeByScore(index Index, min, max int64, opt rangeOptions) (action string, result RangeResult, err error) {
	key := JoinIndexKey(s.eng.name, index)
	byOpt := redis.ZRangeBy{
		Offset: opt.offset,
		Count:  opt.count,
	}
	if min == MinScore {
		byOpt.Min = "-inf"
	} else {
		byOpt.Min = typeconv.ToString(min)
	}
	if min == MaxScore {
		byOpt.Max = "+inf"
	} else {
		byOpt.Max = typeconv.ToString(max)
	}
	if opt.rev {
		if opt.withScores {
			action = action_cache_zrevrangebyscorewithscores
			result, err = s.cache.ZRevRangeByScoreWithScores(key, byOpt)
		} else {
			action = action_cache_zrevrangebyscore
			result, err = s.cache.ZRevRangeByScore(key, byOpt)
		}
	} else {
		if opt.withScores {
			action = action_cache_zrangebyscorewithscores
			result, err = s.cache.ZRangeByScoreWithScores(key, byOpt)
		} else {
			action = action_cache_zrangebyscore
			result, err = s.cache.ZRangeByScore(key, byOpt)
		}
	}
	return
}

func (s *session) indexRangeByLex(index Index, min, max string, opt rangeOptions) (action string, result RangeLexResult, err error) {
	key := JoinIndexKey(s.eng.name, index)
	byOpt := redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: opt.offset,
		Count:  opt.count,
	}
	if opt.rev {
		action = action_cache_zrevrangebylex
		result, err = s.cache.ZRevRangeByLex(key, byOpt)
	} else {
		action = action_cache_zrangebylex
		result, err = s.cache.ZRangeByLex(key, byOpt)
	}
	return
}

func (s *session) addRecord(key string, member interface{}, unixstamp int64) (action string, err error) {
	key = JoinKey(s.Name(), key)
	_, err = s.cache.ZAdd(key, redis.Z{Member: member, Score: float64(unixstamp)})
	if err != nil {
		return action_cache_zadd, err
	}
	return action_null, nil
}

func (s *session) getRecordsByTime(key string, startUnixstamp, endUnixstamp int64) (string, RangeResult, error) {
	key = JoinKey(s.Name(), key)
	opt := redis.ZRangeBy{Min: typeconv.ToString(startUnixstamp), Max: typeconv.ToString(endUnixstamp)}
	result, err := s.cache.ZRangeByScoreWithScores(key, opt)
	if err != nil {
		return action_cache_zrangebyscorewithscores, nil, err
	}
	return action_null, result, nil
}

func (s *session) getRecordsByPage(key string, pageSize int, startRank int64) (string, RangeResult, error) {
	key = JoinKey(s.Name(), key)
	result, err := s.cache.ZRangeWithScores(key, startRank, startRank+int64(pageSize-1))
	if err != nil {
		return action_cache_zrangewithscores, nil, err
	}
	return action_null, result, nil
}
