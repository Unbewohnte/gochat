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

	"unbewohnte/gochat/api"
	"unbewohnte/gochat/log"
)

// Server structure that glues api logic and http/websocket server together
type Server struct {
	workingDir       string
	keyFile          string
	certFile         string
	http             *http.Server
	db               *api.DB
	websockets       *api.WSHolder
	incomingMessages chan api.Message
}

const (
	pagesDirName       string = "pages"
	staticDirName      string = "static"
	scriptsDirName     string = "scripts"
	attachmentsDirName string = "attachments"
)

// Create a new configured and ready-to-launch server
func New(workingDir string, dbPath string, port uint, keyFile string, certFile string) (*Server, error) {
	var server = Server{
		keyFile:          keyFile,
		certFile:         certFile,
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
	attachmentsDirPath := filepath.Join(workingDir, attachmentsDirName)
	err = os.MkdirAll(attachmentsDirPath, os.ModePerm)
	if err != nil {
		log.Error("could not create attachments directory: %s", err)
		os.Exit(1)
	}

	pagesDirPath := filepath.Join(workingDir, pagesDirName)

	serveMux := http.NewServeMux()

	serveMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(workingDir, staticDirName)))))
	serveMux.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir(filepath.Join(workingDir, scriptsDirName)))))
	serveMux.Handle("/attachments/", http.StripPrefix("/attachments/", http.FileServer(http.Dir(attachmentsDirPath))))

	serveMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/":
			requestedPage, err := GetPage(pagesDirPath, "base.html", "index.html")
			if err != nil {
				log.Error("error getting page on route %s: %s", req.URL.Path, err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			requestedPage.ExecuteTemplate(w, "index.html", nil)

		default:
			if strings.HasPrefix(req.URL.Path, api.RouteBase) {
				return
			} else if strings.Contains(req.URL.Path, "favicon.ico") {
				// remove that annoying favicon error by simply ingoring the thing
				return
			}

			requestedPage, err := GetPage(pagesDirPath, "base.html", req.URL.Path[1:]+".html")
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
	// attachments handler
	serveMux.HandleFunc(api.RoutePostAttachemnts, server.handleAttachments)

	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      serveMux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	server.http = &httpServer

	log.Info("Created server instance")

	return &server, nil
}

// Fire up the server
func (s *Server) Start() {
	defer s.db.ShutDown()
	// broadcast messages
	go s.BroadcastMessages()
	// clean attachments storage from time to time
	// max attachment filesize * 50 is the limit, check every 5 sec
	go manageAttachmentsStorage(filepath.Join(s.workingDir, attachmentsDirName), api.MaxAttachmentSize*50, time.Second*5)

	// fire up either a TLS or non-TLS server
	if s.keyFile != "" && s.certFile != "" {
		log.Info("Using TLS")
		log.Info("Working on %s", s.http.Addr)

		err := s.http.ListenAndServeTLS(s.certFile, s.keyFile)
		if err != nil {
			log.Error("Fatal server error: %s", err)
		}
	} else {
		log.Info("Not using TLS")
		log.Info("Working on %s", s.http.Addr)

		err := s.http.ListenAndServe()
		if err != nil {
			log.Error("Fatal server error: %s", err)
		}
	}
}
