package utils

import (
	"testing"
	"time"
)

const SYSINFO_UPTIME uint64 = 391413

func TestTimeFormat(t *testing.T) {

	uptime := time.Now().Add(-time.Second * time.Duration(SYSINFO_UPTIME))

	t.Logf("uptime value: %+v", TimeFormat(uptime, DEF_TIME_FORMAT))
	t.Logf("uptime value: %+v", TimeFormat(uptime, TIME_FORMAT_YMD))
	t.Logf("uptime value: %+v", TimeFormat(uptime, TIME_FORMAT_YYYMMDD))
}
