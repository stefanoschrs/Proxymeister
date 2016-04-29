package crawler

import (
    "fmt"
    "regexp"
    "strings"
	"strconv"
	"io/ioutil"
    "net/http"
)

const url string = "http://www.sslproxies.org/"

type Proxy struct {
    Ip 		string
    Port 	int
}

func FetchProxies() []Proxy{
	var proxies []Proxy

	response, err := http.Get(url)
    if err != nil {
        fmt.Printf("%s", err)
    } else {
        defer response.Body.Close()

        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            fmt.Printf("%s", err)
        } else {

            re1, _ := regexp.Compile(`<tr>(<td>.+<\/td>)+<\/tr>`)
            result:= re1.FindAllStringSubmatch(string(contents), -1)
			for _, v := range result {
                array := strings.Split(v[0][8:len(v[0])-10], "</td><td>")

                p := Proxy{}
				p.Ip = array[0]
				if s, err := strconv.ParseInt(array[1], 10, 32); err == nil {
					p.Port = int(s)
				}

                proxies = append(proxies, p)
            }
        }
    }

	return proxies
}
