package types

import "time"

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	SecretKey string    `json:"secretkey"`
	RoleId    int64     `json:"roleid"`
	Balance   int64     `json:"balance"`
	Created   time.Time `json:"created"`
}
