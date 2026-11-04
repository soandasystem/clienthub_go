package ports

import (
	"context"

	"clienthub_go/core/models"

	"github.com/antoniomarfa/hexatools/repository"
)

// UserRepositoy interface
type ProspectRepository interface {
	repository.Repository
}

// SaleService interface
type ProspectService interface {
	Create(ctx context.Context, prospect models.Prospect) (string, error)
	GetAll(ctx context.Context, filter map[string]interface{}) (*models.ProspectListResponse, error)
	GetByID(ctx context.Context, ID string) (models.Prospect, error)
	Update(ctx context.Context, ID string, prospect models.UpdateProspect) error
	Delete(ctx context.Context, ID string) error
}
