package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	color "github.com/fatih/color"
	sitemap "github.com/yterajima/go-sitemap"
)

var (
	sitemapURL = flag.String("sitemap", "", "sitemap.xml/.xml.gz的URL")

	status100 = regexp.MustCompile(`1\d\d`)
	status200 = regexp.MustCompile(`2\d\d`)
	status300 = regexp.MustCompile(`3\d\d`)
	status400 = regexp.MustCompile(`4\d\d`)
	status500 = regexp.MustCompile(`5\d\d`)
)
var complete chan int = make(chan int)

func loop(smap sitemap.Sitemap) {
	ic := 0
	cr := 0
	er := 0
	c4 := 0
	c5 := 0

	for _, URL := range smap.URL {
		time.Sleep(time.Second)
		ic++
		resp, err := http.Head(URL.Loc)
		if err != nil {
			fmt.Println(err)
			er++
			continue
		}
		defer resp.Body.Close()
		switch statusType(resp.StatusCode) {
		case 100:
			color.Cyan(resp.Status + " " + URL.Loc)
		case 200:
			color.Green(resp.Status + " " + URL.Loc)
			continue
		case 300:
			color.Magenta(resp.Status + " " + URL.Loc)
			cr++
		case 400:
			color.Red(resp.Status + " " + URL.Loc)
			er++
			c4++
		case 500:
			color.Yellow(resp.Status + " " + URL.Loc)
			er++
			c5++
		default:
			color.White(resp.Status + " " + URL.Loc)
		}
	}
	color.Yellow("total" + ":" + strconv.Itoa(ic))
	color.Yellow("error" + ":" + strconv.Itoa(er))
	color.Yellow("redirect" + ":" + strconv.Itoa(cr))
	color.Yellow("404" + ":" + strconv.Itoa(c4))
	color.Yellow("500" + ":" + strconv.Itoa(c5))
	complete <- 0 // 执行完毕了，发个消息
}
func main() {
	flag.Parse()

	if *sitemapURL == "" {
		fmt.Println("請指定 sitemap 的網址 ex: -sitemap  http://dayi.demo.ready-market.com/sitemap.xml")
		return
	}

	smap, err := sitemap.Get(*sitemapURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	go loop(smap)
	<-complete // 直到线程跑完, 取到消息. main在此阻塞住
}

func statusType(statusCode int) int {
	var statusType int

	statusCodeStr := strconv.Itoa(statusCode)
	switch {
	case status100.MatchString(statusCodeStr):
		statusType = 100
	case status200.MatchString(statusCodeStr):
		statusType = 200
	case status300.MatchString(statusCodeStr):
		statusType = 300
	case status400.MatchString(statusCodeStr):
		statusType = 400
	case status500.MatchString(statusCodeStr):
		statusType = 500
	default:
		statusType = 0
	}

	return statusType
}
