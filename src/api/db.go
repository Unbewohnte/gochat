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
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	UsersTablename    string = "users"
	MessagesTablename string = "messages"
)

// SQL database wrapper
type DB struct {
	*sql.DB
}

// Closes database connection
func (db *DB) ShutDown() {
	db.Close()
}

func (db *DB) createUsersTable() error {
	command := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s 
		(username TEXT NOT NULL PRIMARY KEY, secret_hash TEXT NOT NULL)`, UsersTablename)
	_, err := db.Exec(command)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) createMessagesTable() error {
	command := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s
		(id INTEGER NOT NULL PRIMARY KEY, content TEXT NOT NULL,
		sender TEXT NOT NULL, timestamp INTEGER, FOREIGN KEY(sender) REFERENCES %s(username))`,
		MessagesTablename, UsersTablename,
	)

	_, err := db.Exec(command)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) setUpTables() error {
	err := db.createUsersTable()
	if err != nil {
		return fmt.Errorf("error creating users table: %s", err)
	}

	err = db.createMessagesTable()
	if err != nil {
		return fmt.Errorf("error creating messages table: %s", err)
	}

	return nil
}

// Creates a database if does not exist
func CreateDB(path string) error {
	dbFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	dbFile.Close()

	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	db := DB{
		sqlDB,
	}

	err = db.setUpTables()
	if err != nil {
		return fmt.Errorf("could not set up tables: %s", err)
	}

	return nil
}

// Opens a DB
func OpenDB(path string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db := DB{
		sqlDB,
	}

	return &db, nil
}
