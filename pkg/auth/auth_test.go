package auth

import (
	"net/http"
	"testing"
)

func BenchmarkCreateIDToCookie(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CreateIDToCookie()
		}
	})
}

func BenchmarkGetIDFromCookie(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		cookieValue := "30316566363932392d373133392d363036632d393466312d303031353564336462623865da1bc3dd8b597c82ec439df707f6fc8e988ef19082aac7950957830c1bb4aa2a"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			GetIDFromCookie(cookieValue)
		}
	})
}

func BenchmarkValidateCookie(b *testing.B) {
	b.Run("default", func(b *testing.B) {
		cookie := &http.Cookie{
			Name:  "user_id",
			Value: "30316566363932392d373133392d363036632d393466312d303031353564336462623865da1bc3dd8b597c82ec439df707f6fc8e988ef19082aac7950957830c1bb4aa2a",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ValidateCookie(cookie)
		}
	})
}
