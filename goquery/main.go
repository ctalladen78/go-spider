// 使用 github.com/PuerkitoBio/goquery
// 抓取豆瓣Top250页面对应电影条目的ID和标题
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func fetch(url string) (doc *goquery.Document, err error) {
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

	doc, err = goquery.NewDocumentFromReader(res.Body)

	return
}

func parseContent(url string) (err error) {
	var (
		doc       *goquery.Document
		title, id string
	)

	if doc, err = fetch(url); err != nil {
		log.Println(err)
		return
	}

	// 使用 goquery 语法
	doc.Find("ol.grid_view li").Each(func(i int, selection *goquery.Selection) {
		var (
			url    string
			exists bool
		)

		url, exists = selection.Find("hd a").Attr("href")
		if exists {
			id = strings.Split(url, "/")[4]
		}
		title = selection.Find("hd .title").Eq(0).Text()

		fmt.Println(title, id)
	})

	return
}

func main() {
	var (
		startTime time.Time
		interval  time.Duration
		wg        sync.WaitGroup
	)

	startTime = time.Now()

	for i := 0; i < 10; i++ {
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
