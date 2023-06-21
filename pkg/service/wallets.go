package service

import (
	"context"
	"github.com/RakhimovAns/Alif/pkg/postgresql"
	"github.com/RakhimovAns/Alif/types"
)

type WalletService struct {
	service postgresql.WalletService
}

func NewWalletService(service *postgresql.WalletService) *WalletService {
	return &WalletService{service: *service}
}

func (s *WalletService) CreateWallet(ctx context.Context, customer *types.Customer) error {
	return s.service.CreateWallet(ctx, customer)
}
func (s *WalletService) CheckWallet(ctx context.Context, id string) error {
	return s.service.CheckWallet(ctx, id)
}
