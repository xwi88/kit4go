package mysql

import (
	"log"
	"testing"
	"time"
)

func Test_Dial(t *testing.T) {

	connMaxLifetime, _ := time.ParseDuration("30m")
	timeout, _ := time.ParseDuration("5s")
	readTimeout, _ := time.ParseDuration("1s")
	writeTimeout, _ := time.ParseDuration("1s")

	options := MySql{
		addr:     "127.0.0.1:3306",
		user:     "root",
		password: "root1234",
		dbName:   "test",
		// charset:            "utf8mb4",
		collation:          "utf8mb4_general_ci",
		maxOpenConnections: 64,
		maxIdleConnections: 32,
		debug:              true,
		parseTime:          true,
		loc:                "",
		tls:                false, // if true must set certificate
		connMaxLifetime:    connMaxLifetime,
		timeout:            timeout,
		readTimeout:        readTimeout,
		writeTimeout:       writeTimeout,
	}
	client, err := initMySql(options)
	if err != nil {
		log.Fatalf("initMySql err:%v", err.Error())
		return
	}
	defer func() {
		_ = client.Close()
	}()

	err = client.MDB.Ping()
	if err != nil {
		log.Panicf("Ping err:%v", err.Error())
		return
	}
}

// initMySql create mysql client instance
func initMySql(options MySql) (client *Client, err error) {
	charset := Charset(options.charset)
	collation := Collation(options.collation)
	maxOpenConnections := MaxOpenConnections(options.maxOpenConnections)
	maxIdleConnections := MaxIdleConnections(options.maxIdleConnections)
	debug := Debug(options.debug)
	parseTime := ParseTime(options.parseTime)
	loc := Loc(options.loc)
	tls := TLS(options.tls)
	connMaxLifetime := ConnMaxLifetime(options.connMaxLifetime)
	timeout := Timeout(options.timeout)
	readTimeout := ReadTimeout(options.readTimeout)
	writeTimeout := WriteTimeout(options.writeTimeout)

	client, err = Dial(options.addr, options.user, options.password, options.dbName,
		charset, collation, maxOpenConnections, maxIdleConnections,
		debug, parseTime, loc, tls, connMaxLifetime,
		timeout, readTimeout, writeTimeout)
	if err != nil {
		return nil, err
	}
	return client, err
}
