package consts

// Snow Flake
const (
	SnowFlakeStartTimeStamp     = int64(1631548800000) // timestamp millisecond, from 2021-09-14 00:00:00
	SnowFlakeMachineIdBits      = uint(10)
	SnowFlakeSequenceBits       = uint(12)
	SnowFlakeSequenceMask       = int64(-1 ^ (-1 << SnowFlakeSequenceBits))
	SnowFlakeMachineIdShiftLeft = SnowFlakeSequenceBits                          // machine id shift left bits
	SnowFlakeTimestampShiftLeft = SnowFlakeSequenceBits + SnowFlakeMachineIdBits // timestamp shift left bits
)

type AppStatusCode int

// Client error codes and error message
const ()

// Json Web Token
const (
	JWTOk          AppStatusCode = 200100
	JWTInvalid     AppStatusCode = -400100
	JWTExpired     AppStatusCode = -400101
	JWTFormatError AppStatusCode = -400102
)

// Request
const (
	RequestCommonSucceeded AppStatusCode = 200001

	RequestInvalidParameters AppStatusCode = -400200
	RequestInvalidAPIVersion AppStatusCode = -400201
)

// Security or authority
const (
	SecurityInvalidAccessAuthority AppStatusCode = -400400
	SecurityBlockedIPAccess        AppStatusCode = -400401
)

// App internal error
const (
	AppCommonInternalError AppStatusCode = 500001
)
