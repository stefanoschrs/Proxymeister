package database

import (
	"errors"
	"log"
	"testing"

	"github.com/stefanoschrs/proxymeister/pkg/types"
	"github.com/stefanoschrs/proxymeister/pkg/utils"
)

type createProxyRes struct {
	Proxy   types.Proxy
	Created bool
	Error   error
}

func TestDB_CreateProxy(t *testing.T) {
	db := getDatabase()

	testCases := []struct {
		Proxy types.Proxy
		IsNew bool
		err   error
	}{
		{
			Proxy: types.Proxy{
				Ip:     "161.202.226.194",
				Port:   80,
				Source: "sslproxies.org",
			},
			IsNew: true,
		},
		{
			Proxy: types.Proxy{
				Ip:     "161.202.226.194",
				Port:   80,
				Source: "sslproxies.org",
			},
			IsNew: false,
		},
		{
			Proxy: types.Proxy{
				Ip:     "x161.202.226.194",
				Port:   80,
				Source: "sslproxies.org",
			},
			err: errors.New(types.ErrInvalidIp),
		},
		{
			Proxy: types.Proxy{
				Ip:     "161.202.226.194",
				Port:   66000,
				Source: "sslproxies.org",
			},
			err: errors.New(types.ErrInvalidPort),
		},
	}

	for _, testCase := range testCases {
		func() {
			var res createProxyRes
			defer createProxy(t, db, testCase.Proxy, &res)()

			if res.Error != nil {
				if testCase.err == nil || res.Error.Error() != testCase.err.Error() {
					t.Fatal(res.Error)
				}

				return
			}
			if !testCase.IsNew && res.Created {
				t.Fatal(errors.New("proxy already in"))
			}
		}()
	}
}

func TestMain(m *testing.M) {
	err := utils.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func getDatabase() (db DB) {
	db, err := Init()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func createProxy(t *testing.T, db DB, p types.Proxy, res *createProxyRes) func() {
	proxy, created, err := db.CreateProxy(p)
	res.Proxy = proxy
	res.Created = created
	res.Error = err

	return func() {
		if err != nil {
			return
		}

		err = db.DeleteProxy(proxy.ID)
		if err != nil {
			t.Fatal(err)
		}
	}
}
