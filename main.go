package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/b1n/proto-book-store"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Service struct {
	BookStore book_store.BookStoreClient
}

func (s *Service) GetBook(ctx *gin.Context) {
	qId := ctx.Query("id")
	id, err := strconv.Atoi(qId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	getBookRequest := &book_store.GetBookRequest{
		Id: int32(id),
	}
	response, err := s.BookStore.GetBook(context.Background(), getBookRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
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
	target := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	config := &tls.Config{}

	conn, err := grpc.Dial(
		target,
		grpc.WithUnaryInterceptor(interceptor),
		grpc.WithPerRPCCredentials(tokenAuth),
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		log.Println(err)
	}
	return book_store.NewBookStoreClient(conn)
}

func interceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)

	fmt.Printf(`--
	call=%v
	req=%#v
	reply=%#v
	time=%v
	err=%v
`, method, req, reply, time.Since(start), err)

	return err
}
