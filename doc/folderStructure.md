Folder Structure
===

Possible directory structure, no idea yet now.

project/
├── cmd/
│   └── cli/
│       ├── main.go            # Entry point for the CLI application
│       └── commands/          # Folder containing CLI operation commands
│           ├── insert.go      # Command to insert stock holding information
│           ├── delete.go      # Command to delete stock holding information
│           ├── update.go      # Command to update stock holding information
│           └── query.go       # Command to query stock holding information
│
└── internal/
    ├── app/
    │   ├── database/
    │   │   └── sqlite.go       # Operations related to SQLite
    │   ├── stocks/
    │   │   └── stock.go        # Structure and operations for stock holding information
    │   └── cli/
    │       └── handler.go      # Handler for CLI operations
    └── config/
        └── config.go           # Configuration file

