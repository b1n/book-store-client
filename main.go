package main

import (
	"context"
	"github.com/b1n/proto-book-store"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	bookStore := GetBookStoreConn()
	getBookRequest := &book_store.GetBookRequest{
		Id: 1,
	}
	response, err := bookStore.GetBook(context.Background(), getBookRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(response)
}

type tokenAuth struct {
	token string
}

// Return value is mapped to request headers.
func (t *tokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"access-token": t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return false
}

func GetBookStoreConn() book_store.BookStoreClient {
	tokenAuth := &tokenAuth{token: os.Getenv("TOKEN")}
	conn, err := grpc.Dial(os.Getenv("ADDRESS"), grpc.WithPerRPCCredentials(tokenAuth), grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	return book_store.NewBookStoreClient(conn)
}
