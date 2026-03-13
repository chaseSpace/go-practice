package lumberjack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type FileMode uint32

// MarshalJSON implements json.Marshaler
func (fm *FileMode) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"0%03o"`, uint32(*fm))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (fm *FileMode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// 解析八进制字符串
	mode, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return err
	}
	*fm = FileMode(mode)
	return nil
}

// ToOS Convert FileMode to os.FileMode
func (fm FileMode) ToOS() os.FileMode {
	return os.FileMode(fm)
}
