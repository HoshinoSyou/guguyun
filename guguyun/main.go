package main

import (
	"guguyun/cmd"
	"guguyun/dao"
)

func main() {
	dao.DB = dao.Init()
	defer dao.DB.Close()
	cmd.Entrance()
}
