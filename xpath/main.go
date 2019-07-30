// 使用 github.com/antchfx/htmlquery
// 抓取豆瓣Top250页面对应电影条目的ID和标题
// / 从当前节点选取直接子节点
// // 从当前节点选取子孙节点
// . 选取当前节点
// .. 选取当前节点的父节点
// @ 选取属性
package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func fetch(url string) (doc *html.Node, err error) {
	var (
		res    *http.Response
		client http.Client
		req    *http.Request
	)

	log.Printf("Fetch Url %s \n", url)

	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	// 设置请求头信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")

	client = http.Client{}

	if res, err = client.Do(req); err != nil {
		return
	}
	defer res.Body.Close()

	doc, err = htmlquery.Parse(res.Body)

	return
}

func parseContent(url string) (err error) {
	var (
		doc       *html.Node
		nodes     []*html.Node
		title, id string
	)

	if doc, err = fetch(url); err != nil {
		log.Println(err)
		return
	}

	nodes = htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
	for _, node := range nodes {
		// 寻找符合的第一个节点，并转化为字符串
		title = htmlquery.InnerText(htmlquery.FindOne(node, `.//span[@class="title"]/text()`))
		id = strings.Split(htmlquery.InnerText(htmlquery.FindOne(node, `./a/@href`)), "/")[4]

		fmt.Println(title, id)
	}

	return
}

func main() {
	var (
		startTime time.Time
		interval  time.Duration
		wg        sync.WaitGroup
	)

	startTime = time.Now()

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			parseContent("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
		}(i)
	}

	wg.Wait()

	interval = time.Since(startTime)
	fmt.Printf("Take %s", interval)
}
