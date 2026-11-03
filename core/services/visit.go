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

func (p *visitService) GetInforme(ctx context.Context, filter map[string]interface{}) (*models.VisitSveltResponse, error) {
	// Obtiene los roles desde el repositorio
	result, err := p.repository.GetInf(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convierte los resultados
	if len(result) == 0 {
		return &models.VisitSveltResponse{
			Items:      []models.VisitSvelt{},
			TotalCount: 0,
		}, nil
	}

	response, ok := result[0].(models.VisitSveltResponse)
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

	dbVisit.StartedAt = visit.StartedAt
	if visit.StartLatitude != nil {
		dbVisit.StartLatitude = visit.StartLatitude
	}
	if visit.StartLongitude != nil {
		dbVisit.StartLongitude = visit.StartLongitude
	}
	if visit.StartAddress != nil {
		dbVisit.StartAddress = visit.StartAddress
	}
	if visit.EndedAt != nil {
		dbVisit.EndedAt = visit.EndedAt
	}
	if visit.EndLatitude != nil {
		dbVisit.EndLatitude = visit.EndLatitude
	}
	if visit.EndLongitude != nil {
		dbVisit.EndLongitude = visit.EndLongitude
	}
	if visit.EndAddress != nil {
		dbVisit.EndAddress = visit.EndAddress
	}
	if visit.DeviceInfo != nil {
		dbVisit.DeviceInfo = visit.DeviceInfo
	}
	if visit.FinalComments != nil {
		dbVisit.FinalComments = visit.FinalComments
	}
	if visit.DeviceInfo != nil {
		dbVisit.DeviceInfo = visit.DeviceInfo
	}
	if visit.FinalComments != nil {
		dbVisit.FinalComments = visit.FinalComments
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
