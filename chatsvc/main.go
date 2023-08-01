// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type User struct {
	Name string `json:"name"`
}

type Session struct {
	User string `json:"user"`
}

type Chatspace struct {
	Name     string              `json:"name"`
	Users    map[string]*User    `json:"users"`
	Sessions map[string]*Session `json:"sessions"`
}

var addr = flag.String("addr", ":8443", "http service address")
var useTls = flag.Bool("tls", true, "Use HTTPS")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Exactly one command-line argument expected")
	}
	var chatspaces map[string]*Chatspace
	err := json.Unmarshal([]byte(flag.Arg(0)), &chatspaces)
	if err != nil {
		log.Fatal("Command-line argument is not valid JSON", err)
	}
	for chatspaceName, chatspace := range chatspaces {
		hub := newHub(chatspace)
		go hub.run()
		currentChatspaceName := chatspaceName
		http.HandleFunc(fmt.Sprintf("/%s/login", chatspaceName), func(w http.ResponseWriter, r *http.Request) {
			sessionId := r.URL.Query().Get("sessionId")
			var cookieModifiers string
			if *useTls {
				cookieModifiers = "; HttpOnly; SameSite=Strict; Secure"
			} else {
				cookieModifiers = "; HttpOnly; SameSite=Strict"
			}
			w.Header().Add("Set-Cookie", fmt.Sprintf("sessionId=%s%s", sessionId, cookieModifiers))
			http.Redirect(w, r, fmt.Sprintf("../%s", currentChatspaceName), http.StatusSeeOther)
		})
		http.HandleFunc(fmt.Sprintf("/%s", chatspaceName), serveHome)
		http.HandleFunc(fmt.Sprintf("/%s/ws", chatspaceName), func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
	}
	if *useTls {
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
	} else {
		log.Fatal(http.ListenAndServe(*addr, nil))
	}
}
