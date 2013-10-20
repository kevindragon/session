// Provides session and infrastructure for custom session backends
package session

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// 清理过期的session
func init() {
	go storegc()
}

type session struct {
	sid    string
	m      map[string]interface{}
	create time.Time
	maxage time.Duration
}

type Session interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	SessionId() string
}

func (s *session) Get(key string) interface{} {
	create := s.create
	maxage := time.Duration(s.maxage)
	if create.Add(maxage).Unix() < time.Now().Unix() {
		del(s.sid)
		return nil
	} else {
		if v, ok := s.m[key]; ok {
			return v
		}
	}
	return nil
}

func (s *session) SessionId() string {
	return s.sid
}

func (s *session) Set(key string, value interface{}) {
	s.m[key] = value
}

func newSession(sid string) Session {
	now := time.Now()
	v := make(map[string]interface{})
	session := &session{
		sid:    sid,
		m:      v,
		create: now,
		maxage: time.Second * 70}
	return session
}

func Start(w http.ResponseWriter, r *http.Request) Session {
	var sess Session
	sessionid, err := r.Cookie("sessionid")
	if sessionid == nil || err != nil {
		sid := newSha1()
		sessionid = &http.Cookie{
			Name:   "sessionid",
			Value:  sid,
			Path:   "/",
			MaxAge: 3600}
		http.SetCookie(w, sessionid)
		sess = newSession(sid)
		write(sess)
	} else {
		sid := sessionid.Value
		sess = read(sid)
		if sess == nil {
			sess = newSession(sid)
			write(sess)
		}
	}

	return sess
}

func newSessionId() string {
	ns := strconv.Itoa(time.Now().Nanosecond())
	ss := strconv.Itoa(time.Now().Second())
	t := sha1.New()
	io.WriteString(t, ns+"."+ss)
	sh := fmt.Sprintf("%x", t.Sum(nil))
	return sh
}

func newSha1() string {
	ns := strconv.Itoa(time.Now().Nanosecond())
	ss := strconv.Itoa(time.Now().Second())
	t := sha1.New()
	io.WriteString(t, ns+"."+ss)
	sh := fmt.Sprintf("%x", t.Sum(nil))
	return sh
}
