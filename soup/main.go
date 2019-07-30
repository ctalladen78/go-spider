// 使用 github.com/anaskhan96/soup
// 抓取豆瓣Top250页面对应电影条目的ID和标题
package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func fetch(url string) (doc soup.Root, err error) {
	var (
		resp string
	)

	log.Printf("Fetch Url %s \n", url)

	// 设置请求头信息
	soup.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36")

	resp, err = soup.Get(url)
	doc = soup.HTMLParse(resp)

	return
}

func parseContent(url string) (err error) {
	var (
		doc       soup.Root
		roots     []soup.Root
		title, id string
	)

	if doc, err = fetch(url); err != nil {
		log.Println(err)
		return
	}

	roots = doc.Find("ol", "class", "grid_view").FindAll("div", "class", "hd")
	for _, r := range roots {
		url = r.Find("a").Attrs()["href"]
		id = strings.Split(url, "/")[4]
		title = r.Find("span", "class", "title").Text()

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
