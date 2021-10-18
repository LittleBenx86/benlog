package uuid_test

import (
	"testing"

	AppUUID "github.com/LittleBenx86/Benlog/internal/utils/uuid"
)

func Test_SnowFlakeUUIDGeneration(t *testing.T) {
	t.Logf("%d\n", AppUUID.NewSnowFlake(0, 1, 0).GenerateId())
}
