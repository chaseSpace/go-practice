package util

import (
	"fmt"
	"strings"
)

func formatDuration(seconds int, inChinese ...bool) string {
	var hours, minutes, secs int

	// 计算小时数
	hours = seconds / 3600
	// 计算剩余的分钟数
	minutes = (seconds % 3600) / 60
	// 计算剩余的秒数
	secs = seconds % 60

	// 构建输出字符串
	var output string
	if hours > 0 {
		output += fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		output += fmt.Sprintf("%dm", minutes)
	}
	if secs > 0 {
		output += fmt.Sprintf("%ds", secs)
	} else if output == "" { // 如果没有计算出任何时间单位，则默认为秒
		output = fmt.Sprintf("%ds", seconds)
	}
	if len(inChinese) > 0 && inChinese[0] {
		output = strings.Replace(output, "h", "小时", 1)
		output = strings.Replace(output, "m", "分钟", 1)
		output = strings.Replace(output, "s", "秒", 1)
	}
	return output
}
