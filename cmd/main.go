package main

import (
	"context"
	"github.com/RakhimovAns/Alif/cmd/server"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgresql://postgres:Ansar@localhost:5432/test"
	if err := execute(host, port, dsn); err != nil {
		log.Fatal(err)
	}
}
func execute(host string, port string, dsn string) (err error) {
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, err := pgxpool.Connect(connectCtx, dsn)
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()
	router := mux.NewRouter()
	server := server.NewServer(router)
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()

}
