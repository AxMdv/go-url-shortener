package auth_test

import (
	"fmt"
	"net/http"

	"github.com/AxMdv/go-url-shortener/pkg/auth"
)

const cookieName = "user_id"

func ExampleValidateCookie() {

	cookieValue := "30316566376539342d346130332d363633312d616163312d3132366664393030643461313d236637f1eb14e449f64a77d044f63a120fd5ee06d57c4f06684ae32c36e344"

	cookie := &http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	}
	// Trying to validate cookie
	valid, err := auth.ValidateCookie(cookie)
	if err != nil {
		panic(err)
	}
	if valid {
		fmt.Println("Cookie is valid")
	} else {
		fmt.Println("Cookie is not valid")
	}

	// Output:
	// Cookie is valid
}

func ExampleCreateIDToCookie() {
	id, cookieValue, err := auth.CreateIDToCookie()
	if err != nil {
		panic(err)
	}
	fmt.Printf("id:%s, cookie value:%s\n", id, cookieValue)
}

func ExampleGetIDFromCookie() {
	cookieValue := "30316566376539342d346130332d363633312d616163312d3132366664393030643461313d236637f1eb14e449f64a77d044f63a120fd5ee06d57c4f06684ae32c36e344"
	id := auth.GetIDFromCookie(cookieValue)
	fmt.Printf("ID is %s\n", id)

	// Output:
	// ID is 01ef7e94-4a03-6631-aac1-126fd900d4a1
}
