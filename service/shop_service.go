package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/e-commerce-microservices/shop-service/pb"
	"github.com/e-commerce-microservices/shop-service/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shopRepository interface {
	CreateShop(ctx context.Context, arg repository.CreateShopParams) error
}

// ShopService ...
type ShopService struct {
	shopStore     shopRepository
	authClient    pb.AuthServiceClient
	userClient    pb.UserServiceClient
	productClient pb.ProductServiceClient

	pb.UnimplementedShopServiceServer
}

// NewShopService ...
func NewShopService(shopStore shopRepository, authClient pb.AuthServiceClient, userClient pb.UserServiceClient, productClient pb.ProductServiceClient) *ShopService {
	service := &ShopService{
		shopStore:     shopStore,
		authClient:    authClient,
		userClient:    userClient,
		productClient: productClient,
	}

	return service
}

// RegisterShop ...
func (srv *ShopService) RegisterShop(ctx context.Context, req *pb.RegisterShopRequest) (*pb.GeneralResponse, error) {
	// auth
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "can't parse context")
	}
	ctx = metadata.NewOutgoingContext(ctx, md)

	// TODO: two phase commit
	// update user role
	_, err := srv.userClient.SupplierRegister(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	me, err := srv.userClient.GetMe(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	// create shop
	err = srv.shopStore.CreateShop(ctx, repository.CreateShopParams{
		SellerID: me.GetId(),
		Name:     req.GetName(),
		Avatar: sql.NullString{
			String: req.GetAvatar(),
			Valid:  false,
		},
	})
	if err != nil {
		return nil, err
	}

	return &pb.GeneralResponse{
		Message: "Đăng kí thành công, hãy thử bán hàng ngay lập tức",
	}, nil
}

var _empty = &emptypb.Empty{}

// AddProduct ...
func (srv *ShopService) AddProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	// authorization for supplier or admin
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	claims, err := srv.authClient.GetUserClaims(ctx, _empty)
	if err != nil {
		return nil, err
	}
	if claims.GetUserRole() == pb.UserRole_customer {
		return nil, errors.New("unauthorization request")
	}
	supplierID, _ := strconv.ParseInt(claims.Id, 10, 64)
	req.SupplierId = supplierID

	// add product
	return srv.productClient.CreateProduct(ctx, req)
}

// Ping pong
func (ShopService) Ping(ctx context.Context, _ *empty.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Message: "pong",
	}, nil
}
