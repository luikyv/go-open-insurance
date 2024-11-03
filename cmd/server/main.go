package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/capitalizationtitle"
	"github.com/luikyv/go-open-insurance/internal/consent"
	"github.com/luikyv/go-open-insurance/internal/customer"
	"github.com/luikyv/go-open-insurance/internal/endorsement"
	"github.com/luikyv/go-open-insurance/internal/quoteauto"
	"github.com/luikyv/go-open-insurance/internal/resource"
	"github.com/luikyv/go-open-insurance/internal/user"
	"github.com/luikyv/go-open-insurance/internal/webhook"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbSchema           = getEnv("MOCKIN_DB_SCHEMA", "mockin")
	dbStringConnection = getEnv("MOCKIN_DB_CONNECTION", "mongodb://localhost:27017/mockin")
	port               = getEnv("MOCKIN_PORT", "80")
	host               = getEnv("MOCKIN_HOST", "https://mockin.local")
	mtlsHost           = getEnv("MOCKIN_MTLS_HOST", "https://matls-mockin.local")
	apiPrefixOIDC      = "/auth"
	apiPrefixOPIN      = "/open-insurance"
)

type ConsentServerV2 = consent.ServerV2
type CustomerServerV1 = customer.ServerV1
type ResourceServerV2 = resource.ServerV2
type CapitalizationTitleServerV1 = capitalizationtitle.ServerV1
type EndorsementServerV1 = endorsement.ServerV1
type QuoteAutoServerV1 = quoteauto.ServerV1
type opinServer struct {
	ConsentServerV2
	CustomerServerV1
	ResourceServerV2
	CapitalizationTitleServerV1
	EndorsementServerV1
	QuoteAutoServerV1
}

func main() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}

	userStorage := user.NewStorage()
	consentStorage := consent.NewStorage(db)

	userService := user.NewService(userStorage)
	consentService := consent.NewService(consentStorage, userService)

	// OpenID Provider.
	op, err := openidProvider(
		db,
		userService,
		consentService,
		host,
		mtlsHost,
		apiPrefixOIDC,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Storage.
	idempotencyStorage := api.NewIdempotencyStorage(db)
	resourceStorage := resource.NewStorage()
	customerStorage := customer.NewStorage()
	capitalizationtitleStorage := capitalizationtitle.NewStorage()
	quoteAutoStorage := quoteauto.NewStorage(db)

	// Services.
	webhookService := webhook.NewService(op, httpClientFunc())
	idempotencyService := api.NewIdempotencyService(idempotencyStorage)
	resourceService := resource.NewService(resourceStorage, consentService)
	customerService := customer.NewService(customerStorage)
	capitalizationtitleService := capitalizationtitle.NewService(
		capitalizationtitleStorage,
		resourceService,
	)
	endorsementService := endorsement.NewService(consentService, resourceService)
	quoteAutoService := quoteauto.NewService(quoteAutoStorage, webhookService)

	// Server.
	server := opinServer{
		ConsentServerV2:             consent.NewServerV2(consentService),
		CustomerServerV1:            customer.NewServerV1(customerService),
		ResourceServerV2:            resource.NewServerV2(resourceService),
		CapitalizationTitleServerV1: capitalizationtitle.NewServerV1(capitalizationtitleService),
		EndorsementServerV1:         endorsement.NewServerV1(endorsementService),
		QuoteAutoServerV1:           quoteauto.NewServerV1(quoteAutoService),
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		log.Fatal(err)
	}
	strictHandler := api.NewStrictHandlerWithOptions(
		server,
		[]nethttp.StrictHTTPMiddlewareFunc{
			api.MetaMiddleware(mtlsHost),
			api.CacheControlMiddleware(),
			api.AuthPermissionMiddleware(consentService),
			api.AuthScopeMiddleware(op),
			api.IdempotencyMiddleware(idempotencyService),
			api.FAPIIDMiddleware(),
		},
		api.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  api.RequestErrorMiddleware,
			ResponseErrorHandlerFunc: api.ResponseErrorMiddleware,
		},
	)

	opinHandler := api.HandlerFromMux(strictHandler, http.NewServeMux())
	opinHandler = api.SchemaValidationMiddleware(opinHandler, router)
	opinHandler = api.ResponseEncodingMiddleware(opinHandler)

	mux := http.NewServeMux()
	mux.Handle(apiPrefixOIDC+"/", op.Handler())
	mux.Handle(apiPrefixOPIN+"/", opinHandler)

	// Run.
	if err := loadMocks(
		userService,
		customerService,
		resourceService,
		capitalizationtitleService,
	); err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Handler: mux,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func dbConnection() (*mongo.Database, error) {
	conn, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(dbStringConnection).SetBSONOptions(&options.BSONOptions{
			UseJSONStructTags: true,
			NilMapAsEmpty:     true,
			NilSliceAsEmpty:   true,
		}),
	)
	if err != nil {
		return nil, err
	}

	return conn.Database(dbSchema), nil
}

// getEnv retrieves an environment variable or returns a fallback value if not found
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
