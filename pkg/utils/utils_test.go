package utils

import (
	"github.com/stefanoschrs/proxymeister/pkg/types"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestHttpProxyfiedGet(t *testing.T) {
	_ = os.Setenv("PROXYMEISTER_URL", "http://localhost:8080")

	body, err := HttpProxyfiedGet(types.ProxifiedGetOptions{TargetUrl: "https://www.google.com"})
	//body, err := HttpProxyfiedGet(types.ProxifiedGetOptions{TargetUrl: "https://ipinfo.io"})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("DEBUG") == "true" {
		log.Println(string(body))
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
