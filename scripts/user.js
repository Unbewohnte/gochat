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


async function register() {
    let auth_form = document.forms["auth_form"];

    let username = String(auth_form.elements["username"].value).trim();
    if (username.length < 3) {
        document.getElementById("error_output").innerHTML = "Username must be >= 3 characters long";
        return;
    }

    let password = String(auth_form.elements["password"].value).trim();
    if (password.length < 3) {
        document.getElementById("error_output").innerHTML = "Password must be >= 3 characters long";
        return;
    }

    sha256(password).then((password_sha256) => {
        let post_data = {
            username: username,
            secret_hash: password_sha256,
        };
    
        fetch(API_USER_ENDPOINT, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(post_data),
        })
        .then(response => {
            if (response.status == 200) {
                localStorage.setItem(LOCALSTORAGE_NAME_KEY, username);
                localStorage.setItem(LOCALSTORAGE_PASSWORD_HASH_KEY, password_sha256); 
                window.location.replace("/")
            }
        })
        .catch((error) => {
            document.getElementById("error_output").innerHTML = error;
        });
    });
};

async function areCredentialsValid(user_credentials) {    
    let valid = false;

    await fetch(API_USER_ENDPOINT, {
        method: "GET",
        headers: {"AUTH_INFO": user_credentials.username+":"+user_credentials.secret_hash},
    })
    .then(response => {
        return response.text();
    })
    .then(response_data => {
        let response = JSON.parse(response_data)
        console.log(response);
        if (response.valid === true) {
            valid = true;
        } else {
            valid = false;
        }
    });

    return valid;
}

async function login() {
    let auth_form = document.forms["auth_form"];

    let username = String(auth_form.elements["username"].value).trim();
    if (username.length < 3) {
        document.getElementById("error_output").innerHTML = "Username must be >= 3 characters long";
        return;
    }

    let password = String(auth_form.elements["password"].value).trim();
    if (password.length < 3) {
        document.getElementById("error_output").innerHTML = "Password must be >= 3 characters long";
        return;
    }

    await sha256(password).then((password_sha256) => {
        let user_credentials = {
            "username": username,
            "secret_hash": password_sha256,
        };

        areCredentialsValid(user_credentials)
        .then(valid => {
            if (valid === true) {
                localStorage.setItem(LOCALSTORAGE_NAME_KEY, username);
                localStorage.setItem(LOCALSTORAGE_PASSWORD_HASH_KEY, password_sha256); 
                window.location.replace("/");
            } else {
                document.getElementById("error_output").innerHTML = "Invalid credentials";
            }
        })
    });
}