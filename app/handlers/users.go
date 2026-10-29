package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"clienthub_go/config"
	"clienthub_go/core/models"
	"clienthub_go/core/ports"
	"clienthub_go/util"

	"github.com/gin-gonic/gin"
)

// SetUsersRoutes creates users routes
func SetUsersRoutes(ctx context.Context, cfg config.Config, r *gin.Engine, p ports.UsersService) {
	//	r.Use(middlewares.Recover())
	authGroup := r.Group("api/v2027_01/")
	{
		authGroup.POST("users", createUsers(p))
		authGroup.GET("users", getAllUsers(p))
		authGroup.GET("users/:id", getUsersByID(p))
		authGroup.PATCH("users/:id", updateUsers(p))
		authGroup.DELETE("users/:id", deleteUsers(p))
	}
}

// @Summary Create users
// @Description Creates a new users
// @Tags users
// @Param user body models.CreateUsersReq true "New users to be created"
// @Success 201 {object} models.UsersResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/users [post]
func createUsers(p ports.UsersService) gin.HandlerFunc {
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
		var users models.CreateUsersReq
		err = json.Unmarshal(body, &users)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}

		// Llamar al servicio para crear los users
		id, err := p.Create(ctx, users)
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

// @Summary Get all users
// @Description Gets all the users
// @Tags users
// @Success 200 {array} models.UsersResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/users [get]
func getAllUsers(p ports.UsersService) gin.HandlerFunc {
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

// @Summary Get users by ID
// @Description Gets a users by ID
// @Tags users
// @Param id path string true "ID"
// @Success 200 {object} models.UsersResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/users/{id} [get]
func getUsersByID(p ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El ctx ya tiene timeout + schema
		ctx := c.Request.Context()

		//	var params = mux.Vars(c.Request)
		result, err := p.GetByID(ctx, c.Param("id"))
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}
		//	utils.ResponseJSON(c.Writer, c.Request, nil, http.StatusAccepted, users)
		response := util.NewSuccessResponse(result, http.StatusOK)
		c.JSON(response.StatusCode, response)
	}
}

// @Summary Update users
// @Description Updates a users
// @Tags users
// @Param id path string true "ID"
// @Param User body models.UpdateUsersReq true "users"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/users/{id} [patch]
func updateUsers(p ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El ctx ya tiene timeout + schema
		ctx := c.Request.Context()

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}

		params := c.Params
		var users models.UpdateUsersReq
		err = json.Unmarshal(body, &users)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}
		id := params.ByName("id")
		err = p.Update(ctx, params.ByName("id"), users)
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}
		result := util.NewMessageResponse("Registro con id "+id+" actualizado", id)

		response := util.NewSuccessResponse(result, http.StatusAccepted)
		c.JSON(response.StatusCode, response)
	}
}

// @Summary Delete users
// @Description Delete a users
// @Tags users
// @Param id path string true "ID"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Security Bearer
// @Router /api/v2027_01/users/{id} [delete]
func deleteUsers(p ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// El ctx ya tiene timeout + schema
		ctx := c.Request.Context()

		//	var params = mux.Vars(c.Request)
		id := c.Param("id")
		err := p.Delete(ctx, c.Param("id"))
		if err != nil {
			response := util.NewErrorResponse(err, http.StatusInternalServerError)
			c.JSON(response.StatusCode, response)
			return
		}
		result := util.NewMessageResponse("Registro con id "+id+" eliminado ", id)

		response := util.NewSuccessResponse(result, http.StatusAccepted)
		c.JSON(response.StatusCode, response)
	}
}
