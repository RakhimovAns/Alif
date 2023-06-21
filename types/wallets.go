package types

type Wallet struct {
	Balance int64 `json:"balance"`
}
type Actions struct {
	WalletId int64
	Amount   int64
	Sum      int64
}
