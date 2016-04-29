package validator

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	myIp := GetMyIp()
	fmt.Println(Validate(myIp, "http://190.74.213.185:8088"))
	fmt.Println(Validate(myIp, "http://165.139.149.169:3128"))
}
