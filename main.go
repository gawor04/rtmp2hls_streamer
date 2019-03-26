package main

import (
	"./rtmp2hls/session"
	log "github.com/cihub/seelog"
	"net/http"
)

func main() {
	defer log.Flush()

	/* need to restart ffmpeg after connection lost */
	session.SessionsObserver()

	/* http requests handlers */
	http.HandleFunc("/play/", session.HLSSessionHandler)
	http.HandleFunc("/sessions", session.NewSessionHandler)
	http.HandleFunc("/sessions/", session.DeleteSessionHandler)

	http.ListenAndServe("127.0.0.1:80", nil)
}