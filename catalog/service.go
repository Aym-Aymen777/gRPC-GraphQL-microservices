package catalog

import (
	"context"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog/types"
	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/utils"
)

type Service interface {
	AddProduct(ctx context.Context, product *types.Product) error
	GetProductDetails(ctx context.Context, id string) (*types.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]*types.Product, error)
	UpdateProduct(ctx context.Context, product *types.Product) error
	RemoveProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, skip, limit int) ([]*types.Product, error)
	SearchForProducts(ctx context.Context, query string, skip, limit int) ([]*types.Product, error)
}

type CatalogService struct {
	repo Repository
}

func NewService(repo Repository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) AddProduct(ctx context.Context, product *types.Product) error {
	product.ID = utils.GenerateID()
	return s.repo.CreateProduct(ctx, product)
}

func (s *CatalogService) GetProductDetails(ctx context.Context, id string) (*types.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *CatalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]*types.Product, error) {
	return s.repo.GetProductsByIDs(ctx, ids)
}

func (s *CatalogService) UpdateProduct(ctx context.Context, product *types.Product) error {
	return s.repo.UpdateProduct(ctx, product)
}

func (s *CatalogService) RemoveProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *CatalogService) ListProducts(ctx context.Context, skip, limit int) ([]*types.Product, error) {
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 || (skip == 0 && limit == 0) {
		limit = 100
	}
	return s.repo.GetAllProducts(ctx, skip, limit)
}

func (s *CatalogService) SearchForProducts(ctx context.Context, query string, skip, limit int) ([]*types.Product, error) {
	return s.repo.SearchForProducts(ctx, query, skip, limit)
}
