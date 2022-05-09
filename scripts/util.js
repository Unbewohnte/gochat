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


function sha256(string) {
    const encoded_string = new TextEncoder().encode(string);
    return crypto.subtle.digest('SHA-256', encoded_string).then((hash_sum) => {
      const hashArray = Array.from(new Uint8Array(hash_sum));
      const hashHex = hashArray
        .map((bytes) => bytes.toString(16).padStart(2, '0'))
        .join('');
      return hashHex;
    });
}


function to_hex_string(byteArr) {
    return Array.from(byteArr, function(byte) {
        return ('0' + (byte & 0xFF).toString(16).slice(-2));
    }).join('')
}

function is_image(filename) {
    image_exts = ["jpe", "jpeg", "jpg", "png", "ppm", "gif"]
    
    let is_img = false;

    image_exts.find(ext => {
        if (filename.includes(ext)) {
            is_img = true;
            return true;
        }
    });

    return is_img;
}