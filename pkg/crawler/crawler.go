package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/stefanoschrs/proxymeister/pkg/types"
)

func FetchProxies() (proxies []types.Proxy, err error) {
	sources, err := ioutil.ReadDir(os.Getenv("CRAWLER_SOURCES"))
	if err != nil {
		return
	}

	// TODO: Change to async execution
	for _, source := range sources {
		if !strings.HasSuffix(source.Name(), ".so") {
			continue
		}

		sourcePlugin, pluginErr := plugin.Open(path.Join(os.Getenv("CRAWLER_SOURCES"), source.Name()))
		if pluginErr != nil {
			err = fmt.Errorf("[%s] failed to open plugin. %w\n", source.Name(), pluginErr)
			return
		}

		name, lookupNameErr := sourcePlugin.Lookup("Name")
		if lookupNameErr != nil {
			err = fmt.Errorf("[%s] failed to lookup 'name'. %w\n", source.Name(), lookupNameErr)
			return
		}
		url, lookupUrlErr := sourcePlugin.Lookup("Url")
		if lookupUrlErr != nil {
			err = fmt.Errorf("[%s] failed to lookup 'url'. %w\n", source.Name(), lookupUrlErr)
			return
		}
		fetch, lookupFetchErr := sourcePlugin.Lookup("Fetch")
		if lookupFetchErr != nil {
			err = lookupFetchErr
			err = fmt.Errorf("[%s] failed to lookup 'fetch'. %w\n", source.Name(), lookupFetchErr)
			return
		}

		if os.Getenv("DEBUG") == "true" {
			log.Printf("Fetching %s (%s)\n", name, url)
		}

		sourceProxies, fetchErr := fetch.(func() ([]types.Proxy, error))()
		if fetchErr != nil {
			err = fmt.Errorf("[%s] failed to fetch proxies. %w\n", source.Name(), fetchErr)
			return
		}

		if os.Getenv("DEBUG") == "true" {
			log.Printf("Found %d proxies for %s\n", len(sourceProxies), name)
		}

		proxies = append(proxies, sourceProxies...)
	}

	return
}
