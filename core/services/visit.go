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
type visitService struct {
	config     config.Config
	repository ports.VisitRepository
}

// NewUvisitService creates a new user service
func NewVisitService(cfg config.Config, repo ports.VisitRepository) ports.VisitService {
	return &visitService{
		config:     cfg,
		repository: repo,
	}
}

// Create roles
func (p *visitService) Create(ctx context.Context, visit models.Visit) (string, error) {

	insertedID, err := p.repository.Create(ctx, models.Visit(visit))
	if err != nil {
		return "", err
	}

	return insertedID, err
}

// GetAll users
func (p *visitService) GetAll(ctx context.Context, filter map[string]interface{}) (*models.VisitListResponse, error) {
	// Obtiene los roles desde el repositorio
	result, err := p.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convierte los resultados
	if len(result) == 0 {
		return &models.VisitListResponse{
			Items:      []models.Visit{},
			TotalCount: 0,
		}, nil
	}

	response, ok := result[0].(models.VisitListResponse)
	if !ok {
		return nil, fmt.Errorf("tipo de respuesta inesperado del repositorio")
	}

	return &response, nil
}

// GetByID user
func (p *visitService) GetByID(ctx context.Context, ID string) (resp models.Visit, err error) {
	roles, err := p.repository.GetByID(ctx, ID)

	if err != nil {
		return resp, fmt.Errorf("roles con ID %s no encontrado", ID)
	}

	if roles == nil {
		// Si no se encuentra el colegio (colegios es nil), devolver un valor en blanco y un error
		return models.Visit{}, fmt.Errorf("roles con ID %s no encontrado", ID)
	}

	resp = *roles.(*models.Visit)

	return
}

// Update user
func (p *visitService) Update(ctx context.Context, ID string, visit models.Visit) (err error) {

	dbVisit, err := p.GetByID(ctx, ID)
	if err != nil {
		return
	}
	// Actualizar los campos solo si no son nil
	// Actualizar la fecha de modificación
	//now := time.Now().UTC()
	if visit.CustomerType != "" {
		dbVisit.CustomerType = *&visit.CustomerType
	}

	// Llamar al repositorio para actualizar la entidad
	err = p.repository.Update(ctx, ID, models.Visit(dbVisit))

	return err
}

// Delete user
func (p *visitService) Delete(ctx context.Context, ID string) (err error) {
	err = p.repository.Delete(ctx, ID, nil)
	if errors.Is(err, wrappers.NonExistentErr) {
		err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
	}

	return err
}
