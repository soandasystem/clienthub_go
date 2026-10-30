package ports

import (
	"context"

	"clienthub_go/core/models"

	"github.com/antoniomarfa/hexatools/repository"
)

// UserRepositoy interface
type VisitRepository interface {
	repository.Repository
}

// SaleService interface
type VisitService interface {
	Create(ctx context.Context, roles models.Visit) (string, error)
	GetAll(ctx context.Context, filter map[string]interface{}) (*models.VisitListResponse, error)
	GetByID(ctx context.Context, ID string) (models.Visit, error)
	Update(ctx context.Context, ID string, roles models.Visit) error
	Delete(ctx context.Context, ID string) error
}
