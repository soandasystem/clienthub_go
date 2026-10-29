package models

import "time"

type TokenSchemaRegistry struct {
	Id         int64     `json:"id"`
	Token      string    `json:"token"`
	SchemaName string    `json:"schema_name"`
	CompanyId  int64     `json:"company_id"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (TokenSchemaRegistry) TableName() string {
	return "global.token_schema_registry" // Nombre de la tabla con esquema en PostgreSQL
}


type TokenSchemaRegistryResponse struct {
	Items      []TokenSchemaRegistry `json:"items"`
	TotalCount int64                 `json:"totalCount"`
}
