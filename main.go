package main

import (
	"changfang/database"
	SearchModel "changfang/search"
	"fmt"
)

func main()  {

	fmt.Println("---------开始爬取数据----------")
	//数据库初始化
	//database.InitDatabase()


	db := database.InitDatabase()
	//开启爬取模块
	SearchModel.Start(db)
}