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

// rolesService adapter of an user service
type prospectService struct {
	config     config.Config
	repository ports.ProspectRepository
}

// NewURolesService creates a new user service
func NewProspectService(cfg config.Config, repo ports.ProspectRepository) ports.ProspectService {
	return &prospectService{
		config:     cfg,
		repository: repo,
	}
}

// Create roles
func (p *prospectService) Create(ctx context.Context, prospect models.Prospect) (string, error) {

	insertedID, err := p.repository.Create(ctx, models.Prospect(prospect))
	if err != nil {
		return "", err
	}

	return insertedID, err
}

// GetAll users
func (p *prospectService) GetAll(ctx context.Context, filter map[string]interface{}) (*models.ProspectListResponse, error) {
	// Obtiene los roles desde el repositorio
	result, err := p.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		return nil, err
	}

	// Convierte los resultados
	if len(result) == 0 {
		return &models.ProspectListResponse{
			Items:      []models.Prospect{},
			TotalCount: 0,
		}, nil
	}

	response, ok := result[0].(models.ProspectListResponse)
	if !ok {
		return nil, fmt.Errorf("tipo de respuesta inesperado del repositorio")
	}

	return &response, nil
}

// GetByID user
func (p *prospectService) GetByID(ctx context.Context, ID string) (resp models.Prospect, err error) {
	prospect, err := p.repository.GetByID(ctx, ID)

	if err != nil {
		return resp, fmt.Errorf("prospect con ID %s no encontrado", ID)
	}

	if prospect == nil {
		return models.Prospect{}, fmt.Errorf("prospect con ID %s no encontrado", ID)
	}

	resp = *prospect.(*models.Prospect)

	return
}

// Update user
func (p *prospectService) Update(ctx context.Context, ID string, Prospect models.UpdateProspect) (err error) {

	dbProspect, err := p.GetByID(ctx, ID)
	if err != nil {
		return
	}
	// Actualizar los campos solo si no son nil
	// Actualizar la fecha de modificación
	if Prospect.CompanyName != "" {
		dbProspect.CompanyName = Prospect.CompanyName
	}
	if Prospect.ContactName != nil {
		dbProspect.ContactName = Prospect.ContactName
	}
	if Prospect.Phone != nil {
		dbProspect.Phone = Prospect.Phone
	}
	if Prospect.Email != nil {
		dbProspect.Email = Prospect.Email
	}
	if Prospect.Address != nil {
		dbProspect.Address = Prospect.Address
	}
	if Prospect.City != nil {
		dbProspect.City = Prospect.City
	}
	if Prospect.BusinessActivity != nil {
		dbProspect.BusinessActivity = Prospect.BusinessActivity
	}
	if Prospect.Status != "" {
		dbProspect.Status = Prospect.Status
	}
	if Prospect.AssignedSellerID != nil {
		dbProspect.AssignedSellerID = Prospect.AssignedSellerID
	}
	if Prospect.ConvertedCRMCustomerID != nil {
		dbProspect.ConvertedCRMCustomerID = Prospect.ConvertedCRMCustomerID
	}
	if Prospect.Type != nil {
		dbProspect.Type = Prospect.Type
	}
	if Prospect.CreatedBy != nil {
		dbProspect.CreatedBy = *Prospect.CreatedBy
	}
	if Prospect.City != nil {
		dbProspect.City = Prospect.City
	}
	if Prospect.Region != nil {
		dbProspect.Region = Prospect.Region
	}
	if Prospect.Type != nil {
		dbProspect.Type = Prospect.Type
	}

	// Llamar al repositorio para actualizar la entidad
	err = p.repository.Update(ctx, ID, dbProspect)

	return err
}

// Delete user
func (p *prospectService) Delete(ctx context.Context, ID string) (err error) {
	err = p.repository.Delete(ctx, ID, nil)
	if errors.Is(err, wrappers.NonExistentErr) {
		err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
	}

	return err
}
