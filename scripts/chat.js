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


let chatbox = document.getElementById("chatbox");
let text_input = document.getElementById("text_input");
let error_output = document.getElementById("error_output");
let send_button = document.getElementById("send_button");
let disconnect_button = document.getElementById("disconnect_button");


// check whether credentials are right or not
areCredentialsValid({username:username, secret_hash:secret_hash})
.then(valid => {
    if (valid === false) {
        // credentials are not valid ! To the registration page you go
        window.location.replace("/register");
    } else {
        // they are valid, connect as a websocket
        document.getElementById("logged_username").innerHTML = username;
        let socket = new WebSocket(`ws://${HOST}/${API_WEBSOCKET_ENDPOINT}`);
    
        // send auth data right away, as API tells us to do
        socket.onopen = (_) => {
            console.log("Connected to the server");
            let auth_data = JSON.stringify({
                username: username,
                secret_hash: secret_hash,
            });
            socket.send(auth_data);
        }
    
        // notify the user if the socket has been closed
        socket.onclose = (event) => {
            console.log("Closed connection: ", event);
            socket.close();

            chatbox.innerHTML += "Connection closed";
            chatbox.scrollTop = chatbox.scrollHeight;
        };
    
        // display messages
        socket.onmessage = event => {
            let message = JSON.parse(event.data); 
            let from_user = message.from.username;
            let date = new Date(message.timestamp).toLocaleString();

            let constructed_message_to_display;

            // if message starts with ">" - display greentext
            if (String(message.contents).startsWith(">")) {
                constructed_message_to_display = 
                `<p><small>[${date}]</small> <b>${from_user}</b>: <i class="greentext">${message.contents}</i></p>`;
            } else {
                constructed_message_to_display = 
                `<p><small>[${date}]</small> <b>${from_user}</b>: ${message.contents}</p>`;
            }

            chatbox.innerHTML += constructed_message_to_display;
            chatbox.scrollTop = chatbox.scrollHeight;
        }

        // make buttons do intended
        send_button.addEventListener("click", (event) => {
            // clear previous error message
            error_output.innerHTML = "";

            let file_input = document.getElementById("file_input");
            let file = file_input.files[0];
            // no attachment
            if (file == undefined || file == null) {
                // just send a usual message
                socket.send(
                    JSON.stringify({
                        contents: String(text_input.value),
                    }));
    
                text_input.value = "";
                text_input.focus();
                return;
            }

            // file "probably exists and is real"
            if (file.size >= 1 && file.name.length >= 1) {
                // try to send it as attachment 
                let formdata = new FormData();
                formdata.set(ATTACHMENT_FORM_KEY, file);
                
                let response_status;
                fetch(API_ATTACHMENT_ENDPOINT, {
                    method: "POST",
                    headers: {"AUTH_INFO": username+":"+secret_hash},
                    body: formdata
                })
                .then(response => {
                    response_status = response.status;
                    return response.text();
                })
                .then(response_text => {
                    if (response_status != 200) {
                        error_output.innerHTML = "Error uploading file: " + response_text;
                    } else {
                        // no errors, attachment is on the server, append link to it and send message
                        let response_json = JSON.parse(response_text);
                        
                        let attachment_url;
                        if (is_image(file.name) === true) {
                            // it's an image so let's display it right away
                            attachment_url = 
                            `<a href="${response_json.url}" target="_blank" rel="noopener noreferrer"><img class="chat-image" src="${response_json.url}"></a>`;
                        } else {
                            // just link to it
                            attachment_url =
                            `<a href="${response_json.url}" target="_blank" rel="noopener noreferrer">${file.name}</a>`;
                        }

                        
                        socket.send(
                            JSON.stringify({
                                contents: `${text_input.value} ${attachment_url}`,
                            }));
            
                        text_input.value = "";
                        file_input.value = null;
                        text_input.focus();
                    }
                })
                .catch(error => {
                    error_output.innerHTML = "Error sending attachment: " + error;
                    file_input.value = null;
                });
            }
        });

        // disconnect/reconnect button
        disconnect_button.addEventListener("click", (_) => {
            socket.close();
            disconnect_button.value = "Reconnect";
            disconnect_button.addEventListener("click", (_) => {
            location.reload();
            })
        });

        // make message to be sent via pressing 'Enter' key
        document.addEventListener("keypress", (event) => {
            if (event.key == "Enter") {
                send_button.click();               
            }
        });

        // and focus the text input so the user does not have to click it manually
        document.getElementById("text_input").focus();
    }
});