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

// A message struct that is converted back and forth to JSON in order to communicate with frontend
type Message struct {
	TimeStamp uint64 `json:"timestamp"`
	From      User   `json:"from"`
	Contents  string `json:"contents"`
}

/* PROBABLY come back later and implement a message-remembering feature */

// func (db *DB) AddMessage(message Message) error {
// 	command := fmt.Sprintf("INSERT INTO %s(from, contents) VALUES(%s, %s)", MessagesTablename, message.From.Name, message.Contents)
// 	_, err := db.Exec(command)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *DB) GetAllMessages() (*[]Message, error) {
// 	command := fmt.Sprintf("SELECT * FROM %s", MessagesTablename)
// 	rows, err := db.Query(command)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var messages []Message

// 	/*
// 		(id INTEGER NOT NULL PRIMARY KEY, content TEXT NOT NULL,
// 		 to_name TEXT NOT NULL, from_name TEXT NOT NULL, FOREIGN KEY(from_name, to_name) REFERENCES %s(username, username))`,
// 	*/
// 	for rows.Next() {
// 		var message Message
// 		rows.Scan(&message.ID, &message.Contents, message.From.Name)
// 		messages = append(messages, message)
// 	}

// 	return &messages, nil
// }
