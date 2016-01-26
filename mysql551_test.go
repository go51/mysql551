package mysql551_test

import (
	"github.com/go51/mysql551"
	"testing"
)

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

func TestQuery(t *testing.T) {
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
	affected, id := m.Exec(insert, "okumura", "test")

	if affected != 1 {
		t.Errorf("クエリ実行が失敗しました。\nAffected: %d", affected)
	}
	if id < 1 {
		t.Errorf("クエリ実行が失敗しました。\nId: %d", id)
	}
}
