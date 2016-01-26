package mysql551

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	config *Config
	db     *sql.DB

	open bool
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

func (m *Mysql) Query(query string, param ...interface{}) *sql.Rows {
	rows, err := m.db.Query(query, param...)
	if err != nil {
		panic(err)
	}

	return rows

}

func (m *Mysql) Exec(query string, param ...interface{}) (rowsAffected int64, lastInsertId int64) {
	res, err := m.db.Exec(query, param...)
	if err != nil {
		panic(err)
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
