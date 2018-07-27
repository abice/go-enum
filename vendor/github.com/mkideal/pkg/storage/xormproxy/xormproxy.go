package xormproxy

import (
	"bytes"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/mkideal/pkg/storage"
)

// proxy implements storage.DatabaseProxy interface
type proxy struct {
	engine *xorm.Engine
}

func New(engine *xorm.Engine) storage.DatabaseProxy {
	return &proxy{engine: engine}
}

func (p *proxy) NewSession() storage.DatabaseProxySession {
	return &proxySession{session: p.engine.NewSession()}
}

// proxySession implements storage.DatabaseProxySession interface
type proxySession struct {
	session *xorm.Session
}

func (p *proxySession) Begin() error    { return p.session.Begin() }
func (p *proxySession) Commit() error   { return p.session.Commit() }
func (p *proxySession) Rollback() error { return p.session.Rollback() }
func (p *proxySession) Close()          { p.session.Close() }

func (p *proxySession) Insert(table storage.Table) (int64, error) {
	return p.session.InsertOne(table)
}

func (p *proxySession) Update(table storage.Table, fields ...string) (int64, error) {
	return p.session.ID(table.Key()).Cols(fields...).Update(table)
}

func (p *proxySession) Remove(tableName, keyName string, keys ...interface{}) (int64, error) {
	n := len(keys)
	if n == 0 {
		return 0, nil
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "delete from %s where ", tableName)
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprintf(&buf, " or ")
		}
		fmt.Fprintf(&buf, "%s = ?", keyName)
	}
	result, err := p.session.Exec(buf.String(), keys...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (p *proxySession) Get(table storage.Table, fields ...string) (bool, error) {
	return p.session.Cols(fields...).Get(table)
}
