package util

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`    // Opcional: mensaje descriptivo
	Error      string      `json:"error,omitempty"`      // Error legible (si Success = false)
	StatusCode int         `json:"-"`                    // Usado internamente (no se serializa a JSON)
	Page       int         `json:"page,omitempty"`       // Número de página
	PageSize   int         `json:"pageSize,omitempty"`   // Tamaño de la página
	TotalPages int64       `json:"totalPages,omitempty"` // Total de páginas
	TotalCount int64       `json:"totalCount,omitempty"` // Total de páginas
	Data       interface{} `json:"data,omitempty"`       // Datos principales (ej: items)
}

type ApiMessage struct {
	ReturnID string `json:"return_id"`
	Message  string `json:"message,omitempty"`
}

func NewSuccessResponse(data interface{}, statusCode int) *ApiResponse {
	return &ApiResponse{
		Success:    true,
		Data:       data,
		StatusCode: statusCode,
	}
}

func NewSuccessResponsePage(totalCount int64, page int, pageSize int, totalPages int64, statusCode int, data interface{}) *ApiResponse {
	return &ApiResponse{
		Success:    true,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
		Data:       data,
		StatusCode: statusCode,
	}
}

func NewMessageResponse(message string, returnid string) *ApiMessage {
	return &ApiMessage{
		ReturnID: returnid,
		Message:  message,
	}
}

func NewErrorResponse(err error, statusCode int) *ApiResponse {
	return &ApiResponse{
		Success:    false,
		Error:      err.Error(),
		StatusCode: statusCode,
	}
}
