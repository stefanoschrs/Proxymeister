package validator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/stefanoschrs/proxymeister/pkg/types"
)

const validationTimeout = 10 * time.Second

func getValidationUrl(ssl bool) string {
	urls := []string{
		// Very Good
		"myip.dnsomatic.com",
		"bot.whatismyipaddress.com",
		"ipv4bot.whatismyipaddress.com",
		"icanhazip.com",

		// Sometimes SQUID issues
		//"cpanel.com/showip.shtml",
		//"ipecho.net/plain",
		//"checkip.amazonaws.com",
		//"myexternalip.com/raw",
		//"ip-api.com/line/?fields=query",
		//"api.duckduckgo.com/?q=my+ip&format=xml",
	}

	prefix := "http"
	if ssl {
		prefix += "s"
	}

	return prefix + "://" + urls[rand.Intn(len(urls))]
}

func Validate(myIp, ip string, port int, ssl bool) (latency int64, err error) {
	startTime := time.Now()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: fmt.Sprintf("%s:%d", ip, port),
			}),
		},
		Timeout: validationTimeout,
	}
	req, err := http.NewRequest(http.MethodGet, getValidationUrl(ssl), nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// SQUID misconfiguration
	if strings.Contains(string(body), "<body id=\"ERR_ACCESS_DENIED\">") {
		err = errors.New(types.ErrValidationBadSquid)
		return
	}

	//fmt.Printf("My: %s\tProxy: %s\tBody: %s\n", myIp, ip, string(body))

	rgx, err := regexp.Compile("[0-9]{1,3}(\\.[0-9]{1,3}){3}")
	if err != nil {
		return
	}
	resultIp := rgx.Find(body)
	if string(resultIp) == myIp {
		err = errors.New(types.ErrValidationPassthrough)
		return
	}
	if net.ParseIP(string(resultIp)) == nil {
		err = errors.New(types.ErrInvalidIp)
		return
	}

	latency = (time.Now().UnixNano() - startTime.UnixNano()) / 1e6
	return
}
