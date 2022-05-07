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
	"time"

	"github.com/gorilla/websocket"
	"unbewohnte.xyz/Unbewohnte/gochat/api"
	"unbewohnte.xyz/Unbewohnte/gochat/log"
)

// Websocket creation handler
func (s *Server) HandlerWebsockets(w http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{}
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		http.Error(w, "could not upgrade to a websocket", http.StatusInternalServerError)
		return
	}

	var authInfo api.User
	err = socket.ReadJSON(&authInfo)
	if err != nil {
		log.Error("%s did not send auth info on websocket connection at all: %s", socket.RemoteAddr(), err)
		api.Disconnect(socket)
		return
	}

	if !s.db.DoesUserExist(authInfo.Name) {
		log.Error("%s sent an invalid authentication data on websocket: user does not exist", socket.RemoteAddr())
		api.Disconnect(socket)
		return
	}

	userDB, _ := s.db.GetUser(authInfo.Name)

	if authInfo.SecretHash != userDB.SecretHash {
		log.Error("%s did not send valid password hash on websocket", socket.RemoteAddr())
		api.Disconnect(socket)
		return
	}

	newWS := api.WS{
		Socket: socket,
		User:   authInfo,
	}

	// add this new authorized socket to the broadcast
	s.websockets.Sockets = append(s.websockets.Sockets, &newWS)

	log.Info("a new websocket connection has been established with %s as \"%s\"", newWS.Socket.RemoteAddr(), newWS.User.Name)
	go s.websockets.HandleNewWebSocketMessages(&newWS, s.incomingMessages)

	// send this new socket all previous messages to display
	allMessages, err := s.db.GetAllMessages()
	if err != nil {
		// well, or not
		log.Error("could not get all previous messages to display for %s: %s", newWS.User.Name, err)
	} else {
		for _, message := range *allMessages {
			newWS.Socket.WriteJSON(&message)
		}
	}

	// notify chat that a new user has connected
	newConnectionMessage := api.Message{
		From:     api.UserSystem,
		Contents: fmt.Sprintf("%s has connected", newWS.User.Name),
	}
	s.incomingMessages <- newConnectionMessage
}

// Send incoming messages from each websocket to all connected ones
func (s *Server) BroadcastMessages() {
	for {
		message, ok := <-s.incomingMessages
		if !ok {
			break
		}

		if message.From.Name == api.UserSystem.Name {
			// add timestapm manually for the system user
			message.TimeStamp = uint64(time.Now().Unix())
		}

		// add incoming message to the db
		err := s.db.AddMessage(message)
		if err != nil {
			log.Error("could not add a new message from %s to the db: %s", message.From.Name, err)
		}

		for _, ws := range s.websockets.Sockets {
			ws.Socket.WriteJSON(&message)
		}
	}
}
