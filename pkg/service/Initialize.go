package service

import (
	"HermInvest/pkg/repository"
)

func InitializeService() *service {
	db := repository.GetDBConnection()

	// init transactionRepository
	repo := repository.NewRepository(db)

	serv := NewService(repo)

	return serv
}
