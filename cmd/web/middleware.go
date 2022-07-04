package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// func I18n(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		_, err := r.Cookie("lang")
// 		if err != nil {
// 			cookie := &http.Cookie{
// 				Name:  "lang",
// 				Value: "en",
// 				Path:  "/",
// 			}
// 			http.SetCookie(rw, cookie)
// 		}
// 		next.ServeHTTP(rw, r)
// 	})
// }

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
