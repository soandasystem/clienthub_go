package handlers

import (
	"clienthub_go/config"
	"clienthub_go/core/models"
	"clienthub_go/core/ports"
	"clienthub_go/util"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetAuthRoutes define las rutas de autenticación públicas
// En este caso inyectamos los servicios necesarios para validar contra la BD
func SetAuthRoutes(ctx context.Context, cfg config.Config, router *gin.Engine, companyScv ports.CompanyService, userSvc ports.UsersService) {
	//authGroup := router.Group("/api/auth")
	authGroup := router.Group("api/v2027_01/")
	{
		authGroup.POST("/login", loginHandler(userSvc))
		authGroup.POST("/restore", restoreHandler(userSvc))
	}
}

// @Summary Iniciar sesión (Auth)
// @Description Permite iniciar sesión como usuario, apoderado o curso y devuelve un JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciales de login"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/v2027_01/login [post]
func loginHandler(userSvc ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := c.Request.Context() // Este contexto contiene el schema inyectado por el middleware

		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan datos o el formato es incorrecto"})
			return
		}

		var token string
		var err error

		// Dependiendo del login_type, realizamos la validación en la base de datos
		switch req.LoginType {
		case util.AuthTypeUser:
			if userSvc == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Servicio de usuarios no configurado"})
				return
			}
			// Buscar en la BD mediante el servicio
			filter := map[string]interface{}{
				"username": req.Username,
				"password_hash": req.Password, // Nota: en producción, no almacenes ni filtres por texto plano. Debes hashear.
			}
			usersList, errSvc := userSvc.GetAll(reqCtx, filter)
			if errSvc != nil || usersList == nil || len(usersList.Items) == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de usuario inválidas"})
				return
			}

			// Tomamos el primer usuario encontrado
			user := usersList.Items[0]
			roleIDStr := "1" //fmt.Sprintf("%d", user.RolesId)

			token, err = util.GenerateUserToken(user.ID, roleIDStr)

		case util.AuthTypeCourse:
			/*
				if cursoSvc == nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Servicio de cursos no configurado"})
					return
				}
				// Buscar en la BD
				filter := map[string]interface{}{
					"rutapod":  req.Rutapod,
					"password_hash": req.Password,
				}
				cursosList, errSvc := cursoSvc.GetAll(reqCtx, filter)
				if errSvc != nil || cursosList == nil || len(cursosList.Items) == 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de curso inválidas"})
					return
				}

				curso := cursosList.Items[0]
				token, err = util.GenerateCourseToken(curso.ID, curso.Rutapod)
			*/
		case util.AuthTypeAccess:
			/*
				if saleSvc == nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Servicio de ventas (sales) no configurado"})
					return
				}
				// Buscar la venta por accesscode
				filter := map[string]interface{}{
					"accesscode": req.AccessCode,
				}
				salesList, errSvc := saleSvc.GetAll(reqCtx, filter)
				if errSvc != nil || salesList == nil || len(salesList.Items) == 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Código de acceso inválido"})
					return
				}

				// Tomamos la primera coincidencia
				sale := salesList.Items[0]
				token, err = util.GenerateAccessCodeToken(sale.ID)
			*/
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de login no soportado. Debe ser 'user', 'course' o 'access_code'"})
			return
		}

		// Si ocurrió un error interno al firmar el token
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar credenciales"})
			return
		}

		// Devolver el token al frontend
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"type":  req.LoginType,
		})
	}
}

// @Summary Validar credenciales para restauración de contraseña
// @Description Valida la existencia del usuario o curso y retorna los datos para la recuperación
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciales para validación"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /api/v2027_01/restore [post]
func restoreHandler(userSvc ports.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := c.Request.Context() // Este contexto contiene el schema inyectado por el middleware

		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan datos o el formato es incorrecto"})
			return
		}

		var token string
		var emailStr string
		var userId string
		var err error

		// Dependiendo del login_type, realizamos la validación en la base de datos
		switch req.LoginType {
		case util.AuthTypeUser:
			if userSvc == nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Servicio de usuarios no configurado"})
				return
			}
			// Buscar en la BD mediante el servicio
			filter := map[string]interface{}{
				"username":      req.Username,
				"password_hash": req.Password,
			}
			usersList, errSvc := userSvc.GetAll(reqCtx, filter)
			fmt.Println("user list ", usersList)
			if errSvc != nil || usersList == nil || len(usersList.Items) == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de usuario inválidas"})
				return
			}

			// Tomamos el primer usuario encontrado
			user := usersList.Items[0]
			roleIDStr := "1" // fmt.Sprintf("%d",user.RolesId)
			emailStr = user.Email
			userId = user.ID
			token, err = util.GenerateUserToken(user.ID, roleIDStr)

		case util.AuthTypeCourse:
			/*
				if cursoSvc == nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Servicio de cursos no configurado"})
					return
				}
				// Buscar en la BD
				filter := map[string]interface{}{
					"rutapod": req.Rutapod,
				}
				cursosList, errSvc := cursoSvc.GetAll(reqCtx, filter)
				if errSvc != nil || cursosList == nil || len(cursosList.Items) == 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales de curso inválidas"})
					return
				}

				curso := cursosList.Items[0]
				userId = curso.ID
				emailStr = curso.Correo
				token, err = util.GenerateCourseToken(curso.ID, curso.Rutapod)
			*/
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de login no soportado. Debe ser 'user', 'course' o 'access_code'"})
			return
		}

		// Si ocurrió un error interno al firmar el token
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar credenciales"})
			return
		}

		// Devolver el token al frontend
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"type":  req.LoginType,
			"email": emailStr,
			"id":    userId,
		})
	}
}
