package ServiceModels

import "time"

type UrlContent struct {
	ID        int64     `json:"id" gorm:"PRIMARY_KEY"`       //ID
	Author    string    `json:"author"`                      //发布人
	Url       string    `json:"url"`                         //详情URL
	Title     string    `json:"title"`                       //标题       title
	District  string    `json:"district"`                    //区域
	Area      int       `json:"area"`                        //面积        area
	Money     int       `json:"money"`                       //租金        money
	Floor     string    `json:"floor"`                       //楼层
	Detailed  string    `json:"detailed"`                    //详细描述    detailed
	Degree    string    `json:"degree"`                      //亮点      degree
	Images    string    `json:"image" gorm:"column:fileIDs"` //图片
	CreatedAt time.Time `json:"created_at" gorm:"column:createtime"`
}

func (UrlContent) TableName() string {
	return "url_content"
}
