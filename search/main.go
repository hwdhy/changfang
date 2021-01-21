package SearchModel

import "gorm.io/gorm"

func Start(db *gorm.DB) {

	//开启列表爬取
	SearchStartTask(db)

	//详情数据爬取
	SearchDataTask(db)
}
