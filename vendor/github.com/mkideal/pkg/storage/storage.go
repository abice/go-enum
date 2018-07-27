package storage

import (
	"errors"
	"math"
)

const (
	Unused = 0

	ErrorHandlerDepth = 4

	InvalidRank  = -1
	InvalidScore = math.MinInt64

	MinScore = math.MinInt64
	MaxScore = math.MaxInt64
)

var (
	ErrNotFound            = errors.New("not found")
	ErrFieldNotFound       = errors.New("field not found")
	ErrUnexpectedLength    = errors.New("unexpected length")
	ErrViewRefFieldMissing = errors.New("view ref field missing")
	ErrTableNotFoundInView = errors.New("table not found in view")
	ErrTypeAssert          = errors.New("type assert failed")
	ErrNotImplemented      = errors.New("not implemented")
)

const (
	action_null                             = ""
	action_cache_hmget                      = "cache.HMGet"
	action_cache_hmset                      = "cache.HMSet"
	action_cache_hdel                       = "cache.HDel"
	action_cache_del                        = "cache.Del"
	action_cache_zadd                       = "cache.ZAdd"
	action_cache_zscore                     = "cache.ZScore"
	action_cache_zrank                      = "cache.ZRank"
	action_cache_zrange                     = "cache.ZRange"
	action_cache_zrevrange                  = "cache.ZRevRange"
	action_cache_zrangebyscore              = "cache.ZRangeByScore"
	action_cache_zrevrangebyscore           = "cache.ZRevRangeByScore"
	action_cache_zrangebylex                = "cache.ZRangeByLex"
	action_cache_zrevrangebylex             = "cache.ZRevRangeByLex"
	action_cache_zrangewithscores           = "cache.ZRangeWithScores"
	action_cache_zrevrangewithscores        = "cache.ZRevRangeWithScores"
	action_cache_zrangebyscorewithscores    = "cache.ZRangeByScoreWithScores"
	action_cache_zrevrangebyscorewithscores = "cache.ZRevRangeByScoreWithScores"
	action_db_insert                        = "db.Insert"
	action_db_update                        = "db.Update"
	action_db_remove                        = "db.Remove"
	action_db_get                           = "db.Get"
)

func action_get_field(table, field string) string {
	return "table `" + table + "` GetField `" + field + "`"
}
func action_set_field(table, field string) string {
	return "table `" + table + "` SetField `" + field + "`"
}

func JoinKey(engineName, originKey string) string {
	return engineName + "@" + originKey
}

func JoinField(key, originField string) string {
	return key + ":" + originField
}

func JoinIndexKey(engineName string, index Index) string {
	return engineName + "@" + index.TableMeta().Name() + ":" + index.Name()
}

// GetOption represents options for Get/Find/FindView operations
type GetOption func(*getOptions)

type getOptions struct {
	syncFromDatabase bool
}

// WithSyncFromDatabase returns a GetOption which would get data from database if data not found in cache
func WithSyncFromDatabase() GetOption {
	return syncFromDatabase
}

func syncFromDatabase(opt *getOptions) {
	opt.syncFromDatabase = true
}

// RangeOption represents options for IndexRange operations
type RangeOption func(*rangeOptions)

type rangeOptions struct {
	withScores bool
	rev        bool
	offset     int64
	count      int64
}

func RangeRev() RangeOption {
	return func(opts *rangeOptions) {
		opts.rev = true
	}
}

func RangeWithScores() RangeOption {
	return func(opts *rangeOptions) {
		opts.withScores = true
	}
}

func RangeOffset(offset int64) RangeOption {
	return func(opts *rangeOptions) {
		opts.offset = offset
	}
}

func RangeCount(count int64) RangeOption {
	return func(opts *rangeOptions) {
		opts.count = count
	}
}

// RangeLexResult represents result of range by lex
type RangeLexResult interface {
	KeyList
}

// RangeLexResult represents result of range by rank,score etc.
type RangeResult interface {
	KeyList
	Score(i int) int64
}

func ContainsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}
