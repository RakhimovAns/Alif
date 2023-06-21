package postgresql

import (
	"context"
	"encoding/hex"
	"github.com/RakhimovAns/Alif/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type WalletService struct {
	pool *pgxpool.Pool
}

func NewWalletService(pool *pgxpool.Pool) *WalletService {
	return &WalletService{pool: pool}
}

func (s *WalletService) CreateWallet(ctx context.Context, customer *types.Customer) error {
	var hash string
	var id string
	var roleId int
	err := s.pool.QueryRow(ctx, `
    select id, password ,role_id from customers where login=$1
`, customer.Login).Scan(&id, &hash, &roleId)
	if err == pgx.ErrNoRows {
		return types.ErrNoSuchUser
	}
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		log.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(customer.Password))
	if err != nil {
		return types.ErrInvalidPassword
	}
	if customer.Balance > 10_000 {
		return types.ErrNotIdentified
	}
	_, err = s.pool.Exec(ctx, `
insert into wallets(balance, customer_id,role_id) VALUES ($1,$2,$3)
`, customer.Balance, id, roleId)
	return nil
}

func (s *WalletService) CheckWallet(ctx context.Context, id string) error {
	var ID int64
	err := s.pool.QueryRow(ctx, `
    select id from wallets where wallets.customer_id=$1
`, id).Scan(&ID)
	if err == pgx.ErrNoRows {
		return types.ErrNoSuchUser
	}
	if err != nil {
		return err
	}
	return nil
}
