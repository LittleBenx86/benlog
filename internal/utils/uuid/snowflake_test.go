package uuid

import (
	"testing"
)

func Test_SnowFlakeUUIDGeneration(t *testing.T) {
	t.Logf("%d\n", NewSnowFlake(0, 1, 0).GenerateId())
}
