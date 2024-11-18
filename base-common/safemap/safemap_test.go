package safemap

import (
	"fmt"
	"testing"
)

func TestSafeMap(t *testing.T) {

	safeMap := NewSafeMap(10)
	safeMap.Set("a", "aaa")
	fmt.Println(safeMap.Get("a"))

}
