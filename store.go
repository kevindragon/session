package session

import (
	"fmt"
	"time"
)

var store = make(map[string]Session)

func write(sess Session) {
	store[sess.SessionId()] = sess
}

func read(sid string) Session {
	return store[sid]
}

func del(sid string) {
	delete(store, sid)
}

func storegc() {
	fmt.Println("start session gc")
	var needRemove []string
	for k, _ := range store {
		sess, _ := store[k].(*session)
		create := sess.create
		maxage := time.Duration(sess.maxage)
		if create.Add(maxage).Unix() < time.Now().Unix() {
			needRemove = append(needRemove, k)
		}
	}
	for _, k := range needRemove {
		delete(store, k)
	}
	fmt.Println("session gc end")
	time.AfterFunc(time.Second*3600, storegc)
}
