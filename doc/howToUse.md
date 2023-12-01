How to Use
===

## Stock Management CRUD

### 1. Add Stock
- **Input**: [id], stockNo, type, quantity, unitPrice, [date=today]
- **Calculations**: Calculate totalAmount, taxes
- **Action**: Insert into `tblTransaction` with fields - id, stockNo, type, quantity, unitPrice, date, totalAmount, taxes

### 2. Update Stock Unit Price
- **Input**: id, unitPrice
- **Calculations**: Recalculate totalAmount, taxes
- **Action**: Update `tblTransaction` with fields - id, unitPrice, totalAmount, taxes

### 3. Delete Stock
- **Input**: id
- **Confirmation**: Yes or No
- **Action**: Delete from `tblTransaction` with id

### 4. Query Stock
- **Input**: id or stockNo or type or date [summary]
- **Query Action**: Retrieve data from `tblTransaction` based on id, stockNo, type, or date, with stockName from `tblStockMapping`
- **Output Fields**: id, stockNo, stockName, type, quantity, unitPrice, date, totalAmount, taxes

## Database Schema

### Table: tblStockMapping
- **Columns**:
  - stockNo: TEXT (NOT NULL, UNIQUE)
  - stockName: TEXT (NOT NULL)
- **Primary Key**: stockNo

### Table: tblTransaction
- **Columns**:
  - id: INTEGER (NOT NULL, UNIQUE)
  - stockNo: TEXT (NOT NULL, Foreign Key to tblStockMapping)
  - date: TEXT
  - quantity: INTEGER (NOT NULL)
  - type: INTEGER (NOT NULL)
  - unitPrice: REAL (NOT NULL)
  - totalAmount: INTEGER
  - taxes: INTEGER
- **Primary Key**: id
- **Foreign Key Reference**: stockNo (References tblStockMapping stockNo)


