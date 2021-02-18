package main

import (
	"context"
	"fmt"
	"github.com/b1n/proto-book-store"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Service struct {
	BookStore book_store.BookStoreClient
}

func (s *Service) GetBook(ctx *gin.Context) {
	qId := ctx.Query("id")
	id, err := strconv.Atoi(qId)
	getBookRequest := &book_store.GetBookRequest{
		Id: int32(id),
	}
	response, err := s.BookStore.GetBook(context.Background(), getBookRequest)
	if err != nil {
		log.Println(err)
	}

	ctx.JSON(http.StatusOK, response)
}

func main() {
	s := &Service{}
	s.BookStore = GetBookStoreConn()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/test", s.GetBook)

	if err := router.Run(":" + os.Getenv("HTTP_PORT")); err != nil {
		log.Println(err)
	}
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
	target := fmt.Sprintf("%s:%s",os.Getenv("GRPC_HOST"),os.Getenv("GRPC_PORT"))

	conn, err := grpc.Dial(target, grpc.WithPerRPCCredentials(tokenAuth), grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	return book_store.NewBookStoreClient(conn)
}
