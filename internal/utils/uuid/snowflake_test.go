package uuid

import (
	"testing"
)

func Test_SnowFlakeUUIDGeneration(t *testing.T) {
	id := NewSnowFlake(0, 1, 0).GenerateId()
	t.Logf("%d\n", id)
	var idCpy uint
	idCpy = uint(id)
	t.Logf("%d\n", idCpy)
}
