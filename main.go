package main

import (
	"./rtmp2hls/session"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	
	defer log.Flush()

	args := os.Args[1:]

	if len(args) != 2 {

		log.Critical("Wrong number of arguments! ip_address and port are required")
	}

	ip   := args[0]
	port := args[1]

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	/* clear filesystem on exit */
	go func() {
		sig := <-sigs
		log.Critical(sig)

		session.SessionsClear()
		os.Exit(0)
	}()

	/* need to restart ffmpeg after connection lost */
	session.SessionsObserver()

	/* http requests handlers */
	http.HandleFunc("/play/", session.HLSSessionHandler)
	http.HandleFunc("/sessions", session.NewSessionHandler)
	http.HandleFunc("/sessions/", session.DeleteSessionHandler)

	http.ListenAndServe( ip + ":" + port, nil)
}
