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
)

const HashLength uint = 64

type User struct {
	Name       string `json:"username"`
	SecretHash string `json:"secret_hash"`
	Created    uint64
}

// Used to declare connection/disconnection messages in chat
var UserSystem User = User{
	Name:       "System",
	SecretHash: "_",
	Created:    0,
}

// Get existing user from database by unique name
func (db *DB) GetUser(name string) (User, error) {
	command := fmt.Sprintf("SELECT * FROM %s WHERE username=\"%s\"", UsersTablename, name)
	result := db.QueryRow(command)
	err := result.Err()
	if err != nil {
		return User{}, err
	}

	user := User{}
	err = result.Scan(&user.Name, &user.SecretHash)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Create a new user with a unique name
func (db *DB) CreateUser(user *User) error {
	command := fmt.Sprintf("INSERT INTO %s (username, secret_hash) VALUES (?, ?)", UsersTablename)
	_, err := db.Exec(command, user.Name, user.SecretHash)
	if err != nil {
		return err
	}

	return nil
}

// Tells if the user already exists in the db
func (db *DB) DoesUserExist(name string) bool {
	_, err := db.GetUser(name)
	if err != nil {
		return false
	}

	return true
}

// func (db *DB) GetAllUsers() (*[]User, error) {
// 	command := fmt.Sprintf("SELECT * FROM %s", UsersTablename)
// 	rows, err := db.Query(command)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var users []User

// 	for rows.Next() {
// 		var user User
// 		rows.Scan(&user.Name, &user.SecretHash)
// 		users = append(users, user)
// 	}

// 	return &users, nil
// }

// Used to send validation check results for user credentials to frontend
type UserCredentialsValidation struct {
	Valid bool `json:"valid"`
}
