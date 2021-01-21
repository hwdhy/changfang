package ServiceModels

type QueueUrl struct {
	ID        int64  `json:"id" gorm:"PRIMARY_KEY"`
	Title     string `json:"title"`                   //标题
	Author    string `json:"author"`                  //发布者
	DetailURL string `json:"detail_url" gorm:"index"` //详情URL
}

func (QueueUrl) TableName() string {
	return "queue_url"
}
