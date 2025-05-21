package snowflake

import "github.com/bwmarrin/snowflake"

var defaultNode *snowflake.Node

func InitDefaultSnowflakeNode(node int64) {
	if node == 0 {
		node = 1
	}
	var err error
	defaultNode, err = snowflake.NewNode(node)
	if err != nil {
		panic(err)
	}
}

func GetSnowflakeId() int64 {
	return defaultNode.Generate().Int64()
}
