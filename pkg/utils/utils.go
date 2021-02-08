package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stefanoschrs/proxymeister/internal/logging"
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := LoadEnv()
	if err != nil {
		return err
	}

	return nil
}

func LoadEnv() (err error) {
	p := os.Getenv("ENV")
	if p == "" {
		p = ".env"
	}

	return godotenv.Load(p)
}

func GetMyIP() (ip string, err error) {
	res, err := http.Get("http://bot.whatismyipaddress.com")
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	ip = string(body)
	return
}

func HttpGetJson(targetUrl string, result interface{}) (err error) {
	res, err := http.Get(targetUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	return
}

func HttpProxyfiedGet(options types.ProxifiedGetOptions) (body []byte, err error) {
	var activeProxies []types.ProxyFE
	if os.Getenv("PROXYMEISTER_URL") != "" {
		limit := uint(10)
		if options.Tries != nil {
			limit = *options.Tries
		}
		err = HttpGetJson(fmt.Sprintf("%s/proxies?limit=%d&fields=address,ip,port", os.Getenv("PROXYMEISTER_URL"), limit), &activeProxies)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest(http.MethodGet, options.TargetUrl, nil)
	if err != nil {
		err = fmt.Errorf("http.NewRequest: %w", err)
		return
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	var client *http.Client

	if len(activeProxies) == 0 {
		client = &http.Client{}
		res, doErr := client.Do(req)
		if doErr != nil {
			err = fmt.Errorf("client.Do: %w", doErr)
			return
		}
		defer res.Body.Close()

		return ioutil.ReadAll(res.Body)
	}

	rand.Shuffle(len(activeProxies), func(i, j int) { activeProxies[i], activeProxies[j] = activeProxies[j], activeProxies[i] })

	var resBody *io.ReadCloser
	for i, proxy := range activeProxies {
		logging.Debugf("Try %d/%d. Using Proxy: %s", i+1, len(activeProxies), *proxy.Address)

		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Host: *proxy.Address,
				}),
			},
			Timeout: 15 * time.Second,
		}

		res, doErr := client.Do(req)
		if doErr != nil {
			continue
		}
		if res.StatusCode != http.StatusOK {
			continue
		}

		resBody = &res.Body
		break
	}
	if resBody == nil {
		err = fmt.Errorf("cannot reach: %s", options.TargetUrl)
		return
	}
	defer (*resBody).Close()

	return ioutil.ReadAll(*resBody)
}
