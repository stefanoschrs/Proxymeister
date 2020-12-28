package utils

import (
	"io/ioutil"
	"net/http"
	"os"

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

func HttpGet(url string) (body []byte, err error) {
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
