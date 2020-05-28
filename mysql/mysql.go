// Package mysql client
package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // init mysql
)

// Client MySql transaction db
type Client struct {
	*TxDB
}

type MySql struct {
	addr               string
	user               string
	password           string
	dbName             string
	sslMode            bool
	maxOpenConnections int
	maxIdleConnections int
	charset            string
	connMaxLifetime    time.Duration
	parseTime          bool   // deal time.Time, set true
	loc                string // Local

	tablePrefix string
	debug       bool
}

// Option configures MySql using the functional options paradigm popularized by Rob Pike and Dave Cheney.
type Option interface {
	apply(do *MySql)
}

type optionFunc func(do *MySql)

func (fn optionFunc) apply(do *MySql) {
	fn(do)
}

// SSLMode ...
func SSLMode(b bool) Option {
	return optionFunc(func(do *MySql) {
		if b {
			do.sslMode = true
		}
	})
}

// ParseTime ...
func ParseTime(parseTime bool) Option {
	return optionFunc(func(do *MySql) {
		if parseTime {
			do.parseTime = parseTime
		}
	})
}

// MaxOpenConnections ...
func MaxOpenConnections(i int) Option {
	return optionFunc(func(do *MySql) {
		if i > 0 {
			do.maxOpenConnections = i
		}
	})
}

// MaxIdleConnections ...
func MaxIdleConnections(i int) Option {
	return optionFunc(func(do *MySql) {
		if i > 0 {
			do.maxIdleConnections = i
		}
	})
}

// Charset ...
func Charset(s string) Option {
	return optionFunc(func(do *MySql) {
		if s != "" {
			do.charset = s
		}
	})
}

// DialLoc ...
func Loc(s string) Option {
	return optionFunc(func(do *MySql) {
		if s != "" {
			do.loc = s
		}
	})
}

// ConnMaxLifetime ...
func ConnMaxLifetime(d time.Duration) Option {
	return optionFunc(func(do *MySql) {
		if d.Seconds() > 0 {
			do.connMaxLifetime = d
		}
	})
}

// TablePrefix specifies default table prefix
func TablePrefix(s string) Option {
	return optionFunc(func(do *MySql) {
		if s != "" {
			do.tablePrefix = s
		}
	})
}

// Debug specifies whether output the debug info or not
func Debug(d bool) Option {
	return optionFunc(func(do *MySql) {
		if d {
			do.debug = true
		}
	})
}

// Dial ...
func Dial(addr, user, password, dbName string, options ...Option) (c *Client, err error) {
	do := MySql{
		addr:     addr,
		user:     user,
		password: password,
		dbName:   dbName,
	}
	for _, option := range options {
		option.apply(&do)
	}

	urlBuf := bytes.NewBufferString(fmt.Sprintf("%s:%s@%s(%s)/%s", do.user, do.password, "tcp", do.addr, do.dbName))
	if do.parseTime {
		urlBuf.WriteString("?parseTime=True")
	}

	if do.charset != "" {
		urlBuf.WriteString(fmt.Sprintf("&charset=%s", do.charset))
	} else {
		urlBuf.WriteString("&charset=utf8mb4")
	}

	if do.loc != "" {
		loc := strings.Replace(do.loc, "/", "%2F", -1)
		urlBuf.WriteString(fmt.Sprintf("&loc=%s", loc))
	} else {
		urlBuf.WriteString(fmt.Sprintf("&loc=%s", "Local"))
	}

	db, err := sql.Open("mysql", urlBuf.String())
	if err != nil {
		panic(fmt.Sprintf("sql.Open()[%v] failed: %s", urlBuf.String(), err.Error()))
	}
	db.SetMaxIdleConns(do.maxIdleConnections)
	db.SetMaxOpenConns(do.maxOpenConnections)

	connMaxLifeTime := do.connMaxLifetime
	if connMaxLifeTime < 30 {
		connMaxLifeTime = 30
	}
	db.SetConnMaxLifetime(time.Duration(connMaxLifeTime) * time.Second)
	txDB := &TxDB{MDB: db}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c = &Client{txDB}
	return
}

// Close ...
func (d *Client) Close() error {
	if d == nil || d.TxDB == nil {
		return nil
	}
	return d.MDB.Close()
}

// Ping ...
func (d *Client) Ping() error {
	return d.MDB.Ping()
}
