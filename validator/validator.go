package validator

import (
	"regexp"

    "github.com/parnurzeal/gorequest"
)

var pattern, _ = regexp.Compile(`^[0-9]{1,3}(\.[0-9]{1,3}){3}$`)

func GetMyIp() string {
	_, body, errs := gorequest.New().Get("http://bot.whatismyipaddress.com").End()
	if errs != nil || !pattern.MatchString(body) {
		return ""
	}

	return body
}

func Validate(myIp string, proxy string) bool {
	_, body, errs := gorequest.New().Proxy(proxy).Get("http://bot.whatismyipaddress.com").End()
	if errs != nil || !pattern.MatchString(body) || myIp == body {
		return false
	}

	return true
}

// func main(){
// 	myIp := GetMyIp()
// 	fmt.Println(myIp)
// 	fmt.Println(Validate(myIp, "http://190.147.220.37:8080"))
// }
