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

# Detect the OS
ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
    detected_OS := Windows
else
    detected_OS := $(shell uname)  # same as "uname -s"  is Linux
endif

# Windows-specific variables and commands
ifeq ($(detected_OS),Windows)
    BIN_EXT=.exe
    RM=del /f
else
    BIN_EXT=
    RM=rm -f
endif

all: exec

$(BIN_NAME)$(BIN_EXT): $(SRC_PATH).go
	$(MAKE) build

build:
	$(GOBUILD) -o $(BIN_NAME)$(BIN_EXT) $(SRC_PATH).go

exec: $(BIN_NAME)$(BIN_EXT)
	./$(BIN_NAME)$(BIN_EXT)

clean:
	$(GOCLEAN)
	$(RM) $(BIN_NAME)$(BIN_EXT)

run:
	$(GOCMD) run $(SRC_PATH).go