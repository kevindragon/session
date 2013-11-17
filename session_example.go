package main

import (
	"fmt"
	"github.com/kevindragon/session"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sess := session.Start(w, r)
		count := sess.Get("count")
		if count == nil {
			sess.Set("count", 1)
		} else {
			sess.Set("count", count.(int)+1)
		}
		fmt.Fprintln(w, sess)
	})
	log.Println("listen at 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
