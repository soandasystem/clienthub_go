package handlers

import (
	"clienthub_go/config"
	"clienthub_go/core/models"
	"clienthub_go/core/ports"
	"clienthub_go/util"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetUserRoutes creates sale routes
func SetSchemaRegistryRoutes(ctx context.Context, cfg config.Config, r *gin.Engine, p ports.SchemaRegistryService) {
	//	r.Use(middlewares.Recover())
	authGroup := r.Group("api/v2027_01/")
	{
		authGroup.POST("schema_registry", createSchema(p))
		authGroup.GET("schema_registry", getAllSchema(p))
	}

}

// @Summary Create sale
// @Description Creates a new sale
// @Tags sale
// @Param user body models.CreateSaleReq true "New sale to be created"
// @Success 201 {object} models.SaleResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/sale [post]
func createSchema(p ports.SchemaRegistryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El ctx ya tiene timeout + schema
		ctx := c.Request.Context()

		// Leer el cuerpo de la solicitud
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}

		// Parsear el JSON al modelo adecuado
		var schema models.TokenSchemaRegistry
		err = json.Unmarshal(body, &schema)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}

		// Llamar al servicio para crear los sale
		id, err := p.Create(ctx, schema)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}

		// Devolver la respuesta con el código de estado adecuado
		result := util.NewMessageResponse("Registro Creado con id: "+id, id)

		response := util.NewSuccessResponse(result, http.StatusAccepted)
		c.JSON(response.StatusCode, response)
	}
}

// @Summary Get all schema
// @Description Gets all the schema
// @Tags schema
// @Success 200 {array} models.TokenSchemaRegistry "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/sale [get]
func getAllSchema(p ports.SchemaRegistryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El ctx ya tiene timeout + schema
		ctx := c.Request.Context()

		filter := make(map[string]interface{})

		// Obtener todos los parámetros de consulta
		queryParams := c.Request.URL.Query()
		// Convertir a un map[string]interface{}
		for key, values := range queryParams {
			if len(values) > 1 {
				// Si hay múltiples valores, los almacenamos como slice
				filter[key] = values
			} else {
				// Si hay un solo valor, lo almacenamos directamente
				filter[key] = values[0]
			}
		}

		result, err := p.GetAll(ctx, filter)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}
		// Agregar el total count al header (fuera del wrapper JSON)
		c.Header("X-Total-Count", strconv.FormatInt(result.TotalCount, 10))

		// Enviar solo los items en el wrapper de éxito
		response := util.NewSuccessResponse(result.Items, http.StatusOK)
		c.JSON(response.StatusCode, response)
	}
}
