package validator

import (
	"regexp"
	"time"

    "github.com/parnurzeal/gorequest"
)

const REQ_TIMEOUT = time.Second * 15

var pattern, _ = regexp.Compile(`^[0-9]{1,3}(\.[0-9]{1,3}){3}$`)

func GetMyIp() string {
	_, body, errs := gorequest.New().Get("http://bot.whatismyipaddress.com").End()
	if errs != nil || !pattern.MatchString(body) {
		return ""
	}

	return body
}

func Validate(myIp string, proxy string) bool {
	_, body, errs := gorequest.New().Proxy(proxy).Timeout(REQ_TIMEOUT).Get("http://bot.whatismyipaddress.com").End()
	if errs != nil || !pattern.MatchString(body) || myIp == body {
		return false
	}

	return true
}
