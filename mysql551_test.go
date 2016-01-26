package mysql551_test
import (
	"testing"
	"github.com/go51/mysql551"
)

func TestNew(t *testing.T) {
	config := &mysql551.Config{
		Host:"localhost",
		User:"root",
		Password:"",
		Database:"mysql551_test",
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
		Host:"localhost",
		User:"root",
		Password:"",
		Database:"mysql551_test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mysql551.New(config)
	}
}