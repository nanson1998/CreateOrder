package common

import (
	"fmt"
	"math/rand"
	"time"
)

var DefaultId = 0

func init() {
	DefaultId = 1
}

func GetTransID(apptime int64) string {
	DefaultId = rand.Intn(1000000)

	DefaultId += 1
	t := time.Now().Format("150405") //HH:mm:ss
	return fmt.Sprintf("%v_%v%05d", GetTimeString(apptime), t, DefaultId)
}

func GetAppTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTimeString(apptime int64) string {
	t := time.Unix(apptime/1000, 0)
	return fmt.Sprintf("%02d%02d%02d", t.Year()%100, int(t.Month()), t.Day())
}

// yymmdd_appid_xxxxxxxxxx
func GetMRefundId(appid string, apptime int64) string {
	if DefaultId >= 100000 {
		DefaultId = 0
	}

	DefaultId += 1
	t := time.Now().Format("150405")
	return fmt.Sprintf("%v_%v_%v%05d", GetTimeString(apptime), appid, t, DefaultId)
}
