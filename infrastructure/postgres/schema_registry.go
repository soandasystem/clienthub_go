package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"clienthub_go/core/models"
	"clienthub_go/core/ports"

	"github.com/antoniomarfa/hexatools/infrastructure"
	"github.com/antoniomarfa/hexatools/wrappers"

	"gorm.io/gorm"
)

// userRepository adapter of an roles repository for postgres
type schemaRegistryRepository struct {
	infrastructure.PostgresRepositoryOrm
}

// NewUserRepository creates a roles repository for postgres
func NewSchemaRegistryRepository(ctx context.Context, db *gorm.DB) ports.SchemaRegistryRepository {
	return &schemaRegistryRepository{
		infrastructure.PostgresRepositoryOrm{
			DB: db,
		},
	}
}

func (s *schemaRegistryRepository) Create(ctx context.Context, schema interface{}) (string, error) {
	//asigna el schma a usar
	DB := infrastructure.GetDBWithSchema(ctx, s.DB)
	// Asegúrate de que el tipo del token es correcto
	u := schema.(models.TokenSchemaRegistry)

	var existingReg models.TokenSchemaRegistry

	err := DB.Where("token = ?", u.Token).First(&existingReg).Error
	if err == nil {
		// Si no hay error, significa que se encontró un token con ese nombre
		return "error", errors.New("El Token '" + u.Token + "' ya existe")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Si el error no es de registro no encontrado, es un error inesperado
		return "error", errors.New("Error al buscar Token: " + err.Error())
	}

	//asigna el schma a usar
	crDB := infrastructure.GetDBWithSchema(ctx, s.DB)
	// Usamos el contexto y creamos el registro en la base de datos
	if err := crDB.WithContext(ctx).Create(&u).Error; err != nil {
		return "", err
	}

	// El ID se asigna automáticamente a la estructura 'u' después de la creación
	return fmt.Sprintf("%d", u.Id), nil
}

func (s *schemaRegistryRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
	//asigna el schma a usar
	DB := infrastructure.GetDBWithSchema(ctx, s.DB)
	var registro []models.TokenSchemaRegistry
	// Crea una consulta base
	query := DB.WithContext(ctx).Model(&models.TokenSchemaRegistry{})

	// Aplica filtros si existen
	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Aplica paginación si 'skip' y 'take' están definidos
	if skip != nil && take != nil {
		query = query.Offset(*skip).Limit(*take)
	}

	// Ejecuta la consulta
	if err := query.Order("id asc").Find(&registro).Error; err != nil {
		return nil, err
	}

	if len(registro) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	// Mapear a la estructura de respuesta
	response := models.TokenSchemaRegistryResponse{
		Items:      registro,
		TotalCount: 0,
	}

	// Devolver como slice de interface{} para cumplir con la interfaz
	return []interface{}{response}, nil
}

func (s *schemaRegistryRepository) GetInf(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
	return nil, nil
}

func (s *schemaRegistryRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	DB := infrastructure.GetDBWithSchema(ctx, s.DB)
	var registro models.TokenSchemaRegistry
	if err := DB.WithContext(ctx).Where("id = ?", ID).First(&registro).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, wrappers.NewNonExistentErr(err)
		}
		return nil, wrappers.NewNonExistentErr(err)
	}
	return &registro, nil
}

func (s *schemaRegistryRepository) GetByIDUpdate(ctx context.Context, ID string) (interface{}, error) {
	return s.GetByID(ctx, ID)
}

func (s *schemaRegistryRepository) Update(ctx context.Context, ID string, data interface{}) error {
	DB := infrastructure.GetDBWithSchema(ctx, s.DB)
	if err := DB.WithContext(ctx).Model(&models.TokenSchemaRegistry{}).Where("id = ?", ID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *schemaRegistryRepository) Delete(ctx context.Context, ID string, filter map[string]interface{}) error {
	DB := infrastructure.GetDBWithSchema(ctx, s.DB)
	var registro models.TokenSchemaRegistry
	result := DB.WithContext(ctx).Where("id = ?", ID).Delete(&registro)
	if result.Error != nil {
		return wrappers.NewNonExistentErr(result.Error)
	}
	if result.RowsAffected == 0 {
		return wrappers.NewNonExistentErr(sql.ErrNoRows)
	}
	return nil
}
