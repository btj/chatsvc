// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8443", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	tlsCert, err := base64.StdEncoding.DecodeString(os.Getenv("CHATSVC_TLS_CERT"))
	if err != nil {
		log.Fatal("CHATSVC_TLS_CERT: not valid base64", err)
	}
	tlsPrivateKey, err := base64.StdEncoding.DecodeString(os.Getenv("CHATSVC_TLS_PRIVATE_KEY"))
	if err != nil {
		log.Fatal("CHATSVC_TLS_PRIVATE_KEY: not valid base64", err)
	}
	cert, err := tls.X509KeyPair(tlsCert, tlsPrivateKey)
	if err != nil {
		log.Fatal("CHATSVC_TLS_CERT or CHATSVC_TLS_PRIVATE_KEY: Bad cert or key", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := http.Server{
		Addr:      *addr,
		TLSConfig: tlsConfig,
	}
	log.Fatal(server.ListenAndServeTLS("", ""))
}
