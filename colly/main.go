// 使用 github.com/gocolly/colly
// 抓取豆瓣Top250页面对应电影条目的ID和标题
// 报错：cannot find package "google.golang.org/appengine/urlfetch" in any of
// 原因：google.golang.org/appengine 这个包的代码仓库变了，指向 https://github.com/golang/appengine
// 解决办法：将下载目录下的 github.com/golang/appengine 复制到 $gopath/src/google.golang.org/appengine
package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

func main() {
	var (
		c *colly.Collector
	)

	// 异步抓取
	c = colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36"),
	)

	// 只抓取域名是douban(域名后缀和二级域名不限制)的地址
	// 限制并发是5
	c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5})

	// 请求前
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// 请求过程中发生错误
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// 收到响应后
	// 使用 xPath 解析
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Println(err)
		}

		nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
		for _, node := range nodes {
			// 寻找符合的第一个节点，并转化为字符串
			title := htmlquery.InnerText(htmlquery.FindOne(node, `.//span[@class="title"]/text()`))
			id := strings.Split(htmlquery.InnerText(htmlquery.FindOne(node, `./a/@href`)), "/")[4]

			fmt.Println(title, id)
		}
	})

	// 默认使用 goquery
	c.OnHTML(".hd", func(e *colly.HTMLElement) {
		id := strings.Split(e.ChildAttr("a", "href"), "/")[4]

		// title := e.ChildText("span.title")
		// 类名为 title 的 span 标签有2个，用 ChildText 会直接返回2个标签的全部的值

		// 找第一个符合的文本
		title := e.DOM.Find("span.title").Eq(0).Text()

		fmt.Println(title, id)
	})

	// 如果收到的响应内容是 HTML 调用它
	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
	})

	// 在 OnHTML 回调完成后调用
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://movie.douban.com/top250?start=0&filter=")

	c.Wait()
}
