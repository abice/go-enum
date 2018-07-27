package goredisproxy

import (
	"gopkg.in/redis.v5"

	"github.com/mkideal/pkg/storage"
)

// proxy implements storage.CacheProxy
type proxy struct {
	client redis.Cmdable
}

func New(client redis.Cmdable) storage.CacheProxy {
	return &proxy{client: client}
}

func (p *proxy) NewSession() storage.CacheProxySession {
	return &proxySession{client: p.client}
}

// proxySession implements storage.CacheProxySession
type proxySession struct {
	client redis.Cmdable
}

func (p *proxySession) Begin() error    { return nil }
func (p *proxySession) Commit() error   { return nil }
func (p *proxySession) Rollback() error { return nil }
func (p *proxySession) Close()          {}

func (p *proxySession) error(err error) error {
	if err == redis.Nil {
		return storage.ErrNotFound
	}
	return err
}

func (p *proxySession) HDel(table string, fields ...string) (int64, error) {
	x, err := p.client.HDel(table, fields...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HExists(table, field string) (bool, error) {
	x, err := p.client.HExists(table, field).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HIncrBy(table, field string, incr int64) (int64, error) {
	x, err := p.client.HIncrBy(table, field, incr).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HMGet(table string, fields ...string) ([]interface{}, error) {
	x, err := p.client.HMGet(table, fields...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HMSet(table string, fields map[string]string) (string, error) {
	x, err := p.client.HMSet(table, fields).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) Delete(keys ...string) (int64, error) {
	x, err := p.client.Del(keys...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZAdd(key string, members ...redis.Z) (int64, error) {
	x, err := p.client.ZAdd(key, members...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZRem(key string, members ...interface{}) (int64, error) {
	x, err := p.client.ZRem(key, members...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZRank(key, member string) (int64, error) {
	x, err := p.client.ZRank(key, member).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZScore(key, member string) (int64, error) {
	x, err := p.client.ZScore(key, member).Result()
	err = p.error(err)
	return int64(x), err
}

type stringSliceResult []string

func (ss stringSliceResult) Len() int              { return len(ss) }
func (ss stringSliceResult) Key(i int) interface{} { return ss[i] }
func (ss stringSliceResult) Score(i int) int64     { return 0 }

type zsliceResult []redis.Z

func (zs zsliceResult) Len() int              { return len(zs) }
func (zs zsliceResult) Key(i int) interface{} { return zs[i].Member }
func (zs zsliceResult) Score(i int) int64     { return int64(zs[i].Score) }

func (p *proxySession) ZRange(key string, start, stop int64) (storage.RangeResult, error) {
	res := p.client.ZRange(key, start, stop)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRangeWithScores(key string, start, stop int64) (storage.RangeResult, error) {
	res := p.client.ZRangeWithScores(key, start, stop)
	return zsliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRangeByScore(key string, opt redis.ZRangeBy) (storage.RangeResult, error) {
	res := p.client.ZRangeByScore(key, opt)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRangeByLex(key string, opt redis.ZRangeBy) (storage.RangeLexResult, error) {
	res := p.client.ZRangeByLex(key, opt)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) (storage.RangeResult, error) {
	res := p.client.ZRangeByScoreWithScores(key, opt)
	return zsliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRevRange(key string, start, stop int64) (storage.RangeResult, error) {
	res := p.client.ZRevRange(key, start, stop)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRevRangeWithScores(key string, start, stop int64) (storage.RangeResult, error) {
	res := p.client.ZRevRangeWithScores(key, start, stop)
	return zsliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRevRangeByScore(key string, opt redis.ZRangeBy) (storage.RangeResult, error) {
	res := p.client.ZRevRangeByScore(key, opt)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRevRangeByLex(key string, opt redis.ZRangeBy) (storage.RangeLexResult, error) {
	res := p.client.ZRevRangeByLex(key, opt)
	return stringSliceResult(res.Val()), p.error(res.Err())
}

func (p *proxySession) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) (storage.RangeResult, error) {
	res := p.client.ZRevRangeByScoreWithScores(key, opt)
	return zsliceResult(res.Val()), p.error(res.Err())
}
