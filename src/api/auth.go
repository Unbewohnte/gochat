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
	"net/http"
	"strings"
)

const (
	AuthHeaderKey string = "AUTH_INFO"
	AuthSeparator string = ":"
)

var ErrNotAuthorized error = fmt.Errorf("not authorized")

// Extract information about the sender from request headers.
// If no data is given or just wrong - returns a default user struct
func GetUserAuthHeaderData(req *http.Request) User {
	authData := req.Header.Get(AuthHeaderKey)
	dividedAuthData := strings.Split(authData, AuthSeparator)
	if len(dividedAuthData) != 2 {
		return User{}
	}
	username, passHash := dividedAuthData[0], dividedAuthData[1]

	return User{
		Name:       username,
		SecretHash: passHash,
	}
}

// Checks whether provided request indeed contains SOME authentication information.
// Does not check it for being legitimate or not
func WithAuth(req *http.Request) bool {
	authenticationData := req.Header.Get(AuthHeaderKey)
	if authenticationData == "" || len(authenticationData) < 3 {
		return false
	}

	return true
}

// Checks whether request contains valid and real user authorization data that is in the database
func (db *DB) IsValidAuth(req *http.Request) bool {
	if !WithAuth(req) {
		return false
	}

	headerUser := GetUserAuthHeaderData(req)

	if !db.DoesUserExist(headerUser.Name) {
		return false
	}

	userDB, err := db.GetUser(headerUser.Name)
	if err != nil {
		return false
	}

	if userDB.SecretHash != headerUser.SecretHash {
		return false
	}

	return true
}
