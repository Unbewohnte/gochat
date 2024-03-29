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
	"html/template"
	"path/filepath"
)

// Parse files in pagesDir and return a ready-to-render template
func GetPage(pagesDir string, base string, pagename string) (*template.Template, error) {
	page, err := template.ParseFiles(filepath.Join(pagesDir, base), filepath.Join(pagesDir, pagename))
	if err != nil {
		return nil, err
	}

	return page, nil
}
