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

package api

import (
	"encoding/json"
	"fmt"
	"time"

	"unbewohnte/gochat/log"

	"github.com/gorilla/websocket"
)

// Websocket connection paired with the owner
type WS struct {
	Socket *websocket.Conn
	User   User
}

// All existing websockets wrapper
type WSHolder struct {
	Sockets []*WS
}

// Untrack user's websocket from broadcast
func (holder *WSHolder) Remove(username string) {
	var indexToRemove int
	for index, socket := range holder.Sockets {
		if socket.User.Name == username {
			indexToRemove = index
			break
		}
	}

	if len(holder.Sockets) >= 2 {
		holder.Sockets = append(holder.Sockets[:indexToRemove], holder.Sockets[indexToRemove+1:]...)
	} else {
		holder.Sockets = append(holder.Sockets[:indexToRemove])
	}
}

// Tell the websocket that it is no longer connected
func Disconnect(ws *websocket.Conn) {
	ws.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second*1))
	ws.Close()
}

// Handle all incoming messages from this websocket, redirecting messages to specified output
func (holder *WSHolder) HandleNewWebSocketMessages(ws *WS, output chan Message) {
	for {
		msgType, contents, err := ws.Socket.ReadMessage()
		if err != nil {
			log.Error("could not read message from websocket: %s", err)
			Disconnect(ws.Socket)
			holder.Remove(ws.User.Name)

			disconnectionMessage := Message{
				TimeStamp: uint64(time.Now().UnixMilli()),
				From:      UserSystem,
				Contents:  fmt.Sprintf("%s has disconnected", ws.User.Name),
			}
			output <- disconnectionMessage
			break
		}

		switch msgType {
		case websocket.CloseMessage:
			ws.Socket.Close()
			holder.Remove(ws.User.Name)

			disconnectionMessage := Message{
				TimeStamp: uint64(time.Now().UnixMilli()),
				From:      UserSystem,
				Contents:  fmt.Sprintf("%s has disconnected", ws.User.Name),
			}
			output <- disconnectionMessage

		case websocket.TextMessage:
			var newMessage Message
			err = json.Unmarshal(contents, &newMessage)
			if err != nil {
				log.Error("error unmarshaling text message from \"%s\"'s websocket: %s", ws.User.Name, err)
				break
			}

			if len(newMessage.Contents) < 1 || uint(len(newMessage.Contents)) > MaxMessageContentLen {
				break
			}

			newMessage.From = ws.User
			newMessage.TimeStamp = uint64(time.Now().UnixMilli())
			output <- newMessage
		}

	}
}
