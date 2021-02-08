package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/stefanoschrs/proxymeister/pkg/types"
	"github.com/stefanoschrs/proxymeister/pkg/utils"
)

var Name = "sslproxies.org"
var Url = "https://www.sslproxies.org/"

func Fetch() (proxies []types.Proxy, err error) {
	body, err := utils.HttpProxyfiedGet(types.ProxifiedGetOptions{
		TargetUrl: Url,
	})
	if err != nil {
		return
	}

	return process(body)
}

func process(body []byte) (proxies []types.Proxy, err error) {
	pattern, err := regexp.Compile(`(?s)<textarea.+>(.+)</textarea>`)
	if err != nil {
		return
	}
	result := pattern.FindAllSubmatch(body, -1)

	lines := strings.Split(string(result[0][1]), "\n")
	for _, line := range lines[3 : len(lines)-1] {
		s := strings.Split(line, ":")

		port, err2 := strconv.ParseInt(s[1], 10, 64)
		if err2 != nil {
			//log.Println(err2)
			continue
		}

		proxy := types.Proxy{
			Ip:     strings.TrimSpace(s[0]),
			Port:   int(port),
			Status: types.ProxyStatusFresh,
			Source: Name,
		}
		proxies = append(proxies, proxy)
	}

	return
}

func main() {
	body, err := ioutil.ReadFile("pkg/crawler/testdata/sslproxies.org.golden")
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
