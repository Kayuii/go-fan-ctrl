package utils

import (
	"fmt"
)

// 秒时间格式化为可视字符串显示【110 Days 10:20:30】
func SecTimeHuman(sec uint64) (human string) {
	var year, day, hour, minute, second uint64
	if sec >= 31536000 {
		year = sec / 31536000
		sec = sec % 31536000
		human = fmt.Sprintf("%v Year ", year)
	}

	if sec >= 86400 {
		day = sec / 86400
		sec = sec % 86400
		human += fmt.Sprintf("%v Days ", day)
	}
	if sec >= 3600 {
		hour = sec / 3600
		sec = sec % 3600
	}
	if sec >= 60 {
		minute = sec / 60
		second = sec % 60
	}
	if sec < 60 {
		second = sec
	}
	human += fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	return
}
