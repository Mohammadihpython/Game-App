package main

import (
	"github.com/gobwas/ws"
	"net/http"
)

func main() {
	http.ListenAndServe(
		":8080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				return
			}
			defer conn.Close()

		}),
	)
}
