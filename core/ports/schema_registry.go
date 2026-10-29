package ports

import (
	"context"

	"clienthub_go/core/models"

	"github.com/antoniomarfa/hexatools/repository"
)

// UserRepositoy interface
type SchemaRegistryRepository interface {
	repository.Repository
	GetByIDUpdate(ctx context.Context, ID string) (interface{}, error)
}

// SaleService interface
type SchemaRegistryService interface {
	Create(ctx context.Context, schema models.TokenSchemaRegistry) (string, error)
	GetAll(ctx context.Context, filter map[string]interface{}) (*models.TokenSchemaRegistryResponse, error)
}
