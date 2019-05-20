package logger

import (
	"net/http"
)

//decorator function
func Log(h http.Handler, url string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//fmt.Println(time.Now(), url)
		h.ServeHTTP(res, req)
	})
}
