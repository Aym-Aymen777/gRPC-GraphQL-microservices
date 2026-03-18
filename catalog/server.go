package catalog

import (
	"context"
	"log"
	"net"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog/pb"
	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog/types"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Service
}

func ListenGRPC(s Service, port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{service: s})
	if err := serv.Serve(listen); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *grpcServer) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	product := &types.Product{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
	}
	err := s.service.AddProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return &pb.AddProductResponse{}, nil
}

func (s *grpcServer) GetProductDetails(ctx context.Context, req *pb.GetProductDetailsRequest) (*pb.GetProductDetailsResponse, error) {
	product, err := s.service.GetProductDetails(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductDetailsResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (s *grpcServer) GetProductsByIDs(ctx context.Context, req *pb.GetProductsByIDsRequest) (*pb.GetProductsByIDsResponse, error) {
	products, err := s.service.GetProductsByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}
	}
	return &pb.GetProductsByIDsResponse{
		Products: pbProducts,
	}, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	product := &types.Product{
		ID:          req.Product.Id,
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
	}
	err := s.service.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductResponse{}, nil
}

func (s *grpcServer) RemoveProduct(ctx context.Context, req *pb.RemoveProductRequest) (*pb.RemoveProductResponse, error) {
	err := s.service.RemoveProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveProductResponse{
		Id: req.Id,
	}, nil
}

func (s *grpcServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.service.ListProducts(ctx, int(req.Skip), int(req.Limit))
	if err != nil {
		return nil, err
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}
	}
	return &pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (s *grpcServer) SearchForProducts(ctx context.Context, req *pb.SearchForProductsRequest) (*pb.SearchForProductsResponse, error) {
	products, err := s.service.SearchForProducts(ctx, req.Query, int(req.Skip), int(req.Limit))
	if err != nil {
		return nil, err
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}
	}
	return &pb.SearchForProductsResponse{
		Products: pbProducts,
	}, nil
}
