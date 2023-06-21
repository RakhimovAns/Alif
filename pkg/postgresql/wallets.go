package postgresql

import (
	"context"
	"github.com/RakhimovAns/Alif/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type WalletService struct {
	pool *pgxpool.Pool
}

func NewWalletService(pool *pgxpool.Pool) *WalletService {
	return &WalletService{pool: pool}
}

func (s *WalletService) CreateWallet(ctx context.Context, id string) error {
	var hash string
	var roleId int
	err := s.pool.QueryRow(ctx, `
    select id, password ,role_id from customers where id=$1
`, id).Scan(&id, &hash, &roleId)
	if err == pgx.ErrNoRows {
		return types.ErrNoSuchUser
	}
	_, err = s.pool.Exec(ctx, `
insert into wallets(balance, customer_id,role_id) VALUES ($1,$2,$3)
`, 0, id, roleId)
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
func (s *WalletService) Deposit(ctx context.Context, id string, sum int64) error {
	var ID int64
	err := s.pool.QueryRow(ctx, `
    select id from wallets where wallets.customer_id=$1
`, id).Scan(&ID)
	_, err = s.pool.Exec(ctx, `
update wallets set balance=wallets.balance+$1 where wallets.customer_id=$2
`, sum, id)
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.IntoActions(ctx, ID, sum)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *WalletService) GetBalance(ctx context.Context, id string) (sum int64, err error) {
	err = s.pool.QueryRow(ctx, `
select balance from wallets where customer_id=$1
`, id).Scan(&sum)
	if err != nil {
		return -1, err
	}
	return sum, nil
}

func (s *WalletService) GetActions(ctx context.Context, id string) (*types.Actions, error) {
	var ID int64
	err := s.pool.QueryRow(ctx, `
    select id from wallets where wallets.customer_id=$1
`, id).Scan(&ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	item := &types.Actions{}
	err = s.pool.QueryRow(ctx, `
select wallet_id ,amount,sum from actions where wallet_id=$1
`, ID).Scan(&item.WalletId, &item.Amount, &item.Sum)
	return item, nil
}

func (s *WalletService) IntoActions(ctx context.Context, id int64, sum int64) error {
	err := s.pool.QueryRow(ctx, `
select wallet_id from actions where wallet_id=$1
`, id)
	if err == nil {
		_, err1 := s.pool.Exec(ctx, `
insert into actions(wallet_id,amount, sum) VALUES ($1,$3,$2) 
`, id, sum, 1)
		if err1 != nil {
			log.Println(err)
			return err1
		}
	} else {
		_, err1 := s.pool.Exec(ctx, `
update actions set amount= actions.amount+1 , sum=actions.sum+$1 where wallet_id=$2
`, sum, id)
		if err1 != nil {
			log.Println(err)
			return err1
		}
	}
	return nil
}
