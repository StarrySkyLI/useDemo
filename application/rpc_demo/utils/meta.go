package utils

import (
	"fmt"
	"strings"

	"useDemo/application/rpc_demo/rpc"
)

func GetSuccessMeta() *rpc.Metadata {
	return &rpc.Metadata{
		BizCode: 0,
		Message: "success",
	}
}

func GetMetaWithError(err error) *rpc.Metadata {
	return &rpc.Metadata{
		BizCode: -1,
		Message: err.Error(),
	}
}

func GetMetaWithErrorString(err string) *rpc.Metadata {
	return &rpc.Metadata{
		BizCode: -1,
		Message: err,
	}
}
func GetObjectName(bizName string, Id int64) string {
	return fmt.Sprintf("%s/%d", bizName, Id)
}
func CheckHash(hash1, hash2 string) bool {
	// 将两个哈希值都转换为小写
	fHashLower := strings.ToLower(hash1)
	hashLower := strings.ToLower(hash2)

	// 比较小写哈希值
	return fHashLower == hashLower
}
