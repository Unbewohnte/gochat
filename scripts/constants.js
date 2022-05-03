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

const API_BASE_ENDPOINT = "api";
const API_USER_ENDPOINT = "api/user"
const API_WEBSOCKET_ENDPOINT = "api/ws";

const AUTH_HEADER_KEY = "AUTH_INFO";
const AUTH_SEPARATOR = ":"

const LOCALSTORAGE_NAME_KEY = "username";
const LOCALSTORAGE_PASSWORD_HASH_KEY = "secret_hash";

const BACKEND_URL = window.location.host;