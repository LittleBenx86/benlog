package model

type Authority int

const (
	UNKNOWN_AUTHORITY Authority = 0
	ADMINISTRATOR     Authority = 1
	MEMBER            Authority = 2
)
