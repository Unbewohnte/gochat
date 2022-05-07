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
	"encoding/json"
	"io"
	"net/http"

	"unbewohnte.xyz/Unbewohnte/gochat/api"
	"unbewohnte.xyz/Unbewohnte/gochat/log"
)

// User creation/credentials validation http handler
func (s *Server) HandlerUsers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// unmarshal incoming json
		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		var newUser api.User
		err = json.Unmarshal(requestBody, &newUser)
		if err != nil {
			log.Error("could not unmarshal new user json: %s", err)
			http.Error(w, "invalid user json", http.StatusBadRequest)
			return
		}

		// check if this all is genuine
		if len(newUser.Name) <= 1 || uint(len(newUser.SecretHash)) != api.HashLength {
			log.Error("could not create a new user: incoming name or secret hash is invalid")
			http.Error(w, "invalid user json", http.StatusBadRequest)
			return
		}

		if s.db.DoesUserExist(newUser.Name) {
			// already have !
			http.Error(w, "already exists", http.StatusBadRequest)
			return
		}

		// add user to db
		err = s.db.CreateUser(&newUser)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			log.Error("could not create a new user \"%s\": %s", newUser.Name, err)
			return
		}

		log.Info("successfully created a new user: \"%s\"", newUser.Name)

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		// TODO(redo to make it return true or false whether user credentials are okay or not)
		headerUserData := api.GetUserAuthHeaderData(req)

		if s.db.DoesUserExist(headerUserData.Name) {
			response := api.UserCredentialsValidation{}

			userDB, _ := s.db.GetUser(headerUserData.Name)

			if userDB.SecretHash != headerUserData.SecretHash {
				// passwords do not match
				response.Valid = false
			} else {
				response.Valid = true
			}

			responseBytes, err := json.Marshal(&response)
			if err != nil {
				log.Error("could not marshal user credentials validation response: %s", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			w.Write(responseBytes)

		} else {
			// does not even exist
			responseBytes, err := json.Marshal(&api.UserCredentialsValidation{
				Valid: false,
			})
			if err != nil {
				log.Error("could not marshal user credentials validation response: %s", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			w.Write(responseBytes)
		}

		// users, err := db.getAllUsers()
		// if err != nil {
		// 	http.Error(w, "", http.StatusInternalServerError)
		// 	return
		// }

		// w.Header().Add("Content-type", "application/json")
		// userBytes, err := json.Marshal(*users)
		// if err != nil {
		// 	break
		// }

		// w.Write(userBytes)

	default:
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
	}
}
