package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Visit struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID       int64          `gorm:"not null" json:"seller_id"`
	CustomerType   string         `gorm:"type:varchar(20);not null" json:"customer_type"`
	CRMCustomerID  *string        `gorm:"type:varchar(100)" json:"crm_customer_id"`
	ProspectID     *int64         `json:"prospect_id"`
	ScheduledAt    *time.Time     `json:"scheduled_at"`
	Status         string         `gorm:"type:varchar(20);not null;default:PLANNED" json:"status"`
	StartedAt      *time.Time     `json:"started_at"`
	StartLatitude  *float64       `gorm:"type:numeric(10,7)" json:"start_latitude"`
	StartLongitude *float64       `gorm:"type:numeric(10,7)" json:"start_longitude"`
	StartAddress   *string        `gorm:"type:text" json:"start_address"`
	EndedAt        *time.Time     `json:"ended_at"`
	EndLatitude    *float64       `gorm:"type:numeric(10,7)" json:"end_latitude"`
	EndLongitude   *float64       `gorm:"type:numeric(10,7)" json:"end_longitude"`
	EndAddress     *string        `gorm:"type:text" json:"end_address"`
	DeviceInfo     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"device_info"`
	FinalComments  *string        `gorm:"type:text" json:"final_comments"`
	CreatedBy      int64          `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null;default:now()" json:"updated_at"`
}

func (Visit) TableName() string {
	return "salesvisit.visits"
}

type VisitListResponse struct {
	Items      []Visit `json:"items"`
	TotalCount int64   `json:"totalCount"`
}

type VisitDetail struct {
	ID                int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitID           int64          `gorm:"not null;unique" json:"visit_id"`
	VisitObjective    *string        `gorm:"type:text" json:"visit_objective"`
	TopicsDiscussed   *string        `gorm:"type:text" json:"topics_discussed"`
	ProductsConsulted datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"products_consulted"`
	Competitors       *string        `gorm:"type:text" json:"competitors"`
	Complaints        *string        `gorm:"type:text" json:"complaints"`
	InterestLevel     *string        `gorm:"type:varchar(30)" json:"interest_level"`
	Result            *string        `gorm:"type:varchar(30);column:result" json:"result"`
	Observations      *string        `gorm:"type:text" json:"observations"`
	CreatedAt         time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"not null;default:now()" json:"updated_at"`
}

func (VisitDetail) TableName() string {
	return "salesvisit.visit_details"
}

type VisitDetailListResponse struct {
	Items      []VisitDetail `json:"items"`
	TotalCount int64         `json:"totalCount"`
}

type VisitFollowup struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitID      int64      `gorm:"not null" json:"visit_id"`
	SellerID     int64      `gorm:"not null" json:"seller_id"`
	ActivityType string     `gorm:"type:varchar(30);not null" json:"activity_type"`
	Description  *string    `gorm:"type:text" json:"description"`
	ScheduledAt  time.Time  `gorm:"not null" json:"scheduled_at"`
	Status       string     `gorm:"type:varchar(20);not null;default:PENDING" json:"status"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

func (VisitFollowup) TableName() string {
	return "salesvisit.visit_followups"
}

type VisitFollowuptListResponse struct {
	Items      []VisitFollowup `json:"items"`
	TotalCount int64           `json:"totalCount"`
}

type VisitProduct struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitID             int64     `gorm:"not null" json:"visit_id"`
	CRMProductID        string    `gorm:"type:varchar(100);not null" json:"crm_product_id"`
	ProductNameSnapshot *string   `gorm:"type:varchar(200)" json:"product_name_snapshot"`
	Notes               *string   `gorm:"type:text" json:"notes"`
	CreatedAt           time.Time `gorm:"not null;default:now()" json:"created_at"`
}

func (VisitProduct) TableName() string {
	return "salesvisit.visit_products"
}

type VisitProductListResponse struct {
	Items      []VisitProduct `json:"items"`
	TotalCount int64          `json:"totalCount"`
}

// Sales
type SellerSupervisor struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SupervisorID int64      `gorm:"not null" json:"supervisor_id"`
	SellerID     int64      `gorm:"not null" json:"seller_id"`
	ValidFrom    time.Time  `gorm:"type:date;not null;default:CURRENT_DATE" json:"valid_from"`
	ValidTo      *time.Time `gorm:"type:date" json:"valid_to"`
}

func (SellerSupervisor) TableName() string {
	return "salesvisit.seller_supervisors"
}

type SellerSupervisorListResponse struct {
	Items      []SellerSupervisor `json:"items"`
	TotalCount int64              `json:"totalCount"`
}

type SalesOrder struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber   string     `gorm:"type:varchar(50);not null;unique" json:"order_number"`
	VisitID       int64      `gorm:"not null" json:"visit_id"`
	SellerID      int64      `gorm:"not null" json:"seller_id"`
	CustomerType  string     `gorm:"type:varchar(20);not null" json:"customer_type"`
	CRMCustomerID *string    `gorm:"type:varchar(100)" json:"crm_customer_id"`
	ProspectID    *int64     `json:"prospect_id"`
	OrderDate     time.Time  `gorm:"not null;default:now()" json:"order_date"`
	Currency      string     `gorm:"type:varchar(3);not null;default:CLP" json:"currency"`
	Subtotal      float64    `gorm:"type:numeric(14,2);not null;default:0" json:"subtotal"`
	DiscountTotal float64    `gorm:"type:numeric(14,2);not null;default:0" json:"discount_total"`
	TaxTotal      float64    `gorm:"type:numeric(14,2);not null;default:0" json:"tax_total"`
	Total         float64    `gorm:"type:numeric(14,2);not null;default:0" json:"total"`
	Status        string     `gorm:"type:varchar(25);not null;default:DRAFT" json:"status"`
	Notes         *string    `gorm:"type:text" json:"notes"`
	CRMOrderID    *string    `gorm:"type:varchar(100)" json:"crm_order_id"`
	SentToCRMAt   *time.Time `json:"sent_to_crm_at"`
	SyncStatus    string     `gorm:"type:varchar(20);not null;default:PENDING" json:"sync_status"`
	CreatedAt     time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	// Relaciones
	Details []SalesOrderDetail `gorm:"foreignKey:OrderID;references:ID" json:"details,omitempty"`
}

func (SalesOrder) TableName() string {
	return "salesvisit.sales_orders"
}

type SalesOrderDetail struct {
	ID              int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID         int64   `gorm:"not null" json:"order_id"`
	CRMProductID    string  `gorm:"type:varchar(100);not null" json:"crm_product_id"`
	ProductCode     *string `gorm:"type:varchar(80)" json:"product_code"`
	ProductName     string  `gorm:"type:varchar(200);not null" json:"product_name"`
	Quantity        float64 `gorm:"type:numeric(14,3);not null" json:"quantity"`
	UnitPrice       float64 `gorm:"type:numeric(14,2);not null" json:"unit_price"`
	DiscountPercent float64 `gorm:"type:numeric(7,4);not null;default:0" json:"discount_percent"`
	DiscountAmount  float64 `gorm:"type:numeric(14,2);not null;default:0" json:"discount_amount"`
	TaxRate         float64 `gorm:"type:numeric(7,4);not null;default:0" json:"tax_rate"`
	LineTotal       float64 `gorm:"type:numeric(14,2);not null" json:"line_total"`
	Notes           *string `gorm:"type:text" json:"notes"`
}

func (SalesOrderDetail) TableName() string {
	return "salesvisit.sales_order_details"
}

type SalesOrderListResponse struct {
	Items      []SalesOrder `json:"items"`
	TotalCount int64        `json:"totalCount"`
}

// Quotes
type Quote struct {
	ID            int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	QuoteNumber   string          `gorm:"type:varchar(50);not null;unique" json:"quote_number"`
	VisitID       int64           `gorm:"not null" json:"visit_id"`
	SellerID      int64           `gorm:"not null" json:"seller_id"`
	CustomerType  string          `gorm:"type:varchar(20);not null" json:"customer_type"`
	CRMCustomerID *string         `gorm:"type:varchar(100)" json:"crm_customer_id"`
	ProspectID    *int64          `json:"prospect_id"`
	QuoteDate     time.Time       `gorm:"not null;default:now()" json:"quote_date"`
	ValidUntil    *time.Time      `gorm:"type:date" json:"valid_until"`
	Currency      string          `gorm:"type:varchar(3);not null;default:CLP" json:"currency"`
	Subtotal      decimal.Decimal `gorm:"type:numeric(14,2);not null;default:0" json:"subtotal"`
	DiscountTotal decimal.Decimal `gorm:"type:numeric(14,2);not null;default:0" json:"discount_total"`
	TaxTotal      decimal.Decimal `gorm:"type:numeric(14,2);not null;default:0" json:"tax_total"`
	Total         decimal.Decimal `gorm:"type:numeric(14,2);not null;default:0" json:"total"`
	Status        string          `gorm:"type:varchar(20);not null;default:DRAFT" json:"status"`
	Notes         *string         `gorm:"type:text" json:"notes"`
	PDFPath       *string         `gorm:"type:text" json:"pdf_path"`
	CRMQuoteID    *string         `gorm:"type:varchar(100)" json:"crm_quote_id"`
	CreatedAt     time.Time       `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"not null;default:now()" json:"updated_at"`
	// Relaciones
	Details []QuoteDetail `gorm:"foreignKey:QuoteID;references:ID" json:"details,omitempty"`
}

func (Quote) TableName() string {
	return "salesvisit.quotes"
}

type QuoteDetail struct {
	ID              int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	QuoteID         int64           `gorm:"not null" json:"quote_id"`
	CRMProductID    string          `gorm:"type:varchar(100);not null" json:"crm_product_id"`
	ProductCode     *string         `gorm:"type:varchar(80)" json:"product_code"`
	ProductName     string          `gorm:"type:varchar(200);not null" json:"product_name"`
	Quantity        decimal.Decimal `gorm:"type:numeric(14,3);not null" json:"quantity"`
	UnitPrice       decimal.Decimal `gorm:"type:numeric(14,2);not null" json:"unit_price"`
	DiscountPercent decimal.Decimal `gorm:"type:numeric(7,4);not null;default:0" json:"discount_percent"`
	DiscountAmount  decimal.Decimal `gorm:"type:numeric(14,2);not null;default:0" json:"discount_amount"`
	TaxRate         decimal.Decimal `gorm:"type:numeric(7,4);not null;default:0" json:"tax_rate"`
	LineTotal       decimal.Decimal `gorm:"type:numeric(14,2);not null" json:"line_total"`
	Notes           *string         `gorm:"type:text" json:"notes"`
}

func (QuoteDetail) TableName() string {
	return "salesvisit.quote_details"
}

type QuoteListResponse struct {
	Items      []Quote `json:"items"`
	TotalCount int64   `json:"totalCount"`
}

// clients
type Prospect struct {
	ID                     int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyName            string    `gorm:"type:varchar(200);not null" json:"company_name"`
	ContactName            *string   `gorm:"type:varchar(150)" json:"contact_name"`
	Phone                  *string   `gorm:"type:varchar(40)" json:"phone"`
	Email                  *string   `gorm:"type:varchar(150)" json:"email"`
	Address                *string   `gorm:"type:text" json:"address"`
	City                   *string   `gorm:"type:varchar(100)" json:"city"`
	BusinessActivity       *string   `gorm:"type:varchar(150)" json:"business_activity"`
	Status                 string    `gorm:"type:varchar(30);not null;default:NEW" json:"status"`
	AssignedSellerID       *int64    `json:"assigned_seller_id"`
	ConvertedCRMCustomerID *string   `gorm:"type:varchar(100)" json:"converted_crm_customer_id"`
	CreatedBy              int64     `gorm:"not null" json:"created_by"`
	CreatedAt              time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt              time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (Prospect) TableName() string {
	return "salesvisit.prospects"
}

type ProspectListResponse struct {
	Items      []Prospect `json:"items"`
	TotalCount int64      `json:"totalCount"`
}

// registros externos
type ExternalRecord struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConnectionID    int64      `gorm:"not null" json:"connection_id"`
	EntityType      string     `gorm:"type:varchar(50);not null" json:"entity_type"`
	ExternalID      string     `gorm:"type:varchar(150);not null" json:"external_id"`
	LocalEntityType string     `gorm:"type:varchar(50);not null" json:"local_entity_type"`
	LocalEntityID   int64      `gorm:"not null" json:"local_entity_id"`
	LastSyncedAt    *time.Time `json:"last_synced_at"`
}

func (ExternalRecord) TableName() string {
	return "salesvisit.external_records"
}

type ExternalRecordListResponse struct {
	Items      []ExternalRecord `json:"items"`
	TotalCount int64            `json:"totalCount"`
}
