package uuid

import "testing"

func Test_GenerateRandomStringId(t *testing.T) {
	id, err := GenerateRandomStringId(31)
	if err != nil {
		t.Fatal("random string id error", err)
	}
	t.Log(id)
}

func Benchmark_GenerateRandomStringId(b *testing.B) {
	for i := 0; i < 200; i++ {
		id, err := GenerateRandomStringId(31)
		if err != nil {
			b.Fatal("random string id error", err)
		}
		b.Log(id)
	}
}
