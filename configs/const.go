package configs

import "time"

var SignKey = []byte("asd@#lskd2!aw32k34242WSASdsk32")

const (
	AccessExpireTime  = time.Minute * 1
	RefreshExpireTime = time.Hour * 240
	Layout            = "2006-01-02 15:04:05"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)
