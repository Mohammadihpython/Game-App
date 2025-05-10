package main

import (
	"GameApp/entity"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"

	"net/http"
)

func main() {
	http.ListenAndServe(
		":8080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				panic(err)
			}

			defer conn.Close()

			//	we create an chanel to wait for go routines
			done := make(chan bool)
			go readMessage(conn, done)

		}),
	)
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {

			log.Println(err)
			done <- true
			return
		}
		var notif entity.Notification
		err = json.Unmarshal(msg, &notif)
		if err != nil {
			panic(err)
		}
		fmt.Println(notif)
		fmt.Println(opCode)

	}
}
