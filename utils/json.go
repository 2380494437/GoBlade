package utils

import (
	"encoding/json"
	"errors"
)

// JsonEncode 将任意对象编码为 JSON 字符串
func JsonEncode(data any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JsonDecode 将 JSON 字符串解析为指定结构或 map
func JsonDecode(jsonStr string, target any) error {
	if jsonStr == "" {
		return errors.New("json字符串为空")
	}
	return json.Unmarshal([]byte(jsonStr), target)
}
