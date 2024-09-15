package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/luikyv/go-open-insurance/internal/api"
	"github.com/luikyv/go-open-insurance/internal/consent"
	consentv2 "github.com/luikyv/go-open-insurance/internal/consent/v2"
	"github.com/luikyv/go-open-insurance/internal/user"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Error middleware, validation middleware.

const (
	databaseSchema           = "gopin"
	databaseStringConnection = "mongodb://admin:password@localhost:27018"
	port                     = "80"
	host                     = "https://gopin.localhost"
	mtlsHost                 = "https://matls-gopin.localhost"
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
	consentStorage := consent.NewStorage(db)

	// Services.
	userService := user.NewService(userStorage)
	consentService := consent.NewService(userService, consentStorage)

	// Provider.
	op, err := openidProvider(db, userService, consentService,
		host, mtlsHost, apiPrefixOIDC)
	if err != nil {
		log.Fatal(err)
	}

	// Server.
	consentV2Server := consentv2.NewServer(baseURLOPIN, consentService)
	server := opinServer{
		Server: consentV2Server,
	}
	strictHandler := api.NewStrictHandlerWithOptions(
		server,
		[]nethttp.StrictHTTPMiddlewareFunc{
			api.AuthScopeMiddleware(op),
			api.FAPIIDMiddleware(),
			api.CacheControlMiddleware(),
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
	opinHandler = http.StripPrefix(apiPrefixOPIN, opinHandler)

	mux := http.NewServeMux()
	mux.Handle(apiPrefixOIDC+"/", op.Handler())
	mux.Handle(apiPrefixOPIN+"/", opinHandler)

	// Run.
	if err := loadUsers(userService); err != nil {
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

type opinServer struct {
	consentv2.Server
}
