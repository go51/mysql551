package mysql551_test

import (
	"github.com/go51/mysql551"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	initialize()

	code := m.Run()

	os.Exit(code)

}

func initialize() {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)
	m.Open()
	defer m.Close()

	// truncate
	sql := "truncate table sample_table"
	_, _ = m.Exec(sql)
	sql = "truncate table sample_table_cud"
	_, _ = m.Exec(sql)

	// Insert
	sql = "insert into sample_table (name, description) values (?, ?)"
	_, _ = m.Exec(sql, "pubapp.biz_1", "domain_1")
	_, _ = m.Exec(sql, "pubapp.biz_2", "domain_2")
	_, _ = m.Exec(sql, "pubapp.biz_3", "domain_3")
	_, _ = m.Exec(sql, "pubapp.biz_4", "domain_4")
	_, _ = m.Exec(sql, "pubapp.biz_5", "domain_5")

}

func TestNew(t *testing.T) {
	config := &mysql551.Config{
		Host:     "localhost",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m1 := mysql551.New(config)
	m2 := mysql551.New(config)

	if m1 == nil {
		t.Errorf("インスタンスの生成に失敗しました。%#v\n", m1)
	}
	if m2 == nil {
		t.Errorf("インスタンスの生成に失敗しました。%#v\n", m2)
	}
	if &m1 == &m2 {
		t.Errorf("インスタンスの生成に失敗しました。\n [%p] != [%p]", &m1, &m2)
	}
}

func BenchmarkNew(b *testing.B) {
	config := &mysql551.Config{
		Host:     "localhost",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mysql551.New(config)
	}
}

func TestOpenClose(t *testing.T) {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	m.Open()
	if !m.IsOpen() {
		t.Errorf("データベースとの接続に失敗しました。\nResult: %v\n", m.IsOpen())
	}
	m.Close()
	if m.IsOpen() {
		t.Errorf("データベースとの切断に失敗しました。\nResult: %v\n", m.IsOpen())
	}

}

func BenchmarkOpenClose(b *testing.B) {
	b.SkipNow()
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Open()
		m.Close()
	}
}

type sample struct {
	id          int64
	name        string
	description string
}

func TestQueryAll(t *testing.T) {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	m.Open()
	defer m.Close()

	sqlAll := "select id, name, description from sample_table"
	rows := m.Query(sqlAll)
	defer rows.Close()

	i := 0
	for rows.Next() {
		s := sample{}
		err := rows.Scan(&s.id, &s.name, &s.description)
		if err != nil {
			t.Errorf("クエリ実行が失敗しました。\nResult: %d件", i)
		}
		i++
	}

	if i != 5 {
		t.Errorf("クエリ実行が失敗しました。\nResult: %d件", i)
	}

}

func TestQueryOne(t *testing.T) {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	m.Open()
	defer m.Close()

	sqlAll := "select id, name, description from sample_table where id = ?"
	rows := m.Query(sqlAll, 3)
	defer rows.Close()

	if rows.Next() {
		s := sample{}
		err := rows.Scan(&s.id, &s.name, &s.description)
		if err != nil {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.id != 3 {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.name != "pubapp.biz_3" {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.description != "domain_3" {
			t.Errorf("クエリ実行に失敗しました。")
		}
	}

}

func TestExecInsert(t *testing.T) {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	m.Open()
	defer m.Close()

	insert := "insert into sample_table_cud (name, description) values (?, ?)"
	affected, id := m.Exec(insert, "pubapp.biz", "test")

	if affected != 1 {
		t.Errorf("クエリ実行が失敗しました。\nAffected: %d", affected)
	}
	if id != 1 {
		t.Errorf("クエリ実行が失敗しました。\nId: %d", id)
	}
}

func TestExecUpdate(t *testing.T) {
	config := &mysql551.Config{
		Host:     "tcp(localhost:3306)",
		User:     "root",
		Password: "",
		Database: "mysql551_test",
	}

	m := mysql551.New(config)

	m.Open()
	defer m.Close()

	update := "update sample_table_cud set name = ?, description = ? where id = ?"
	affected, id := m.Exec(update, "pubapp.biz_1", "test_1", 1)

	if affected != 1 {
		t.Errorf("クエリ実行が失敗しました。\nAffected: %d", affected)
	}
	if id != 0 {
		t.Errorf("クエリ実行が失敗しました。\nId: %d", id)
	}

	update = "update sample_table_cud set name = ?, description = ? where id = ?"
	affected, id = m.Exec(update, "pubapp.biz_2", "test_2", 2)

	if affected != 0 {
		t.Errorf("クエリ実行が失敗しました。\nAffected: %d", affected)
	}
	if id != 0 {
		t.Errorf("クエリ実行が失敗しました。\nId: %d", id)
	}

	sql := "select id, name, description from sample_table_cud where id = ?"
	rows := m.Query(sql, 1)

	if rows.Next() {
		s := sample{}
		err := rows.Scan(&s.id, &s.name, &s.description)
		if err != nil {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.id != 1 {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.name != "pubapp.biz_1" {
			t.Errorf("クエリ実行に失敗しました。")
		}
		if s.description != "test_1" {
			t.Errorf("クエリ実行に失敗しました。")
		}
	}
}
