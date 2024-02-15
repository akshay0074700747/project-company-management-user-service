package helpers

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func DialGrpc(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}

func GenUuid() string {
	return uuid.New().String()
}

func PrintErr(err error, messge string) {
	fmt.Println(messge, err)
}

func PrintMsg(msg string)  {
	fmt.Println(msg)
}

func IsValidEmail(email string) (bool,error) {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(emailRegex, email)

	return match,err
}

func IsValidPhoneNumber(phoneNumber string) bool {

	pattern := `^\+?\d{1,3}[-\s]?\d{3}[-\s]?\d{3}[-\s]?\d{4}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(phoneNumber)
}
