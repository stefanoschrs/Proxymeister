package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/stefanoschrs/proxymeister/pkg/types"
	"github.com/stefanoschrs/proxymeister/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

var Name = "hidemy.name"
var Url = "https://hidemy.name/en/proxy-list/?maxtime=3000&type=s&anon=4#list/"

func Fetch() (proxies []types.Proxy, err error) {
	body, err := utils.HttpGet(Url)
	if err != nil {
		return
	}

	return process(body)
}

func process(body []byte) (proxies []types.Proxy, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// IP
		// Port
		// Location
		// Speed
		// Type
		// Anonymity
		// Latest update

		port, err2 := strconv.ParseInt(goquery.NewDocumentFromNode(s.Find("td").Get(1)).Text(), 10, 64)
		if err2 != nil {
			//log.Println(err2)
			return
		}

		proxies = append(proxies, types.Proxy{
			Ip:     goquery.NewDocumentFromNode(s.Find("td").Get(0)).Text(),
			Port:   int(port),
			Source: Name,
		})
	})

	return
}

func main() {
	body, err := ioutil.ReadFile("pkg/crawler/testdata/hidemy.name.golden")
	if err != nil {
		log.Fatal(err)
	}

	proxies, err := process(body)
	if err != nil {
		log.Fatal(err)
	}

	for _, proxy := range proxies {
		log.Printf("%+v\n", proxy)
	}
}
