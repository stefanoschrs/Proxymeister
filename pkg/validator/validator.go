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
	"time"

	"github.com/stefanoschrs/proxymeister/pkg/types"
)

const validationTimeout = 10 * time.Second

func getValidationUrl(ssl bool) string {
	mixUrls := []string{
		"bot.whatismyipaddress.com",
		"ipv4bot.whatismyipaddress.com",
		"ipinfo.io/ip",
		//"icanhazip.com",
	}

	httpUrls := []string{}

	httpsUrls := []string{
		//"nordvpn.com/wp-admin/admin-ajax.php?action=get_user_info_data",
		//"api.myip.com",
		//"ip4.seeip.org",
		//"ipapi.co/ip",
		"api.ipify.org",
		//"api.my-ip.io/ip.txt",
		//"api4.my-ip.io/ip.txt",
	}

	urls := mixUrls
	prefix := "http"
	if ssl {
		prefix += "s"
		urls = append(urls, httpsUrls...)
	} else {
		urls = append(urls, httpUrls...)
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

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("https status: %d - %s", res.StatusCode, res.Status)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// SQUID misconfiguration
	//if strings.Contains(string(body), "<body id=\"ERR_ACCESS_DENIED\">") {
	//	err = errors.New(types.ErrValidationBadSquid)
	//	return
	//}

	rgx, err := regexp.Compile("[0-9]{1,3}(\\.[0-9]{1,3}){3}")
	if err != nil {
		return
	}
	resultIp := rgx.Find(body)

	//fmt.Printf("My: %s\tProxy: %s\tMatch: %s\tBody: %s\n", myIp, ip, resultIp, string(body))

	if string(resultIp) == myIp {
		err = errors.New(types.ErrValidationPassthrough)
		return
	}
	if string(resultIp) == "" {
		err = errors.New(types.ErrEmptyResponse)
		return
	}
	if net.ParseIP(string(resultIp)) == nil {
		err = errors.New(types.ErrInvalidIp)
		return
	}

	latency = (time.Now().UnixNano() - startTime.UnixNano()) / 1e6
	return
}
