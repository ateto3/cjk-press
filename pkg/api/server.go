package api

import "net/http"

func StartServer(addr string) error {
	return http.ListenAndServe(addr, nil)
}
