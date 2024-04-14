SHELL := /bin/bash # For Completion Now

HEAD = head
BASH = bash
TAIL_CMD = tail -n +2 -q
ECHO_CMD = echo "date,time,stockNo,stockName,tranType,quantity,unitPrice,source"
MANUAL_CONV = $(BASH) $(MI_SCRIPT)
COMMISSION_CONV = $(BASH) $(CH_SCRIPT)
MAPPING_CONV = $(BASH) $(SM_SCRIPT)

TXN_REC = ./tblTransactionRecord.csv
STK_MAP = ./tblStockMapping.csv
SRC_DIR = ./commissionHistory/
EXAMPLE_SRC_DIR = ./commissionHistoryExample/
MI_SCRIPT = convertManualInput.sh
CH_SCRIPT = convertCommissionHistory.sh
SM_SCRIPT = convertStockMapping.sh


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

$(TXN_REC): $(SRC_DIR)*.csv $(MI_SCRIPT) $(CH_SCRIPT)
	@$(ECHO_CMD) > $(TXN_REC)
	@$(TAIL_CMD) $(SRC_DIR)Before_2020-12_excel.csv | $(MANUAL_CONV) >> $(TXN_REC)
	@$(TAIL_CMD) $(SRC_DIR)commission_history_* | $(COMMISSION_CONV) >> $(TXN_REC)

$(STK_MAP): $(TXN_REC) $(SM_SCRIPT)
	@echo "stockNo,stockName" > $(STK_MAP)
	@$(TAIL_CMD) $(TXN_REC) | $(MAPPING_CONV) >> $(STK_MAP)

show: clean $(TXN_REC) $(STK_MAP)
	@echo "$(TXN_REC)"
	@$(HEAD) $(TXN_REC)
	@echo ""
	@echo "$(STK_MAP)"
	@$(HEAD) $(STK_MAP)
	@echo "" 
	@echo "Now the SRC_DIR is '$(SRC_DIR)'" 

clean:
	@$(RM) $(TXN_REC)
	@$(RM) $(STK_MAP)