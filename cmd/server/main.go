package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/luikyv/go-opf/internal/consent"
	"github.com/luikyv/go-opf/internal/middleware"
	"github.com/luikyv/go-opf/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	namespace                = "urn:gopin"
	databaseSchema           = "gopin"
	databaseStringConnection = "mongodb://admin:password@localhost:27018"
	port                     = ":80"
	host                     = "https://gopin.localhost"
	mtlsHost                 = "https://matls-gopin.localhost"
	apiPrefixOIDC            = "/auth"
	baseURLOIDC              = host + apiPrefixOIDC
	apiPrefixOPIN            = "/open-insurance"
	baseURLOPIN              = mtlsHost + apiPrefixOPIN
)

func main() {
	db := dbConnection()

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
		panic(err)
	}

	// Routers.
	consentRouter := consent.NewRouter(op, consentService, baseURLOPIN, namespace)

	// APIs.
	server := gin.Default()

	openBankingRouter := server.Group(apiPrefixOPIN)
	openBankingRouter.Use(middleware.CacheControl())
	openBankingRouter.Use(middleware.FAPIID())
	consentRouter.AddRoutesV2(openBankingRouter)

	oidcRouter := server.Group(apiPrefixOIDC)
	oidcRouter.Any("/*w", gin.WrapH(op.Handler()))

	// Run.
	if err := server.Run(port); err != nil {
		panic(err)
	}
}

func dbConnection() *mongo.Database {
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
		panic(err)
	}
	return conn.Database(databaseSchema)
}
