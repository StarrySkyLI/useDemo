package test

import (
	"fmt"
	"gitlab.coolgame.world/go-template/base-common/arLanguage"
	"testing"
)

func TestIsExist(t *testing.T) {
	// 示例数据
	data := map[string]interface{}{
		"foo": "Hello",
		"bar": []interface{}{
			"apple",
			map[string]interface{}{
				"nested": "banana",
			},
			"cherry",
		},
	}

	fmt.Println("Before modification:", data)
	res := arLanguage.SwitchLanguage(data, "india")
	fmt.Println("After modification:", res)
}
