package uuid

import (
	"fmt"
	"testing"
)

func TestGenUUID(t *testing.T) {
	uuid := GenUUID()
	fmt.Println(uuid)
}
