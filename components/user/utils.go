package user

import (
	"fmt"
	"log"
)

// PrintAndReturnErrorf 打印错误信息并返回格式化的错误
// 这个函数用于在记录错误的同时返回错误
func PrintAndReturnErrorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	log.Printf("Error: %v", err)
	return err
}
