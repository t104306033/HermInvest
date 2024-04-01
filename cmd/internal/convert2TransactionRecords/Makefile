SHELL := /bin/bash # For Completion Now

HEAD = head
BASH = bash
TAIL_CMD = tail -n +2 -q
ECHO_CMD = echo "date,time,stockNo,stockName,tranType,unitPrice,quantity,status"
AWK_CMD = awk -F',' '{printf "%s,--,%s,%s,%s,%s,%s,--\n",$$5,$$2,$$1,$$3,$$6,$$4}'
CONVERT_CMD = $(BASH) $(CONVERT_SCRIPT)

OUTPUT_NAME = ./Output.csv
SRC_DIR = ./commissionHistory/
EXAMPLE_SRC_DIR = ./commissionHistoryExample/
CONVERT_SCRIPT = convertCommissionHistory.sh


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
	CP=copy
else
    BIN_EXT=
    RM=rm -f
	CP=cp
endif

# check if SRC_DIR has csv, otherwise assign to EXAMPLE_DIR
ifeq ("$(wildcard $(SRC_DIR)*.csv)", "")
	SRC_DIR = $(EXAMPLE_SRC_DIR)
endif

.PHONY: all show clean

all: show

$(OUTPUT_NAME): $(SRC_DIR)*.csv $(CONVERT_SCRIPT)
	@$(ECHO_CMD) > $(OUTPUT_NAME)
	@$(TAIL_CMD) $(SRC_DIR)Before_2020-12_excel.csv | $(AWK_CMD) >> $(OUTPUT_NAME)
	@$(TAIL_CMD) $(SRC_DIR)commission_history_* | $(CONVERT_CMD) >> $(OUTPUT_NAME)

show: clean $(OUTPUT_NAME)
	@$(HEAD) $(OUTPUT_NAME)
	@echo "" 
	@echo "Now the SRC_DIR is '$(SRC_DIR)'" 

clean:
	@$(RM) $(OUTPUT_NAME)