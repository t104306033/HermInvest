Convert2TransactionRecords
===

Convert2TransactionRecords is used to convert commission history records, sourced from "exported from the stock system" and "manually inputted data." These records are standardized and can be imported into the hermInvest database for use.

[TOC]

### Folder Structure
* `convertManualInput.sh`: Converts manually inputted records.
* `convertCommissionHistory.sh`: Converts records exported from the stock system.
* `commissionHistory/`: Directory for commission history records. Place the file here.
* `commissionHistoryExample/`: Directory for example commission history records.

### How to use

Before proceeding, ensure you have `make` installed on your system. To execute the makefile, run `make` in the terminal.

The Makefile will convert the files of the `commissionHistory/` directory and generate a file named `tblTransactionRecord.csv`, which contains the converted transaction records. After that, You can then import this file using sqlitebrowser.

Note: If the `commissionHistory/` directory is empty, the Makefile will use records from `commissionHistoryExample/` for you.

### Important Notes
1. To ensure data integrity, the "date" and "time" fields for manually inputted records will be converted as specified in convertManualInput.sh.
2. The "source" field is used to distinguish the source of data:
    ```
    1: Manual input
    2: Exported from the stock system
    3: Command Line Interface (CLI)
    4: Web
    ```
