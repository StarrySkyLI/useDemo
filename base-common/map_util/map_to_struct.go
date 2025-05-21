package map_util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// MapToStruct 将 map[string]any 转为结构体 我觉得我现在强的可怕！
func MapToStruct(interfaceMap any, output any) error {

	val := reflect.ValueOf(output)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("output 必须是非空结构体指针")
	}
	if val.Elem().Kind() != reflect.Struct {
		return errors.New("output 必须指向结构体")
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    "mapstructure",
		DecodeHook: createDecodeHook(),
		Result:     output,
	})
	if err != nil {
		return fmt.Errorf("创建解码器失败: %v", err)
	}

	// 执行解码
	if err := decoder.Decode(interfaceMap); err != nil {
		return fmt.Errorf("解码失败: %v", err)
	}
	return nil
}

func createDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// 仅处理源类型为字符串的转换
		if f.Kind() != reflect.String {
			return data, nil
		}
		strVal := data.(string)

		switch t.Kind() {
		case reflect.Int:
			return strconv.Atoi(strVal)
		case reflect.Bool:
			return strconv.ParseBool(strVal)
		case reflect.Float64:
			return strconv.ParseFloat(strVal, 64)
		case reflect.Int64:
			return strconv.ParseInt(strVal, 10, 64)
		case reflect.Slice:
			return handleSliceConversion(t, strVal)

		}

		return data, nil
	}
}

func handleSliceConversion(t reflect.Type, strVal string) (interface{}, error) {
	elemType := t.Elem()

	if strings.HasPrefix(strVal, "[") && strings.HasSuffix(strVal, "]") {
		slicePtr := reflect.New(t)
		if err := json.Unmarshal([]byte(strVal), slicePtr.Interface()); err == nil {
			return slicePtr.Elem().Interface(), nil
		}
	}

	parts := strings.Split(strVal, ",")
	slice := reflect.MakeSlice(t, len(parts), len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)

		switch elemType.Kind() {
		case reflect.String:
			slice.Index(i).SetString(part)
		case reflect.Int:
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("元素 '%s' 无法转换为 int", part)
			}
			slice.Index(i).SetInt(int64(val))
		case reflect.Int64:
			val, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("元素 '%s' 无法转换为 int64", part)
			}
			slice.Index(i).SetInt(val)
		case reflect.Float64:
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("元素 '%s' 无法转换为 float64", part)
			}
			slice.Index(i).SetFloat(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(part)
			if err != nil {
				return nil, fmt.Errorf("元素 '%s' 无法转换为 bool", part)
			}
			slice.Index(i).SetBool(val)
		default:
			return nil, fmt.Errorf("不支持的切片元素类型: %v", elemType.Kind())
		}
	}

	return slice.Interface(), nil
}
