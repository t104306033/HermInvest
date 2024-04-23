package service

import (
	"HermInvest/pkg/repository"
	"fmt"
)

func InitializeService() *service {
	db, err := repository.GetDBConnection()
	if err != nil {
		fmt.Println("Error geting DB connection: ", err)
	}

	// init transactionRepository
	repo := repository.NewRepository(db)

	serv := NewService(repo)

	return serv
}
