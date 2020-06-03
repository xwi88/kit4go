// Package mysql client, only support MySQL 5.5+

// https://github.com/go-sql-driver/mysql#timetime-support
// https://github.com/go-sql-driver/mysql#unicode-support
package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // init and register mysql driver
	"github.com/kdpujie/log4go"
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
	charset            string
	collation          string // utf8mb4_general_ci
	maxOpenConnections int
	maxIdleConnections int
	debug              bool
	// default false, [true, false, skip-verify, preferred]
	// if true must set certificate
	tls             bool
	parseTime       bool // deal time.Time, set true
	connMaxLifetime time.Duration
	timeout         time.Duration // Timeout for establishing connections, aka dial timeout.
	readTimeout     time.Duration
	writeTimeout    time.Duration
	loc             string // default Local
	tablePrefix     string
}

// Option configures MySql using the functional options paradigm popularized by Rob Pike and Dave Cheney.
type Option interface {
	apply(do *MySql)
}

type optionFunc func(do *MySql)

func (fn optionFunc) apply(do *MySql) {
	fn(do)
}

// TLS tls=true enables TLS / SSL encrypted connection to the server. Use skip-verify if you want to use a
// self-signed or invalid certificate (server side) or use preferred to use TLS only when advertised by
// the server. This is similar to skip-verify, but additionally allows a fallback to a connection which is
// not encrypted. Neither skip-verify nor preferred add any reliable security. You can use a custom TLS
// config after registering it with mysql.RegisterTLSConfig.
func TLS(tls bool) Option {
	return optionFunc(func(do *MySql) {
		if tls {
			do.tls = tls
		}
	})
}

// MaxOpenConnections max open connections
// https://github.com/go-sql-driver/mysql#connection-pool-and-timeouts
func MaxOpenConnections(i int) Option {
	return optionFunc(func(do *MySql) {
		if i > 0 {
			do.maxOpenConnections = i
		}
	})
}

// MaxIdleConnections max idle connections
// https://github.com/go-sql-driver/mysql#connection-pool-and-timeouts
func MaxIdleConnections(i int) Option {
	return optionFunc(func(do *MySql) {
		if i > 0 {
			do.maxIdleConnections = i
		}
	})
}

// ConnMaxLifetime connection max lifetime
// https://github.com/go-sql-driver/mysql#connection-pool-and-timeouts
func ConnMaxLifetime(d time.Duration) Option {
	return optionFunc(func(do *MySql) {
		if d.Seconds() > 0 {
			do.connMaxLifetime = d
		}
	})
}

// Charset Sets the charset used for client-server interaction ("SET NAMES <value>").
// If multiple charsets are set (separated by a comma), the following charset is used if setting the charset failes.
// This enables for example support for utf8mb4 (introduced in MySQL 5.5.3) with fallback to utf8 for older servers
// (charset=utf8mb4,utf8).
// Usage of the charset parameter is discouraged because it issues additional queries to the server. Unless you need
// the fallback behavior, please use collation instead.
// https://github.com/go-sql-driver/mysql#charset
func Charset(charset string) Option {
	return optionFunc(func(do *MySql) {
		if charset != "" {
			do.charset = charset
		}
	})
}

// Collation Sets the collation used for client-server interaction on connection. In contrast to charset, collation does
// not issue additional queries. If the specified collation is unavailable on the target server, the connection will fail.
// A list of valid charsets for a server is retrievable with SHOW COLLATION.
// The default collation (utf8mb4_general_ci) is supported from MySQL 5.5. You should use an older collation
// (e.g. utf8_general_ci) for older MySQL.
// https://github.com/go-sql-driver/mysql#collation
func Collation(collation string) Option {
	return optionFunc(func(do *MySql) {
		if collation != "" {
			do.collation = collation
		}
	})
}

// Loc Sets the location for time.Time values (when using parseTime=true). "Local" sets the system's location.
// https://github.com/go-sql-driver/mysql#loc
func Loc(s string) Option {
	return optionFunc(func(do *MySql) {
		if s != "" {
			do.loc = s
		}
	})
}

// ParseTime parseTime=true changes the output type of DATE and DATETIME values to time.Time instead of
// []byte / string The date or datetime like 0000-00-00 00:00:00 is converted into zero value of time.Time.
// https://github.com/go-sql-driver/mysql#parsetime
func ParseTime(parseTime bool) Option {
	return optionFunc(func(do *MySql) {
		if parseTime {
			do.parseTime = parseTime
		}
	})
}

// Timeout for establishing connections, aka dial timeout. The value must be a decimal number with a unit
// suffix ("ms", "s", "m", "h"), such as "30s", "0.5m" or "1m30s".
func Timeout(timeout time.Duration) Option {
	return optionFunc(func(do *MySql) {
		if timeout >= time.Millisecond {
			do.timeout = timeout
		}
	})
}

// ReadTimeout I/O read timeout. The value must be a decimal number with a unit
// suffix ("ms", "s", "m", "h"), such as "30s", "0.5m" or "1m30s".
func ReadTimeout(readTimeout time.Duration) Option {
	return optionFunc(func(do *MySql) {
		if readTimeout >= time.Millisecond {
			do.readTimeout = readTimeout
		}
	})
}

// WriteTimeout I/O write timeout. The value must be a decimal number with a unit suffix ("ms", "s", "m", "h"),
// such as "30s", "0.5m" or "1m30s".
// Default 0
func WriteTimeout(writeTimeout time.Duration) Option {
	return optionFunc(func(do *MySql) {
		if writeTimeout >= time.Millisecond {
			do.writeTimeout = writeTimeout
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

// Dial dial mysql
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
	// Notes: Watch out here, must start ?
	// 为了处理time.Time，您需要包括parseTime作为参数
	if do.parseTime {
		urlBuf.WriteString("?parseTime=True")
	} else {
		urlBuf.WriteString("?parseTime=False")
	}
	urlBuf.WriteString(fmt.Sprintf("&readTimeout=%v", do.readTimeout))
	urlBuf.WriteString(fmt.Sprintf("&writeTimeout=%v", do.writeTimeout))
	urlBuf.WriteString(fmt.Sprintf("&timeout=%v", do.timeout))
	// default utf8mb4_general_ci
	if do.collation != "" {
		urlBuf.WriteString(fmt.Sprintf("&collation=%s", do.collation))
	}
	// not recommend set, replace with collation
	if do.charset != "" {
		urlBuf.WriteString(fmt.Sprintf("&charset=%s", do.charset))
	}

	if do.loc != "" {
		loc := strings.Replace(do.loc, "/", "%2F", -1)
		urlBuf.WriteString(fmt.Sprintf("&loc=%s", loc))
	} else {
		urlBuf.WriteString(fmt.Sprintf("&loc=%s", "Local"))
	}

	if do.tls {
		urlBuf.WriteString("&tls=true")
	}

	if do.debug {
		log4go.Debug("[mysql] dataSourceName:%v", urlBuf.String())
	}
	db, err := sql.Open("mysql", urlBuf.String())
	if err != nil {
		log4go.Error("[mysql] Open()[%v] failed: %s", urlBuf.String(), err.Error())
		return nil, err
	}
	db.SetMaxIdleConns(do.maxIdleConnections)
	db.SetMaxOpenConns(do.maxOpenConnections)

	connMaxLifeTime := do.connMaxLifetime
	// less than 1 second, set 30 second
	if connMaxLifeTime <= time.Second {
		connMaxLifeTime = time.Second * 30
	}
	db.SetConnMaxLifetime(connMaxLifeTime)
	txDB := &TxDB{MDB: db}

	if do.debug {
		log4go.Debug("[mysql] db config:%#v", do)
	}
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
