package uuid

import (
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/utils/uuid/intf"
	"sync"
	"time"
)

type snowflakeId struct {
	sync.Mutex
	timestamp int64
	machineId int64
	sequence  int64
}

func (s *snowflakeId) GenerateId() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1e6
	if s.timestamp != now {
		s.sequence = 0
	} else {
		s.sequence = (s.sequence + 1) & consts.SnowFlakeSequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	}

	s.timestamp = now
	r := (now-consts.SnowFlakeStartTimeStamp)<<consts.SnowFlakeTimestampShiftLeft |
		(s.machineId << consts.SnowFlakeMachineIdShiftLeft) | (s.sequence)
	return r
}

func NewSnowFlake(ts int64, mid int64, seq int64) intf.UUID {
	return &snowflakeId{
		timestamp: ts,
		machineId: mid,
		sequence:  seq,
	}
}
