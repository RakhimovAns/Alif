package postgresql

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/RakhimovAns/Alif/pkg"
	"github.com/RakhimovAns/Alif/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type CustomerService struct {
	pool *pgxpool.Pool
}

func NewCustomerService(pool *pgxpool.Pool) *CustomerService {
	return &CustomerService{pool: pool}
}

func (s *CustomerService) Register(ctx context.Context, customer *types.Customer) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(customer.Password))
	if err != nil {
		log.Println(err)
		return types.ErrInvalidPassword
	}
	id := uuid.New().String()
	_, err = s.pool.Exec(ctx, `
		insert into customers(id,name,login,password,role_id) values ($1,$2,$3,$4,$5) on conflict (login) do nothing 
`, id, customer.Name, customer.Login, hex.EncodeToString(hash), 2)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (s *CustomerService) Login(ctx context.Context, login string, password string, requestBody map[string]string) (string, string, error) {
	var hash string
	var id string
	err := s.pool.QueryRow(ctx, `
		select id, password from customers where login=$1
`, login).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", "", types.ErrNoSuchUser
	}
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return "", "", types.ErrInvalidPassword
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.TokenClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		id,
	})
	TokenStr, err := token.SignedString([]byte("My Key"))
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", "", err
	}
	Key := "Ansar"
	xDigest, err := pkg.CalculateHMACSHA1(requestBody, Key)
	if err != nil {
		return "", "", err
	}
	_, err = s.pool.Exec(ctx, `
update customers  set role_id=2 where customers.id=$1
`, id)
	if err == pgx.ErrNoRows {
		return "", "", types.ErrNoSuchUser
	} else if err != nil {
		return "", "", err
	}
	_, err = s.pool.Exec(ctx, `
update wallets  set role_id=2 where wallets.customer_id=$1
`, id)
	return TokenStr, xDigest, nil
}

func (s *CustomerService) GetID(ctx context.Context, customer *types.Customer) (string, error) {
	var hash string
	var id string
	err := s.pool.QueryRow(ctx, `
		select id, password from customers where login=$1
`, customer.Login).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		return "", types.ErrNoSuchUser
	}
	hashed, err := hex.DecodeString(hash)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(hashed, []byte(customer.Password))
	if err != nil {
		return "", types.ErrInvalidPassword
	}
	return id, nil
}
