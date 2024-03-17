package utils

import (
	"encoding/json"
	"reflect"
	"time"
)

// 一些工具函数
func Paginate[T any](data []T, pageNum int, pageSize int) []T {
	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex > len(data) {
		return []T{}
	}
	if endIndex >= len(data) {
		return data[startIndex:]
	}
	return data[startIndex:endIndex]
}

func SelfMarshal[T any](data T) ([]byte, error) {
	res := make(map[string]interface{})

	v := reflect.ValueOf(data)
	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		// 获取字段名
		fieldName := field.Name
		// 获取字段值
		fieldValue := v.Field(i)

		// 检查字段是否可导出
		if field.PkgPath == "" {
			// 字段可导出，转换为interface{}
			// 断言可以得到具体值
			valueInterface := fieldValue.Interface()
			if fieldValue.Kind() == reflect.Slice && fieldValue.IsNil() {
				res[fieldName] = []int{}
			} else {
				res[fieldName] = valueInterface
			}
		}
	}
	return json.Marshal(res)
}

func GetCurrentTimeStr() string {
	currentTime := time.Now()
	const layout = "2006-01-02 15:04:05"
	timeString := currentTime.Format(layout)
	return timeString
}
