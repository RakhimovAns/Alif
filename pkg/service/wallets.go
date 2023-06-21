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

func (s *WalletService) CreateWallet(ctx context.Context, id string) error {
	return s.service.CreateWallet(ctx, id)
}
func (s *WalletService) CheckWallet(ctx context.Context, id string) error {
	return s.service.CheckWallet(ctx, id)
}

func (s *WalletService) DepositWallet(ctx context.Context, id string, sum int64) error {
	return s.service.Deposit(ctx, id, sum)
}
func (s *WalletService) GetBalance(ctx context.Context, id string) (int64, error) {
	return s.service.GetBalance(ctx, id)
}
func (s *WalletService) GetActions(ctx context.Context, id string) (*types.Actions, error) {
	return s.service.GetActions(ctx, id)
}
