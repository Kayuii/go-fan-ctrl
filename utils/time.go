package utils

import "time"

const (
	// 默认格式化为 Y-M-D h:m:s 的模板,死变态老把时间模板设置到这个时间点格式
	DEF_TIME_FORMAT     = "2006-01-02 15:04:05"
	TIME_FORMAT_YMD     = "20060102"   // 时间格式：20140725
	TIME_FORMAT_YYYMMDD = "2006-01-02" // 时间格式：2014-07-25
)

// 把时间 time.Time 格式化【格式错误返回的信息也是不对版的】
func TimeFormat(t time.Time, fmt string) string {
	// Time类型有这几个函数可以获取对于的值：t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()
	// 另外 t.Format() 函数可以按照模板格式化
	// 源码：go/src/pkg/time/format.go 有定义其它的格式化模板
	return t.Format(fmt)
}
