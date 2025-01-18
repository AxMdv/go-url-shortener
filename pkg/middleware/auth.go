// Package middleware provides handling http requests and responses.
// There are authentication, gzip and logging middlewares.
package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/AxMdv/go-url-shortener/pkg/auth"
)

const cookieName = "user_id"

// SignUpMiddleware sets cookie to http.ResponseWriter if http.Request doesn`t have necessary cookie or cookie is invalid.
// Also it sets user UUID to context of a http.Request.
func SignUpMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(cookieName)
		// если куки нет
		if err != nil {
			id, cookieValue, createErr := auth.CreateIDToCookie()
			if createErr != nil {
				log.Println("can`t create id to cookie", createErr)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: cookieName, Value: cookieValue})
			cr := auth.SetUUIDToRequestContext(r, id)
			h.ServeHTTP(w, cr)
			return
		}

		// кука есть: проверяем айди на валидность
		valid, err := auth.ValidateCookie(cookie)

		if !valid || err != nil {
			// кука не валидна
			log.Println("cookie is not valid or error during validation of cookie", err)

			id, cookieValue, err := auth.CreateIDToCookie()
			if err != nil {
				log.Println("can`t create id to cookie", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{Name: cookieName, Value: cookieValue})
			cr := auth.SetUUIDToRequestContext(r, id)
			h.ServeHTTP(w, cr)
			return //
		}
		// c айди все ок - передаём в контексте реквеста айди
		id := auth.GetIDFromCookie(cookie.Value)
		cr := auth.SetUUIDToRequestContext(r, id)
		h.ServeHTTP(w, cr)
	}
}

// ValidateUserMiddleware validates http.Cookie of a http.Request.
// If cookie is valid user UUID is set to request context.
// If cookie is invalid this function returns http.StatusUnauthorized.
func ValidateUserMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		valid, err := auth.ValidateCookie(cookie)
		if !valid || err != nil {
			// кука не валидна
			log.Println("cookie is not valid or error during validation of cookie", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// c айди все ок - передаём в контексте реквеста айди
		id := auth.GetIDFromCookie(cookie.Value)
		cr := auth.SetUUIDToRequestContext(r, id)
		h.ServeHTTP(w, cr)
	}
}

// TrustedSubnet checks request IP for being in trusted subnet.
func TrustedSubnet(h http.HandlerFunc, subnet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if subnet == "" {
			log.Println("s.Config.TrustedSubnet == empty")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ip, err := auth.GetIPfromRequest(r)
		if err != nil {
			log.Println("getIPfromRequest(r)", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			log.Println("fail to parse cidr", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		fmt.Printf("ipnet: %v, ip: %v\n", ipNet, ip.To16().String())
		trustedRequest := ipNet.Contains(ip)
		if !trustedRequest {
			log.Println("untrusted trustedRequest")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
}
