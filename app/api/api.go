package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"clienthub_go/app/handlers"
	"clienthub_go/config"
	"clienthub_go/core/ports"
	"clienthub_go/core/services"

	//_ "clienthub_go/docs"
	"clienthub_go/infrastructure/b2"
	"clienthub_go/infrastructure/postgres"
	"clienthub_go/util"

	"github.com/antoniomarfa/hexatools/infrastructure/middleware"

	"github.com/antoniomarfa/hexatools/infrastructure"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type api struct {
	config   config.Config
	services svs
}

type svs struct {
	company        ports.CompanyService
	roles          ports.RolesService
	users          ports.UsersService
	schemaRegistry ports.SchemaRegistryService
	upload         ports.UploadService
	visit          ports.VisitService
	visitdetail    ports.VisitDetailService
}

// New creates a new API
func New(ctx context.Context, cfg config.Config) (a api) {
	a.config = cfg

	var companyRepo ports.CompanyRepository
	var rolesRepo ports.RolesRepository
	var usersRepo ports.UsersRepository
	var schemaRegistryRepo ports.SchemaRegistryRepository
	var visitRepo ports.VisitRepository
	var visitdetailRepo ports.VisitDetailRepository

	//abre la base de datos
	db, err := infrastructure.ConnectPostgresOrm(ctx, a.config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	companyRepo = postgres.NewCompanyRepository(ctx, db)
	rolesRepo = postgres.NewRolesRepository(ctx, db)
	usersRepo = postgres.NewUsersRepository(ctx, db)
	schemaRegistryRepo = postgres.NewSchemaRegistryRepository(ctx, db)
	visitRepo = postgres.NewVisitRepository(ctx, db)
	visitdetailRepo = postgres.NewVisitDetailRepository(ctx, db)

	a.services.company = services.NewCompanyService(a.config, companyRepo)
	a.services.roles = services.NewRolesService(a.config, rolesRepo)
	a.services.users = services.NewUsersService(a.config, usersRepo)
	a.services.schemaRegistry = services.NewSchemaRegistryService(a.config, schemaRegistryRepo)
	a.services.visit = services.NewVisitService(a.config, visitRepo)
	a.services.visitdetail = services.NewVisitDetailService(a.config, visitdetailRepo)

	// Inicializar B2 Storage y Upload Service
	var b2Storage ports.UploadStorage
	if a.config.B2Endpoint != "" {
		b2Storage = b2.NewB2Storage(
			ctx,
			a.config.B2KeyID,
			a.config.B2Application,
			a.config.B2Bucket,
			a.config.B2Region,
			a.config.B2Endpoint,
			a.config.B2PublicURL,
		)
		a.services.upload = services.NewUploadService(a.config, b2Storage)
	}

	//	a.services.contrato = services.NewContratoService(a.config, b2Storage)

	return a
}

// Run API
func (a *api) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		// Configurar CORS globalmente
		router.Use(middleware.CORSMiddleware())

		//Habilitar respuesta 505 method Not allowed
		router.HandleMethodNotAllowed = true
		router.NoMethod(func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Método HTTP no permitido para esta ruta",
			})
		})

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Ruta no encontrada",
			})
		})

		//asigna el schema a trabajar
		router.Use(util.SchemaMiddleware())

		// Luego este que inyecta el contexto con timeout y schema:
		router.Use(util.ContextWithSchemaMiddleware(a.config.Timeout.Duration))

		// ==========================================
		// 1. RUTAS PÚBLICAS (No requieren Token)
		// ==========================================
		handlers.SetHealthRoutes(ctx, a.config, router)
		handlers.SetCompanyRoutes(ctx, a.config, router, a.services.company)
		handlers.SetAuthRoutes(ctx, a.config, router, a.services.company, a.services.users)
		handlers.SetEmailRoutes(ctx, a.config, router)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// ==========================================
		// 2. RUTAS PRIVADAS (Requieren Token)
		// ==========================================
		// Al usar Use() aquí, solo afectará a las rutas registradas debajo
		router.Use(util.JWTAuthMiddleware())

		//handlers.SetCompanyRoutes(ctx, a.config, router, a.services.company)
		handlers.SetRolesRoutes(ctx, a.config, router, a.services.roles)
		handlers.SetUsersRoutes(ctx, a.config, router, a.services.users)
		handlers.SetSchemaRegistryRoutes(ctx, a.config, router, a.services.schemaRegistry)
		handlers.SetVisitRoutes(ctx, a.config, router, a.services.visit)
		handlers.SetVisitDetailRoutes(ctx, a.config, router, a.services.visitdetail)

		if a.services.upload != nil {
			handlers.SetUploadRoutes(ctx, a.config, router, a.services.upload)
		}

		log.Printf("Version: %s", a.config.Version)
		log.Printf("Environment: %s", a.config.Environment)
		log.Printf("Database: %s", a.config.Database)
		log.Printf("Listening on port %d", a.config.Port)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.Port),
			Handler: router,
		}
		go shutdown(ctx, server)
		return server.ListenAndServe()
	}
}

func SchemaMiddleware() gin.HandlerFunc {
	panic("unimplemented")
}

func shutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	log.Printf("Shutting down API gracefully...")
	server.Shutdown(ctx)
}
