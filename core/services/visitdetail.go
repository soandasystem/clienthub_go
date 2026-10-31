package services

import (
	"context"
	"errors"
	"fmt"

	"clienthub_go/config"
	"clienthub_go/core/models"
	"clienthub_go/core/ports"

	"github.com/antoniomarfa/hexatools/wrappers"
)

// visitService adapter of an user service
type visitDetailService struct {
	config     config.Config
	repository ports.VisitDetailRepository
}

// NewUvisitService creates a new user service
func NewVisitDetailService(cfg config.Config, repo ports.VisitRepository) ports.VisitDetailService {
	return &visitDetailService{
		config:     cfg,
		repository: repo,
	}
}

// Create roles
func (p *visitDetailService) Create(ctx context.Context, visit models.VisitDetail) (string, error) {

	insertedID, err := p.repository.Create(ctx, models.VisitDetail(visit))
	if err != nil {
		return "", err
	}

	return insertedID, err
}

// GetAll users
func (p *visitDetailService) GetAll(ctx context.Context, filter map[string]interface{}) (*models.VisitDetailListResponse, error) {
	// Obtiene los roles desde el repositorio
	result, err := p.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convierte los resultados
	if len(result) == 0 {
		return &models.VisitDetailListResponse{
			Items:      []models.VisitDetail{},
			TotalCount: 0,
		}, nil
	}

	response, ok := result[0].(models.VisitDetailListResponse)
	if !ok {
		return nil, fmt.Errorf("tipo de respuesta inesperado del repositorio")
	}

	return &response, nil
}

// GetByID user
func (p *visitDetailService) GetByID(ctx context.Context, ID string) (resp models.VisitDetail, err error) {
	roles, err := p.repository.GetByID(ctx, ID)

	if err != nil {
		return resp, fmt.Errorf("roles con ID %s no encontrado", ID)
	}

	if roles == nil {
		// Si no se encuentra el colegio (colegios es nil), devolver un valor en blanco y un error
		return models.VisitDetail{}, fmt.Errorf("roles con ID %s no encontrado", ID)
	}

	resp = *roles.(*models.VisitDetail)

	return
}

// Update user
func (p *visitDetailService) Update(ctx context.Context, ID string, visit models.VisitDetail) (err error) {

	dbVisit, err := p.GetByID(ctx, ID)
	if err != nil {
		return
	}
	// Actualizar los campos solo si no son nil
	// Actualizar la fecha de modificación
	//now := time.Now().UTC()
	if *visit.Complaints != "" {
		dbVisit.Competitors = *&visit.Complaints
	}

	// Llamar al repositorio para actualizar la entidad
	err = p.repository.Update(ctx, ID, models.VisitDetail(dbVisit))

	return err
}

// Delete user
func (p *visitDetailService) Delete(ctx context.Context, ID string) (err error) {
	err = p.repository.Delete(ctx, ID, nil)
	if errors.Is(err, wrappers.NonExistentErr) {
		err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
	}

	return err
}
