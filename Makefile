# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# Paths and names
SRC_NAME=*
SRC_FOLDER=./cmd/hermInvestCli
SRC_PATH=$(SRC_FOLDER)/$(SRC_NAME)
BIN_NAME=hermInvestCli

all: exec

$(BIN_NAME).exe: $(SRC_PATH).go
	$(MAKE) build

build:
	$(GOBUILD) -o $(BIN_NAME).exe $(SRC_PATH).go

exec: $(BIN_NAME).exe
	$(BIN_NAME).exe

clean:
	$(GOCLEAN)
	del /f $(BIN_NAME).exe

run:
	$(GOCMD) run $(SRC_PATH).go