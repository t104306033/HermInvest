SHELL := /bin/bash # For Completion Now

HEAD = head
BASH = bash
TAIL_CMD = tail -n +2 -q
ECHO_CMD = echo "date,time,stockNo,stockName,tranType,quantity,unitPrice,status"
MANUAL_CONV = $(BASH) $(MI_SCRIPT)
COMMISSION_CONV = $(BASH) $(CH_SCRIPT)

OUTPUT_NAME = ./tblTransactionRecord.csv
SRC_DIR = ./commissionHistory/
EXAMPLE_SRC_DIR = ./commissionHistoryExample/
MI_SCRIPT = convertManualInput.sh
CH_SCRIPT = convertCommissionHistory.sh


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

$(OUTPUT_NAME): $(SRC_DIR)*.csv $(MI_SCRIPT) $(CH_SCRIPT)
	@$(ECHO_CMD) > $(OUTPUT_NAME)
	@$(TAIL_CMD) $(SRC_DIR)Before_2020-12_excel.csv | $(MANUAL_CONV) >> $(OUTPUT_NAME)
	@$(TAIL_CMD) $(SRC_DIR)commission_history_* | $(COMMISSION_CONV) >> $(OUTPUT_NAME)

show: clean $(OUTPUT_NAME)
	@$(HEAD) $(OUTPUT_NAME)
	@echo "" 
	@echo "Now the SRC_DIR is '$(SRC_DIR)'" 

clean:
	@$(RM) $(OUTPUT_NAME)