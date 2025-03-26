package util

import (
	"fmt"
	"testing"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds   int
		inChinese bool
		expected  string
	}{
		{3661, false, "1h1m1s"},
		{3661, true, "1小时1分钟1秒"},
		{3601, false, "1h1s"},
		{3601, true, "1小时1秒"},
		{61, false, "1m1s"},
		{61, true, "1分钟1秒"},
		{59, false, "59s"},
		{59, true, "59秒"},
		{0, false, "0s"},
		{0, true, "0秒"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d秒", tt.seconds), func(t *testing.T) {
			got := formatDuration(tt.seconds, tt.inChinese)
			if got != tt.expected {
				t.Errorf("formatDuration(%d) = %s; want %s", tt.seconds, got, tt.expected)
			}
		})
	}
}
