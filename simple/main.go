// 使用标准库
// 抓取豆瓣Top250页面对应电影条目的ID和标题
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func fetch(url string) (content string, err error) {
	var (
		res    *http.Response
		bytes  []byte
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

	if bytes, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}

	content = string(bytes)

	return
}

func parseContent(url string) (err error) {
	var (
		content             string
		hdRe, titleRe, idRe *regexp.Regexp
		items               [][]string
		title, id           string
	)

	if content, err = fetch(url); err != nil {
		log.Println(err)
		return
	}

	// 使用正则表达式
	hdRe = regexp.MustCompile(`<div class="hd">((.|\n)*?)</div>`)
	titleRe = regexp.MustCompile(`<span class="title">(.*?)</span>`)
	idRe = regexp.MustCompile(`<a href="https://movie.douban.com/subject/(\d+)/"`)

	// 把符合正则表达式的结果都解析出来(返回一个切片)
	items = hdRe.FindAllStringSubmatch(content, -1)
	for _, item := range items {
		// 找第一个符合的结果
		title = titleRe.FindStringSubmatch(item[0])[1]
		id = idRe.FindStringSubmatch(item[0])[1]
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
