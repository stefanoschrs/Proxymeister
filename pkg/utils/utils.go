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
