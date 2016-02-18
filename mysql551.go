package mysql551

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Mysql struct {
	config *Config
	db     *sql.DB
	tx     *sql.Tx

	open bool
	tran bool
}

type Config struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func New(config *Config) *Mysql {
	m := Mysql{
		config: config,
		open:   false,
	}

	return &m
}

func (m *Mysql) Open() {
	dataSourceName := m.config.User + ":" + m.config.Password + "@" + m.config.Host + "/" + m.config.Database

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	m.db = db
	m.open = true
}

func (m *Mysql) Close() {
	err := m.db.Close()
	if err != nil {
		panic(err)
	}

	m.db = nil
	m.open = false
}

func (m *Mysql) IsOpen() bool {
	if !m.open {
		return false
	}

	if m.db == nil {
		return false
	}

	return true
}

func (m *Mysql) Begin() {
	if !m.IsOpen() {
		panic(errors.New("[mysql555]It is not connected to the host"))
	}

	if m.IsTransaction() {
		return
	}

	tx, err := m.db.Begin()
	if err != nil {
		panic(err)
	}

	m.tx = tx
	m.tran = true
}

func (m *Mysql) Commit() {
	if !m.IsOpen() {
		panic(errors.New("[mysql555]It is not connected to the host"))
	}

	if !m.IsTransaction() {
		return
	}

	err := m.tx.Commit()
	if err != nil {
		panic(err)
	}

	m.tx = nil
	m.tran = false

}

func (m *Mysql) Rollback() {
	if !m.IsOpen() {
		panic(errors.New("[mysql555]It is not connected to the host"))
	}

	if !m.IsTransaction() {
		return
	}

	err := m.tx.Rollback()
	if err != nil {
		panic(err)
	}

	m.tx = nil
	m.tran = false

}

func (m *Mysql) IsTransaction() bool {
	if !m.tran {
		return false
	}

	if m.tx == nil {
		return false
	}

	return true
}

func (m *Mysql) Query(query string, param ...interface{}) *sql.Rows {
	if !m.IsOpen() {
		panic(errors.New("[mysql555]It is not connected to the host"))
	}

	var rows *sql.Rows = nil
	var err error = nil

	if m.IsTransaction() {
		rows, err = m.tx.Query(query, param...)
	} else {
		rows, err = m.db.Query(query, param...)
	}

	if err != nil {
		message := err.Error()
		message = message + "\n" + "[ SQL ] " + query + "\n" + fmt.Sprintf("[ PRM ] %#v\n", param)
		panic(errors.New(message))
	}

	return rows

}

func (m *Mysql) Exec(query string, param ...interface{}) (rowsAffected int64, lastInsertId int64) {
	if !m.IsOpen() {
		panic(errors.New("[mysql555]It is not connected to the host"))
	}

	var res sql.Result = nil
	var err error = nil

	if m.IsTransaction() {
		res, err = m.tx.Exec(query, param...)
	} else {
		res, err = m.db.Exec(query, param...)
	}

	if err != nil {
		message := err.Error()
		message = message + "\n" + "[ SQL ] " + query + "\n" + fmt.Sprintf("[ PRM ] %#v\n", param)
		panic(errors.New(message))
	}

	lastInsertId, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}

	return
}

var nilTime time.Time = time.Time{}
var nilString string = ""
var nilInt64 int64 = 0

func ToNil(value interface{}) interface{} {
	if val, ok := value.(string); ok {
		// String
		if val == nilString {
			return nil
		} else {
			return value
		}
	} else if val, ok := value.(int64); ok {
		// Int64
		if val == nilInt64 {
			return nil
		} else {
			return value
		}
	} else if val, ok := value.(time.Time); ok {
		// time.Time
		if val == nilTime {
			return nil
		} else {
			return value
		}
	} else {
		return value
	}
}

type GString string
type GTime time.Time
type GInt64 int64
