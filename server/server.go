package server

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type HttpServer struct {
	Addr            string
	Handler         http.Handler
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	KeepAliveEnable bool
	RouterFunc      func(map[string]func(http.ResponseWriter, *http.Request))
	Wg              *sync.WaitGroup
	Mux             map[string]func(http.ResponseWriter, *http.Request)
}

func (s *HttpServer) Startup() {
	s.Wg.Add(1)
	go startupHttpServer(s)
}

func (s *HttpServer) Shutdown() {
	shutdownHttpServer(s)
}

func startupHttpServer(hs *HttpServer) {
	defer hs.Wg.Done()

	ss := &http.Server{
		Addr:           hs.Addr,
		Handler:        hs.Handler,
		ReadTimeout:    hs.ReadTimeout * time.Millisecond,
		WriteTimeout:   hs.WriteTimeout * time.Millisecond,
		MaxHeaderBytes: hs.MaxHeaderBytes,
	}
	if hs.KeepAliveEnable {
		ss.SetKeepAlivesEnabled(true)
	} else {
		ss.SetKeepAlivesEnabled(false)
	}

	hs.RouterFunc(hs.Mux)
	err := ss.ListenAndServe()
	if err != nil {
		log.Fatalf("stat server ListenAndServe error: %v", err)
	}
}

func shutdownHttpServer(hs *HttpServer) {
	//TODO: cleanup env
}
