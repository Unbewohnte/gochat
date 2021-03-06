/*
gochat - A dead simple real time webchat.
Copyright (C) 2022  Kasyanov Nikolay Alexeyevich (Unbewohnte)
This file is a part of gochat
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"unbewohnte.xyz/gochat/api"
	"unbewohnte.xyz/gochat/log"
	"unbewohnte.xyz/gochat/page"
)

// Server structure that glues api logic and http/websocket server together
type Server struct {
	workingDir       string
	http             *http.Server
	db               *api.DB
	websockets       *api.WSHolder
	incomingMessages chan api.Message
}

// Create a new configured and ready-to-launch server
func New(workingDir string, dbPath string, port uint) (*Server, error) {
	var server = Server{
		workingDir:       workingDir,
		websockets:       &api.WSHolder{},
		incomingMessages: make(chan api.Message),
	}

	err := api.CreateDB(dbPath)
	if err != nil {
		log.Error("could not create database: %s", err)
		os.Exit(1)
	}

	db, err := api.OpenDB(dbPath)
	if err != nil {
		log.Error("could not open database: %s", err)
		os.Exit(1)
	}
	server.db = db

	// set up routes and handlers
	const (
		pagesDirName   string = "pages"
		staticDirName  string = "static"
		scriptsDirName string = "scripts"
	)

	pagesDirPath := filepath.Join(workingDir, pagesDirName)

	serveMux := http.NewServeMux()

	serveMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(workingDir, staticDirName)))))
	serveMux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir(filepath.Join(workingDir, scriptsDirName)))))

	serveMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/":
			requestedPage, err := page.Get(pagesDirPath, "base.html", "index.html")
			if err != nil {
				log.Error("error getting page on route %s: %s", req.URL.Path, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			requestedPage.ExecuteTemplate(w, "index.html", nil)

		default:
			if strings.HasPrefix(req.URL.Path, api.RouteBase) {
				return
			}

			requestedPage, err := page.Get(pagesDirPath, "base.html", req.URL.Path[1:]+".html")
			if err != nil {
				log.Error("error getting page on route %s: %s", req.URL.Path, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			requestedPage.ExecuteTemplate(w, req.URL.Path[1:]+".html", nil)
		}
	})

	// user api endpoint
	serveMux.HandleFunc(api.RouteUsers, server.HandlerUsers)
	// ws api endpoint
	serveMux.HandleFunc(api.RouteWebsockets, server.HandlerWebsockets)

	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      serveMux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	server.http = &httpServer

	return &server, nil
}

// Fire up the server
func (s *Server) Start() {
	defer s.db.ShutDown()
	go s.BroadcastMessages()

	// fire up a server
	err := s.http.ListenAndServe()
	if err != nil {
		log.Error("FATAL server error: %s", err)
	}
}
