package storage

import (
	"gopkg.in/redis.v5"
)

type CacheProxy interface {
	NewSession() CacheProxySession
}

type CacheProxySession interface {
	Tx
	HDel(table string, fields ...string) (int64, error)
	HExists(table, field string) (bool, error)
	HIncrBy(table, field string, incr int64) (int64, error)
	HMGet(table string, fields ...string) ([]interface{}, error)
	HMSet(table string, fields map[string]string) (string, error)
	Delete(keys ...string) (int64, error)
	ZAdd(key string, members ...redis.Z) (int64, error)
	ZRem(key string, members ...interface{}) (int64, error)
	ZRank(key, member string) (int64, error)
	ZScore(key, member string) (int64, error)
	// range APIs
	ZRange(key string, start, stop int64) (RangeResult, error)
	ZRangeWithScores(key string, start, stop int64) (RangeResult, error)
	ZRangeByScore(key string, opt redis.ZRangeBy) (RangeResult, error)
	ZRangeByLex(key string, opt redis.ZRangeBy) (RangeLexResult, error)
	ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) (RangeResult, error)
	ZRevRange(key string, start, stop int64) (RangeResult, error)
	ZRevRangeWithScores(key string, start, stop int64) (RangeResult, error)
	ZRevRangeByScore(key string, opt redis.ZRangeBy) (RangeResult, error)
	ZRevRangeByLex(key string, opt redis.ZRangeBy) (RangeLexResult, error)
	ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) (RangeResult, error)
}

type DatabaseProxy interface {
	NewSession() DatabaseProxySession
}

type DatabaseProxySession interface {
	Tx
	Insert(table Table) (int64, error)
	Update(table Table, fields ...string) (int64, error)
	Remove(tableName, keyName string, keys ...interface{}) (int64, error)
	Get(table Table, fields ...string) (bool, error)
}

// nullDatabaseProxy implements a null DatabaseProxy
type nullDatabaseProxy int

const NullDatabaseProxy = nullDatabaseProxy(0)

var _ = (DatabaseProxy)(NullDatabaseProxy)

func (x nullDatabaseProxy) NewSession() DatabaseProxySession                  { return x }
func (nullDatabaseProxy) Begin() error                                        { return nil }
func (nullDatabaseProxy) Commit() error                                       { return nil }
func (nullDatabaseProxy) Rollback() error                                     { return nil }
func (nullDatabaseProxy) Close()                                              { return }
func (nullDatabaseProxy) Insert(table Table) (int64, error)                   { return 1, nil }
func (nullDatabaseProxy) Update(table Table, fields ...string) (int64, error) { return 1, nil }
func (nullDatabaseProxy) Remove(tableName, keyName string, keys ...interface{}) (int64, error) {
	return int64(len(keys)), nil
}
func (nullDatabaseProxy) Get(table Table, fields ...string) (bool, error) { return false, nil }

// nullCacheProxy implements a null CacheProxy
type nullCacheProxy int

const NullCacheProxy = nullCacheProxy(0)

var _ = (CacheProxy)(NullCacheProxy)

func (x nullCacheProxy) NewSession() CacheProxySession                    { return x }
func (nullCacheProxy) Begin() error                                       { return nil }
func (nullCacheProxy) Commit() error                                      { return nil }
func (nullCacheProxy) Rollback() error                                    { return nil }
func (nullCacheProxy) Close()                                             { return }
func (nullCacheProxy) HDel(table string, fields ...string) (int64, error) { return 0, nil }
func (nullCacheProxy) HExists(table, field string) (bool, error)          { return false, nil }
func (nullCacheProxy) HIncrBy(table, field string, incr int64) (int64, error) {
	return 0, ErrNotImplemented
}
func (nullCacheProxy) HMGet(table string, fields ...string) ([]interface{}, error)  { return nil, nil }
func (nullCacheProxy) HMSet(table string, fields map[string]string) (string, error) { return "", nil }
func (nullCacheProxy) Delete(keys ...string) (int64, error)                         { return 0, nil }
func (nullCacheProxy) ZAdd(key string, members ...redis.Z) (int64, error)           { return 0, nil }
func (nullCacheProxy) ZRem(key string, members ...interface{}) (int64, error)       { return 0, nil }
func (nullCacheProxy) ZRank(key, member string) (int64, error)                      { return InvalidRank, nil }
func (nullCacheProxy) ZScore(key, member string) (int64, error)                     { return InvalidScore, nil }
func (nullCacheProxy) ZRange(key string, start, stop int64) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRangeWithScores(key string, start, stop int64) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRangeByScore(key string, opt redis.ZRangeBy) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRangeByLex(key string, opt redis.ZRangeBy) (RangeLexResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRevRange(key string, start, stop int64) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRevRangeWithScores(key string, start, stop int64) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRevRangeByScore(key string, opt redis.ZRangeBy) (RangeResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRevRangeByLex(key string, opt redis.ZRangeBy) (RangeLexResult, error) {
	return nil, ErrNotImplemented
}
func (nullCacheProxy) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) (RangeResult, error) {
	return nil, ErrNotImplemented
}
