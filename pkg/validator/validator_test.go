package validator

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	myIp := getMyIp()

	testCases := []struct {
		Ip   string
		Port int
		Ssl  bool
	}{
		// Success
		{"159.8.114.34", 8123, false},
		{"159.8.114.34", 8123, true},

		// Basic checks
		//{"129.8.114.37", 80, false},
		//{"129.8.114.37", 66000, false},
		//{"256.8.114.37", 80, false},
	}

	for _, testCase := range testCases {
		latency, err := Validate(myIp, testCase.Ip, testCase.Port, false)
		if err != nil {
			if strings.Contains(err.Error(), "deadline exceeded") {
				log.Println("Timeout")
				continue
			}
			if strings.Contains(err.Error(), "invalid port") {
				log.Println("Invalid Port")
				continue
			}
			if strings.Contains(err.Error(), "no such host") {
				log.Println("Invalid IP")
				continue
			}
			if strings.Contains(err.Error(), "Proxy Authentication Required") {
				log.Println("Locked proxy")
				continue
			}

			log.Println(err)
		} else {
			log.Printf("Proxy: ok\tSSL: %v\tLatency: %dms\n", testCase.Ssl, latency)
		}
	}
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	m.Run()
}

func getMyIp() string {
	res, err := http.Get("http://bot.whatismyipaddress.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
