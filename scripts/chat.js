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

let username = localStorage.getItem(LOCALSTORAGE_NAME_KEY);
let secret_hash = localStorage.getItem(LOCALSTORAGE_PASSWORD_HASH_KEY);

// check whether credentials are right or not
areCredentialsValid({username:username, secret_hash:secret_hash})
.then(valid => {
    if (valid === true) {
        // they are valid, connect as a websocket
        document.getElementById("logged_username").innerHTML = username;
        let socket = new WebSocket(`ws://${BACKEND_URL}/${API_WEBSOCKET_ENDPOINT}`);
    
        // send auth data right away, as API tells us to do
        socket.onopen = event => {
            console.log("Connected to the server");
            let auth_data = JSON.stringify({
                username: username,
                secret_hash: secret_hash,
            });
            socket.send(auth_data);
        }
    
        socket.onclose = event => {
            console.log("Closed connection: ", event);
            socket.close();

            let chatbox = document.getElementById("chatbox");
            chatbox.innerHTML += "Connection closed";
            chatbox.scrollTop = chatbox.scrollHeight;
        };
    
        // display it
        socket.onmessage = event => {
            let message = JSON.parse(event.data); 
            let from_user = message.from.username;
            let date = new Date(message.timestamp).toLocaleString();
            let constructed_message_to_display = `[${date}] ${from_user}: ${message.contents}` + "\n";

            let chatbox = document.getElementById("chatbox");
            chatbox.innerHTML += constructed_message_to_display;
            chatbox.scrollTop = chatbox.scrollHeight;
        }

        // make buttons do intended
        let send_button = document.getElementById("send_button");
        send_button.addEventListener("click", (event) => {
            let text_input = document.getElementById("text_input");

            socket.send(
                JSON.stringify({
                    contents: String(text_input.value),
                }));

            text_input.value = "";
            text_input.focus();
        });

        let disconnect_button = document.getElementById("disconnect_button");
        disconnect_button.addEventListener("click", (event) => {
            socket.close();
            disconnect_button.value = "Reconnect";
            disconnect_button.addEventListener("click", (event) => {
               location.reload();
            })
        })

        // make message to be sent via pressing 'Enter' key
        document.addEventListener("keypress", (event) => {
            if (event.key == "Enter") {
                send_button.click();               
            }
        })

        // and focus the text input so the user does not have to click it manually
        document.getElementById("text_input").focus();

    } else {
        // credentials are not valid ! To the registration page you go
        window.location.replace("/register");
    }
})