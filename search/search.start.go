package SearchModel

import (
	ServiceModels "changfang/models"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	//"sync"
	Tools "xingqiu.co/utils/tools"
)

func SearchStartTask(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover err:", err)
		}
	}()

	//多线程爬取数据
	//var wg sync.WaitGroup
	//wg.Add(50)
	//从第50页开始爬取全部的数据---获得页面中的详情URL
	for page := 50; page >= 1; page-- {
		//go SearchStart(db, page, &wg)
		SearchStart(db, page)
	}
	//wg.Wait()
}

//搜索列表数据
//func SearchStart(db *gorm.DB, page int, wg *sync.WaitGroup) {
func SearchStart(db *gorm.DB, page int) {


	//列表的URL
	url := "http://sz.zhaoshang800.com/cf_list-0-1-17-0-0-0-0-0-0-0-0-2-" + strconv.Itoa(page) + ".html"
	//单次存储数据的量
	var count int64
	//获取请求返回的数据
	request := Tools.HttpGet(url)
	//转化为string  保存进str中
	str, err := request.String()

	if err != nil {
		Tools.TimeSleep(3)
		fmt.Println("接口刷新异常 err:", err)
	}
	//去空格
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	//Tools.FileWrite("./", "t.json", str)

	//正则表达式 ---发布者匹配
	strAuthor := `">([^<]*)</a><img src="tpl/newindex/images/jjr_rz.jpg" title="实名认证用户"/>`
	//正则表达式 ---标题匹配
	strTitle := `<h3  style=" ">([^<]*)</h3></a>`
	//正则表达式 ---匹配详情页面URL匹配
	strDetails := `<a class="l_pta fl[^/]*([^"]*)" target="_blank" title=`

	//正则表达式获取所有的URL
	compileDetails := regexp.MustCompilePOSIX(strDetails)
	hwDetailsUrl := compileDetails.FindAllStringSubmatch(str, -1)
	//匹配发布者
	compileAuthor := regexp.MustCompilePOSIX(strAuthor)
	hwAuthors := compileAuthor.FindAllStringSubmatch(str, -1)
	//匹配标题
	compileTitle := regexp.MustCompilePOSIX(strTitle)
	hwTitles := compileTitle.FindAllStringSubmatch(str, -1)



	//存储从页面获取下来的数据----（标题，发布者，详情URL）
	queues := make([]ServiceModels.QueueUrl, 30)
	//遍历将值传进来
	for key, hwDetailUrl := range hwDetailsUrl {
		queues[key].DetailURL = "http://sz.zhaoshang800.com" + hwDetailUrl[1]
	}
	for key, hwTitle := range hwTitles {
		if key >= 3 {
			queues[key-1].Title = hwTitle[1]
		} else {
			queues[key].Title = hwTitle[1]
		}
	}
	for key, hwAuthor := range hwAuthors {
		queues[key].Author = hwAuthor[1]
	}

	//保存不重复的值
	var queuesRepeat []ServiceModels.QueueUrl
	//遍历获取下来的数据
	for _, queue := range queues {

		//是否检测到空， 为空表示当前页面爬取完成
		if queue.DetailURL == ""{
			break
		}
		//判断是否重复
		var detailRepeat ServiceModels.QueueUrl
		//查询数据库中是否已经存在URL
		db.Where("detail_url=?", queue.DetailURL).First(&detailRepeat)
		//detailRepeat.DetailURL为空表示不存在  加入集合
		if detailRepeat.ID == 0 {
			count++
			queuesRepeat = append(queuesRepeat, queue)
		}
	}
	if queuesRepeat != nil {
		db.Create(&queuesRepeat)
	}
	fmt.Printf("爬取完第%d页-- 插入%d条数据\n", page, count)
	//wg.Done()
}

