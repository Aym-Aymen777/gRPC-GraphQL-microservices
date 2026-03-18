package catalog

import (
	"context"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	service := pb.NewCatalogServiceClient(conn)
	return &Client{
		conn:    conn,
		service: service,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) AddProduct(name string, description string, price float64) error {
	_, err := c.service.AddProduct(context.Background(), &pb.AddProductRequest{
		Product: &pb.Product{
			Name:        name,
			Description: description,
			Price:       price,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetProductDetails(id string) (*pb.Product, error) {
	resp, err := c.service.GetProductDetails(context.Background(), &pb.GetProductDetailsRequest{Id: id})	
	if err != nil {
		return nil, err
	}
	return resp.Product, nil
}

func (c *Client) UpdateProduct(id string, name string, description string, price float64) error {
	_, err := c.service.UpdateProduct(context.Background(), &pb.UpdateProductRequest{
		Product: &pb.Product{
			Id:          id,
			Name:        name,
			Description: description,
			Price:       price,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveProduct(id string) error {
	_, err := c.service.RemoveProduct(context.Background(), &pb.RemoveProductRequest{Id: id})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListProducts(skip, limit int) ([]*pb.Product, error) {
	resp, err := c.service.ListProducts(context.Background(), &pb.ListProductsRequest{
		Skip: int32(skip),
		Limit: int32(limit),
	})
	if err != nil {
		return nil, err
	}
	return resp.Products, nil
}

func (c *Client) SearchForProducts(query string, skip, limit int) ([]*pb.Product, error) {
	resp, err := c.service.SearchForProducts(context.Background(), &pb.SearchForProductsRequest{
		Query: query,
		Skip: int32(skip),
		Limit: int32(limit),
	})	
	if err != nil {
		return nil, err
	}
	return resp.Products, nil
}