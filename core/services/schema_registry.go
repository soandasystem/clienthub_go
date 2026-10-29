package services

import (
	"context"
	"fmt"

	"clienthub_go/config"
	"clienthub_go/core/models"
	"clienthub_go/core/ports"
)

// schemaRegistryService adapter of an user service
type schemaRegistryService struct {
	config     config.Config
	repository ports.SchemaRegistryRepository
}

// NewSchemaRegistryService creates a new user service
func NewSchemaRegistryService(cfg config.Config, repo ports.SchemaRegistryRepository) ports.SchemaRegistryService {
	return &schemaRegistryService{
		config:     cfg,
		repository: repo,
	}
}

// Create schema registry
func (p *schemaRegistryService) Create(ctx context.Context, schema models.TokenSchemaRegistry) (string, error) {

	insertedID, err := p.repository.Create(ctx, models.TokenSchemaRegistry(schema))
	if err != nil {
		return "", err
	}

	return insertedID, err
}

// GetAll users
func (p *schemaRegistryService) GetAll(ctx context.Context, filter map[string]interface{}) (*models.TokenSchemaRegistryResponse, error) {
	// Obtiene las vemtas desde el repositorio

	result, err := p.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convierte los resultados
	if len(result) == 0 {
		return &models.TokenSchemaRegistryResponse{
			Items:      []models.TokenSchemaRegistry{},
			TotalCount: 0,
		}, nil
	}

	response, ok := result[0].(models.TokenSchemaRegistryResponse)
	if !ok {
		return nil, fmt.Errorf("tipo de respuesta inesperado del repositorio")
	}

	return &response, nil
}
