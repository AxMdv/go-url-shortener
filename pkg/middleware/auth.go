package middleware

import (
	"log"
	"net/http"

	"github.com/AxMdv/go-url-shortener/pkg/auth"
)

const cookieName = "user_id"

func SignUpMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(cookieName)
		//если куки нет
		if err != nil {
			id, cookieValue, err := auth.CreateIDToCookie()
			if err != nil {
				log.Println("can`t create id to cookie", err)
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
		}
		// c айди все ок - передаём в контексте реквеста айди
		id := auth.GetIDFromCookie(cookie.Value)
		cr := auth.SetUUIDToRequestContext(r, id)
		h.ServeHTTP(w, cr)
	}
}
