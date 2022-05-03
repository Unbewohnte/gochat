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
	"fmt"
	"io"
	"os"
	"path/filepath"

	"unbewohnte.xyz/gochat/log"
	"unbewohnte.xyz/gochat/server"
)

const version string = "0.1.0"

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

	// work out launch flags
	var port uint
	flag.UintVar(&port, "port", 8080, "Set working port")
	flag.Usage = func() {
		fmt.Printf("gochat v%s\n\nFlags\nport [uint] -> specify a port number (default: 8080)\n\n(c) Unbewohnte (Kasyanov Nikolay Alexeyevich)\n", version)
	}
	flag.Parse()

	const dbFilename string = "gochat.db"
	dbPath := filepath.Join(exeDirPath, dbFilename)

	server, err := server.New(exeDirPath, dbPath, port)
	if err != nil {
		log.Error("could not create a new server instance: %s", err)
	}
	server.Start()
}
