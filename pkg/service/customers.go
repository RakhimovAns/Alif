package service

import (
	"context"
	"github.com/RakhimovAns/Alif/pkg/postgresql"
	"github.com/RakhimovAns/Alif/types"
)

type CustomerService struct {
	service postgresql.CustomerService
}

func NewCustomerService(service *postgresql.CustomerService) *CustomerService {
	return &CustomerService{service: *service}
}

func (s *CustomerService) Register(ctx context.Context, customer *types.Customer) error {
	return s.service.Register(ctx, customer)
}

func (s *CustomerService) Login(ctx context.Context, login string, password string, requestBody map[string]string) (string, string, error) {
	return s.service.Login(ctx, login, password, requestBody)
}
func (s *CustomerService) GetId(ctx context.Context, customer *types.Customer) (string, error) {
	return s.service.GetID(ctx, customer)
}
