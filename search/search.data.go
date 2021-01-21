package SearchModel

import (
	ServiceModels "changfang/models"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	Tools "xingqiu.co/utils/tools"
)

func SearchDataTask(db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover err:", err)
		}
	}()
	//开始爬取
	SearchDataStart(db)
}

func SearchDataStart(db *gorm.DB) {

	//插入总数量
	var countAll int64
	//存储插入数量
	var count int64
	//重复数量
	var countRepeat int64

	//存储所有URL的集合
	var SearchDetailUrls []ServiceModels.QueueUrl
	//查询数据表中所有详情URL
	db.Select("detail_url", "author").Find(&SearchDetailUrls).Count(&countAll)

	//遍历所有详情URL
	for dhy, SearchDetailUrl := range SearchDetailUrls {

		if dhy == 5 {
			break
		}

		var UrlRepeat ServiceModels.UrlContent
		//判断数据在表中是否存在，如果不存在添加进数据库中
		db.Where("url=?", SearchDetailUrl.DetailURL).First(&UrlRepeat)
		if UrlRepeat.Url == "" {
			//请求URL
			request := Tools.HttpGet(SearchDetailUrl.DetailURL)
			//转化为string  保存进str中
			str, err := request.String()
			if err != nil {
				Tools.TimeSleep(3)
				fmt.Println("接口刷新异常 err:", err)
			}
			//去特殊符号
			str = strings.Replace(str, "\n", "", -1)
			str = strings.Replace(str, "\t", "", -1)
			str = strings.Replace(str, "<br />", "", -1)
			str = strings.Replace(str, "<br/>", "", -1)

			//Tools.FileWrite("./", "t.json", str)

			//正则表达式class="mt_diqu f_wryh">[公明厂房]</a>
			//标题    1
			strTitle := `</a>([^<]*)</h1>`
			//区域    1
			strDistrict := `<a>所在区域：</a><span>([^<]*)</span>`
			//面积    1
			strArea := `<a>出租总面积：</a><span class="rede94">([^<]*)</span><span>平米</span>`
			//租金    1
			strRent := `<a>[^<]*租金：</a><span class="rede94">([^<]*)</span>`
			//楼层    1
			strFloor := `<a>所在楼层：</a><span>([^<]*)</span>`
			//亮点   多
			strHighlights := `<a class="details-tag-item" target="_blank" href="/cf_list-0-0-0-0-0-0-0-0-0-[^-]*--2.html"><span>([^<]*)</span></a>`
			//详细描述 1
			strDescription := `<div class="xq_cfgk">([^<]*)`
			//图片URL  多
			strImage := `<img class="lazyload" style="width:90%;height:90%;text-align:center" name="[^"]*"[^d]*data-src="([^"]*)"`

			posixTitle := regexp.MustCompilePOSIX(strTitle)
			hwTitle := posixTitle.FindStringSubmatch(str)

			posixDistrict := regexp.MustCompilePOSIX(strDistrict)
			hwDistrict := posixDistrict.FindStringSubmatch(str)

			posixArea := regexp.MustCompilePOSIX(strArea)
			hwArea := posixArea.FindStringSubmatch(str)

			posixRent := regexp.MustCompilePOSIX(strRent)
			hwRent := posixRent.FindStringSubmatch(str)

			posixFloor := regexp.MustCompilePOSIX(strFloor)
			hwFloor := posixFloor.FindStringSubmatch(str)

			posixDescription := regexp.MustCompilePOSIX(strDescription)
			hwDescription := posixDescription.FindStringSubmatch(str)

			posixHighlights := regexp.MustCompilePOSIX(strHighlights)
			hwHighlights := posixHighlights.FindAllStringSubmatch(str, -1)

			posixImage := regexp.MustCompilePOSIX(strImage)
			hwImages := posixImage.FindAllStringSubmatch(str, -1)

			var Highlights string
			for k, hwHighlight := range hwHighlights {
				if k == len(hwHighlights)-1 {
					Highlights = Highlights + hwHighlight[1]
					break
				}
				Highlights = Highlights + hwHighlight[1] + "."
			}

			//image   string类型拼接成数组
			var Images string
			for k, hwImage := range hwImages {

				if k == len(hwImages)-1 {
					Images = Images + hwImage[1]
				}

				Images = Images + hwImage[1] + ","

			}
			hwAreaInt, _ := strconv.Atoi(hwArea[1])
			hwRentInt, _ := strconv.Atoi(hwRent[1])

			/*fmt.Println("hwTitle", hwTitle[1])
			fmt.Println("hwArea", hwArea[1])
			fmt.Println("hwDistrict", hwDistrict[1])
			fmt.Println("hwDescription", hwDescription[1])
			fmt.Println("hwFloor", hwFloor[1])
			fmt.Println("hwRent", hwRent[1])
			fmt.Println("SearchDetailUrl.Author", SearchDetailUrl.Author)
			fmt.Println("SearchDetailUrl.DetailURL", SearchDetailUrl.DetailURL)
			fmt.Println("Degree",Degree)
			fmt.Println("Images",Images)*/

			//加入结构体
			HwUrlContent := ServiceModels.UrlContent{
				Author:   SearchDetailUrl.Author,
				Url:      SearchDetailUrl.DetailURL,
				Title:    hwTitle[1],
				District: hwDistrict[1],
				Area:     hwAreaInt,
				Money:    hwRentInt,
				Floor:    hwFloor[1],
				Detailed: hwDescription[1],
				Degree:   Highlights,
				Images:   Images,
			}
			//加入数据库
			db.Create(&HwUrlContent)
			count++
			fmt.Printf("往数据表url_content插入%d条数据,过滤重复数据%d条  共有%d条数据\n", count, countRepeat, countAll)

			if (count + countRepeat) == countAll {
				fmt.Println("----------------爬取完成------------------")
			}
		} else {
			countRepeat++
			fmt.Printf("往数据表url_content插入%d条数据,过滤重复数据%d条  共有%d条数据\n", count, countRepeat, countAll)
		}
	}
}
