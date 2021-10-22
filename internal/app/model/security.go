package model

type Authority int

const (
	UNKNOWN_AUTHORITY Authority = 0
	ADMINISTRATOR     Authority = 1
	METRICS_MEMBER    Authority = 2
	ANONYMOUS         Authority = 3
)
