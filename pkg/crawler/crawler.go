package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/stefanoschrs/proxymeister/pkg/types"
)

func httpGet(url string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func fetchProxiesFromSource(source types.ProxySource) (proxies []types.Proxy, err error)  {
	body, err := httpGet(source.Url)
	if err != nil {
		return
	}

	pattern, err := regexp.Compile(`(?s)<textarea.+>(.+)<\/textarea>`)
	if err != nil {
		return
	}
	result := pattern.FindAllSubmatch(body, -1)

	lines := strings.Split(string(result[0][1]), "\n")
	for _, line := range lines[3 : len(lines)-1] {
		s := strings.Split(line, ":")

		port, err2 := strconv.ParseInt(s[1], 10, 64)
		if err2 != nil {
			log.Println(err2)
			continue
		}

		proxy := types.Proxy{
			Ip:     s[0],
			Port:   int(port),
			Source: source.Name,
		}
		proxies = append(proxies, proxy)
	}

	return
}

func FetchProxies() (proxies []types.Proxy, err error) {
	// TODO: Move to config file with specific parsers
	proxySources := []types.ProxySource{
		{
			Name: "sslproxies.org",
			//Url: "http://0.0.0.0:5000/sslproxies.org.golden",
			Url: "https://www.sslproxies.org/",
		},
	}

	for _, source := range proxySources {
		var sourceProxies []types.Proxy

		sourceProxies, err = fetchProxiesFromSource(source)
		if err != nil {
			err = fmt.Errorf("failed to fetch %s. %w\n", source.Name, err)
			return
		}

		proxies = append(proxies, sourceProxies...)
	}

	return
}
