BIN_DIR:=bin
SRC_DIR:=src
EXE_NAME:=gochat
PAGES_DIR:=pages
STATIC_DIR:=static
SCRIPTS_DIR:=scripts


all: clean
	mkdir $(BIN_DIR)
	cd $(SRC_DIR) && go build && mv $(EXE_NAME) ../$(BIN_DIR)
	cp -r $(PAGES_DIR) $(BIN_DIR)
	cp -r $(STATIC_DIR) $(BIN_DIR)
	cp -r $(SCRIPTS_DIR) $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)