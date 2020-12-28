package validator

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stefanoschrs/proxymeister/pkg/utils"
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
	ip, err := utils.GetMyIP()
	if err != nil {
		log.Fatal(err)
	}

	return ip
}
