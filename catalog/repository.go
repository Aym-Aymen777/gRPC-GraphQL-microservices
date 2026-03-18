package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog/types"
	"github.com/elastic/go-elasticsearch/v9"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *types.Product) error
	GetProductByID(ctx context.Context, id string) (*types.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]*types.Product, error)
	UpdateProduct(ctx context.Context, product *types.Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetAllProducts(ctx context.Context, skip, limit int) ([]*types.Product, error)
	SearchForProducts(ctx context.Context, query string, skip, limit int) ([]*types.Product, error)
}

type esSearchResponse struct {
	Hits struct {
		Hits []struct {
			ID     string         `json:"_id"`
			Source types.Product  `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticSearchRepository struct {
	db *elasticsearch.Client
}

func NewElasticSearchRepository(cfg elasticsearch.Config) (Repository, error) {
	db, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ElasticSearchRepository{db: db}, nil
}

func (r *ElasticSearchRepository) CreateProduct(ctx context.Context, product *types.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	res, err := r.db.Index(
		"products",
		bytes.NewReader(data),
		r.db.Index.WithDocumentID(product.ID),
		r.db.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (r *ElasticSearchRepository) GetProductByID(ctx context.Context, id string) (*types.Product, error) {
	res, err := r.db.Get("products", id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Handle ES-level errors
	if res.IsError() {
		return nil, fmt.Errorf("error getting document: %s", res.String())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract _source
	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: %v", result)
	}

	// Map to struct
	product := &types.Product{
		ID:          id,
		Name:        source["name"].(string),
		Description: source["description"].(string),
		Price:       source["price"].(float64),
	}

	return product, nil
}

func (r *ElasticSearchRepository) GetProductsByIDs(ctx context.Context, ids []string) ([]*types.Product, error) {
	// Build valid JSON body
	bodyMap := map[string]interface{}{
		"ids": ids,
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	res, err := r.db.Mget(
		bytes.NewReader(bodyBytes),
		r.db.Mget.WithIndex("products"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting documents: %s", res.String())
	}

	// Define structured response
	var response struct {
		Docs []struct {
			ID     string         `json:"_id"`
			Found  bool           `json:"found"`
			Source *types.Product `json:"_source"`
		} `json:"docs"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var products []*types.Product

	for _, doc := range response.Docs {
		if !doc.Found || doc.Source == nil {
			continue // skip missing docs
		}

		doc.Source.ID = doc.ID // inject ES _id into struct
		products = append(products, doc.Source)
	}

	return products, nil
}

func (r *ElasticSearchRepository) UpdateProduct(ctx context.Context, product *types.Product) error {
	// Check if document exists
	res, err := r.db.Exists("products", product.ID)
	if err != nil {
		return err
	}
	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == 404 {
		return fmt.Errorf("product with ID %s not found", product.ID)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %s", res.Status())
	}

	// Prepare partial update
	data, err := json.Marshal(map[string]interface{}{
		"doc": map[string]interface{}{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
		},
	})
	if err != nil {
		return err
	}

	// Perform update
	updateRes, err := r.db.Update(
		"products",
		product.ID,
		bytes.NewReader(data),
		r.db.Update.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer updateRes.Body.Close()

	if updateRes.IsError() {
		return fmt.Errorf("update failed: %s", updateRes.String())
	}

	return nil
}


func (r *ElasticSearchRepository) DeleteProduct(ctx context.Context, id string) error {
	res, err := r.db.Delete(
		"products",
		id,
		// remove refresh in production or make it optional
		r.db.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Handle not found (idempotent delete)
	if res.StatusCode == 404 {
		return nil
	}

	// Handle other errors
	if res.IsError() {
		return fmt.Errorf("delete failed: %s", res.String())
	}

	return nil
}

func (r *ElasticSearchRepository) GetAllProducts(ctx context.Context, skip, limit int) ([]*types.Product, error) {
	res, err := r.db.Search(
		r.db.Search.WithIndex("products"),
		r.db.Search.WithSize(limit),
		r.db.Search.WithFrom(skip),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error listing products: %s", res.String())
	}

	var response esSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var products []*types.Product
	for _, hit := range response.Hits.Hits {
		p := hit.Source
		p.ID = hit.ID
		products = append(products, &p)
	}

	return products, nil
}

func (r *ElasticSearchRepository) SearchForProducts(ctx context.Context, query string, skip, limit int) ([]*types.Product, error) {
	// Proper ES query
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2", "description"},
			},
		},
	}

	bodyBytes, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	res, err := r.db.Search(
		r.db.Search.WithIndex("products"),
		r.db.Search.WithBody(bytes.NewReader(bodyBytes)),
		r.db.Search.WithSize(limit),
		r.db.Search.WithFrom(skip),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.String())
	}

	var response esSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	var products []*types.Product
	for _, hit := range response.Hits.Hits {
		p := hit.Source
		p.ID = hit.ID
		products = append(products, &p)
	}

	return products, nil
}
