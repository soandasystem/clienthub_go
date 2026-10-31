package ports

import (
	"context"

	"clienthub_go/core/models"

	"github.com/antoniomarfa/hexatools/repository"
)

// UserRepositoy interface
type VisitDetailRepository interface {
	repository.Repository
}

// SaleService interface
type VisitDetailService interface {
	Create(ctx context.Context, visit models.VisitDetail) (string, error)
	GetAll(ctx context.Context, filter map[string]interface{}) (*models.VisitDetailListResponse, error)
	GetByID(ctx context.Context, ID string) (models.VisitDetail, error)
	Update(ctx context.Context, ID string, visit models.VisitDetail) error
	Delete(ctx context.Context, ID string) error
}
