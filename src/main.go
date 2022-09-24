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

package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"unbewohnte/gochat/log"
	"unbewohnte/gochat/server"
)

const version string = "0.1.1"

var (
	port        *uint   = flag.Uint("port", 8080, "Set working port")
	tlsKeyFile  *string = flag.String("tlsKeyFile", "", "Specify tls key file")
	tlsCertFile *string = flag.String("tlsCertFile", "", "Specify tls cert file")
)

func main() {
	// set up logging
	exePath, err := os.Executable()
	if err != nil {
		log.Error("could not lookup executable's path: %s", err)
	}

	exeDirPath := filepath.Dir(exePath)
	logsDirPath := filepath.Join(exeDirPath, "logs")

	err = os.MkdirAll(logsDirPath, os.ModePerm)
	if err != nil {
		log.Error("could not create logs directory: %s", err)
	}

	logsFile, err := os.Create(filepath.Join(logsDirPath, "latest.log"))
	if err != nil {
		log.Error("could not create logs file: %s", err)
	}

	if err == nil {
		log.SetOutput(io.MultiWriter(os.Stdout, logsFile))
	} else {
		log.SetOutput(os.Stdout)
	}

	// parse flags
	flag.Parse()

	const dbFilename string = "gochat.db"
	dbPath := filepath.Join(exeDirPath, dbFilename)

	server, err := server.New(exeDirPath, dbPath, *port, *tlsKeyFile, *tlsCertFile)
	if err != nil {
		log.Error("could not create a new server instance: %s", err)
	}
	server.Start()
}
