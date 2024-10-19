package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/capitalizationtitle"
	capitalizationtitlev1 "github.com/luikyv/go-open-insurance/internal/capitalizationtitle/v1"
	"github.com/luikyv/go-open-insurance/internal/consent"
	consentv2 "github.com/luikyv/go-open-insurance/internal/consent/v2"
	"github.com/luikyv/go-open-insurance/internal/customer"
	customersv1 "github.com/luikyv/go-open-insurance/internal/customer/v1"
	"github.com/luikyv/go-open-insurance/internal/resource"
	resourcev2 "github.com/luikyv/go-open-insurance/internal/resource/v2"
	"github.com/luikyv/go-open-insurance/internal/user"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseSchema           = "gopin"
	databaseStringConnection = "mongodb://localhost:27017/gopin"
	port                     = "80"
	host                     = "https://gopin.local"
	mtlsHost                 = "https://matls-gopin.local"
	apiPrefixOIDC            = "/auth"
	baseURLOIDC              = host + apiPrefixOIDC
	apiPrefixOPIN            = "/open-insurance"
	baseURLOPIN              = mtlsHost + apiPrefixOPIN
)

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
	consentService := consent.NewService(consentStorage, userService)
	resourceService := resource.NewService(resourceStorage, consentService)
	customerService := customer.NewService(customerStorage)
	capitalizationtitleService := capitalizationtitle.NewService(
		capitalizationtitleStorage,
		resourceService,
	)
	idempotencyService := api.NewIdempotencyService(idempotencyStorage)

	// Provider.
	op, err := openidProvider(db, userService, consentService,
		host, mtlsHost, apiPrefixOIDC)
	if err != nil {
		log.Fatal(err)
	}

	// Server.
	server := opinServer{
		consentV2Server:  consentv2.NewServer(baseURLOPIN, consentService),
		customerV1Server: customersv1.NewServer(baseURLOPIN, customerService),
		resouceV2Server:  resourcev2.NewServer(baseURLOPIN, resourceService),
		capitalizationTitleV1Server: capitalizationtitlev1.NewServer(
			baseURLOPIN,
			capitalizationtitleService,
		),
	}
	strictHandler := api.NewStrictHandlerWithOptions(
		server,
		[]nethttp.StrictHTTPMiddlewareFunc{
			api.MetaMiddleware(),
			api.CacheControlMiddleware(),
			api.AuthPermissionMiddleware(consentService.VerifyPermissions),
			api.AuthScopeMiddleware(op),
			api.IdempotencyMiddleware(idempotencyService),
			api.FAPIIDMiddleware(),
		},
		api.StrictHTTPServerOptions{
			ResponseErrorHandlerFunc: api.ResponseErrorMiddleware,
		},
	)

	opinMux := http.NewServeMux()
	api.HandlerFromMux(strictHandler, opinMux)
	// Add a validation middleware for open insurance requests.
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	opinHandler := nethttpmiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&nethttpmiddleware.Options{
			ErrorHandler: api.ValidationErrorHandler(),
		},
	)(opinMux)
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
	log.Fatal(s.ListenAndServe())
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
