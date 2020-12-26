package crawler

import (
	"fmt"
	"testing"
)

func TestFetchProxies(t *testing.T) {
	proxies, err := FetchProxies()
	if err != nil {
		t.Fatal(err)
	}

	for _, proxy := range proxies {
		fmt.Printf("%+v\n", proxy)
	}
}
