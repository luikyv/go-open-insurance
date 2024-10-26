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
	capitalizationtitlev1 "github.com/luikyv/go-open-insurance/internal/capitalizationtitle/v1"
	"github.com/luikyv/go-open-insurance/internal/consent"
	consentv2 "github.com/luikyv/go-open-insurance/internal/consent/v2"
	"github.com/luikyv/go-open-insurance/internal/customer"
	customersv1 "github.com/luikyv/go-open-insurance/internal/customer/v1"
	"github.com/luikyv/go-open-insurance/internal/endorsement"
	endorsementv1 "github.com/luikyv/go-open-insurance/internal/endorsement/v1"
	"github.com/luikyv/go-open-insurance/internal/resource"
	resourcev2 "github.com/luikyv/go-open-insurance/internal/resource/v2"
	"github.com/luikyv/go-open-insurance/internal/user"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	databaseSchema           = getEnv("MOCKIN_DB_SCHEMA", "mockin")
	databaseStringConnection = getEnv("MOCKIN_DB_CONNECTION", "mongodb://localhost:27017/mockin")
	port                     = getEnv("MOCKIN_PORT", "80")
	host                     = getEnv("MOCKIN_HOST", "https://mockin.local")
	mtlsHost                 = getEnv("MOCKIN_MTLS_HOST", "https://matls-mockin.local")
	apiPrefixOIDC            = "/auth"
	apiPrefixOPIN            = "/open-insurance"
	baseURLOPIN              = mtlsHost + apiPrefixOPIN
)

func init() {

}

func main() {
	db, err := dbConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Storage.
	userStorage := user.NewStorage()
	idempotencyStorage := api.NewIdempotencyStorage(db)
	consentStorage := consent.NewStorage(db)
	resourceStorage := resource.NewStorage()
	customerStorage := customer.NewStorage()
	capitalizationtitleStorage := capitalizationtitle.NewStorage()

	// Services.
	userService := user.NewService(userStorage)
	idempotencyService := api.NewIdempotencyService(idempotencyStorage)
	consentService := consent.NewService(consentStorage, userService)
	resourceService := resource.NewService(resourceStorage, consentService)
	customerService := customer.NewService(customerStorage)
	capitalizationtitleService := capitalizationtitle.NewService(
		capitalizationtitleStorage,
		resourceService,
	)
	endorsementService := endorsement.NewService(consentService, resourceService)

	// Provider.
	op, err := openidProvider(db, userService, consentService,
		host, mtlsHost, apiPrefixOIDC)
	if err != nil {
		log.Fatal(err)
	}

	// Server.
	server := opinServer{
		consentV2Server:             consentv2.NewServer(baseURLOPIN, consentService),
		customerV1Server:            customersv1.NewServer(baseURLOPIN, customerService),
		resouceV2Server:             resourcev2.NewServer(baseURLOPIN, resourceService),
		capitalizationTitleV1Server: capitalizationtitlev1.NewServer(baseURLOPIN, capitalizationtitleService),
		endorsementV1Server:         endorsementv1.NewServer(endorsementService),
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
			api.MetaMiddleware(),
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
	opinHandler = http.StripPrefix(apiPrefixOPIN, opinHandler)

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
	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}
	clientOpts := options.Client().
		ApplyURI(databaseStringConnection).
		SetBSONOptions(bsonOpts)
	conn, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, err
	}
	return conn.Database(databaseSchema), nil
}

// getEnv retrieves an environment variable or returns a fallback value if not found
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
