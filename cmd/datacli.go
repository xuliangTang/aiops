package main

import "aipos/cmd/datacli"

// 用来导入或更新point数据
// go run cmd/datacli.go -f resources/vecdata/getter.yaml
func main() {
	cmd := datacli.NewDataCliCommand()
	// 执行命令
	cmd.Execute()
}
