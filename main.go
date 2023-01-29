package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/e-commerce-microservices/shop-service/pb"
	"github.com/e-commerce-microservices/shop-service/repository"
	"github.com/e-commerce-microservices/shop-service/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// postgres driver
	_ "github.com/lib/pq"
)

func main() {
	// create grpc server
	grpcServer := grpc.NewServer()

	// init shop db connection
	pgDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWD"), os.Getenv("DB_DBNAME"),
	)

	shopDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer shopDB.Close()
	if err := shopDB.Ping(); err != nil {
		log.Fatal("can't ping to user db", err)
	}

	// init shop queries
	shopQueries := repository.New(shopDB)

	// dial auth client
	authServiceConn, err := grpc.Dial("auth-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can't dial auth service: ", err)
	}
	// create auth client
	authClient := pb.NewAuthServiceClient(authServiceConn)

	// dial user client
	userServiceConn, err := grpc.Dial("user-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can't dial user service: ", err)
	}
	// create auth client
	userClient := pb.NewUserServiceClient(userServiceConn)

	// dial product client
	productServiceConn, err := grpc.Dial("product-service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	productClient := pb.NewProductServiceClient(productServiceConn)

	// create shop service
	shopService := service.NewShopService(shopQueries, authClient, userClient, productClient)
	// register shop service
	pb.RegisterShopServiceServer(grpcServer, shopService)

	// listen and serve
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
