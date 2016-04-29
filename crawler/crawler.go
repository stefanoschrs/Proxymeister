package crawler

import (
    "regexp"
    "strings"
	"strconv"

    "github.com/parnurzeal/gorequest"
)

const url string = "http://www.sslproxies.org/"

type Proxy struct {
    Ip 		string
    Port 	int
}

func FetchProxies() []Proxy{
	var proxies []Proxy

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		return proxies
	}

	pattern, _ := regexp.Compile(`<tr>(<td>.+<\/td>)+<\/tr>`)
	result := pattern.FindAllStringSubmatch(body, -1)
	for _, v := range result {
		array := strings.Split(v[0][8:len(v[0])-10], "</td><td>")

		p := Proxy{}
		p.Ip = array[0]
		if s, err := strconv.ParseInt(array[1], 10, 32); err == nil {
			p.Port = int(s)
		}

		proxies = append(proxies, p)
	}

	return proxies
}
