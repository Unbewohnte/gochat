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
	"fmt"
	"strconv"
)

// A message struct that is converted back and forth to JSON in order to communicate with frontend
type Message struct {
	ID        uint64
	TimeStamp uint64 `json:"timestamp"`
	From      User   `json:"from"`
	Contents  string `json:"contents"`
}

// Add a new message to the database, following the limitations and removing the
// oldest one as well if the capacity has been exceeded
func (db *DB) AddMessage(message Message) error {
	// check how many messages are already stored
	command := fmt.Sprintf("SELECT COUNT(*) FROM %s", MessagesTablename)
	result := db.QueryRow(command)
	err := result.Err()
	if err != nil {
		return err
	}

	var countStr string
	result.Scan(&countStr)
	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		return err
	}

	if count >= uint64(MaxMessagesRemembered) {
		// remove the last one
		command = fmt.Sprintf("DELETE FROM %s WHERE timestamp = (SELECT MIN(timestamp) FROM %s)", MessagesTablename, MessagesTablename)
		_, err := db.Exec(command)
		if err != nil {
			return err
		}
	}

	command = fmt.Sprintf("INSERT INTO %s(sender, content, timestamp) VALUES('%s', '%s', %d)",
		MessagesTablename, message.From.Name, message.Contents, message.TimeStamp)
	_, err = db.Exec(command)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAllMessages() (*[]Message, error) {
	command := fmt.Sprintf("SELECT * FROM %s", MessagesTablename)
	rows, err := db.Query(command)
	if err != nil {
		return nil, err
	}

	var messages []Message

	for rows.Next() {
		var message Message
		rows.Scan(&message.ID, &message.Contents, &message.From.Name, &message.TimeStamp)
		messages = append(messages, message)
	}

	return &messages, nil
}
