## gochat
### A dead simple real time webchat with authorization using websockets

### Building
`make` or `cd src/ && go build && mv gochat ..`

### Project structure
As you are probably here for reconnaissance reasons for your own project here's a
quick overview of the structure:

- `pages` -> html pages
- `scripts` -> js files where chat.js is the "heavy-lifter" which actually implements the websocket chat logic 
- `static` -> icon, stylesheet for pages
- `src/api` -> the fundament of backend logic, authorization, database and websocket handling,all necessary structs, constants
- `src/log` -> a custom logger
- `src/page` -> a shorcut function for handling html templates
- `src/server` -> the main server struct that glues everything together. For a high-level inspection - refer there
- `src/main.go` -> parse command line flags, set up logging and leave the rest for the server

### License
GPLv3