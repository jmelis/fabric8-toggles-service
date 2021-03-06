package main

import (
	"context"
	"fmt"
	"net/http"

	authmiddleware "github.com/fabric8-services/fabric8-auth/goamiddleware"
	"github.com/fabric8-services/fabric8-auth/log"
	"github.com/fabric8-services/fabric8-toggles-service/app"
	"github.com/fabric8-services/fabric8-toggles-service/auth"
	"github.com/fabric8-services/fabric8-toggles-service/configuration"
	"github.com/fabric8-services/fabric8-toggles-service/controller"
	"github.com/fabric8-services/fabric8-toggles-service/jsonapi"
	"github.com/fabric8-services/fabric8-toggles-service/token"
	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
)

func main() {

	// Initialized configuration
	config, err := configuration.GetData()
	if err != nil || config == nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to setup the configuration")
	}

	fmt.Printf("%s", config)

	// Create service
	service := goa.New("feature")

	// Initialize log
	log.InitializeLogger(config.IsLogJSON(), config.GetLogLevel())

	service.WithLogger(goalogrus.New(log.Logger()))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(gzip.Middleware(9))
	service.Use(jsonapi.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	c, err := auth.NewClient(context.Background(), config.GetAuthServiceURL())
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to initialize auth service client")
	}
	tokenParser, err := token.NewParser(c)
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to create token manager")
	}
	// Middleware that extracts and stores the token in the context
	jwtMiddlewareTokenContext := authmiddleware.TokenContext(tokenParser, app.NewJWTSecurity())
	service.Use(jwtMiddlewareTokenContext)
	app.UseJWTMiddleware(service, goajwt.New(tokenParser.PublicKeys(), nil, app.NewJWTSecurity()))
	service.Use(log.LogRequest(config.IsDeveloperModeEnabled()))

	// Mount "features" controller
	featuresCtrl := controller.NewFeaturesController(service, tokenParser, config)
	app.MountFeaturesController(service, featuresCtrl)

	// Mount "status" controller
	statusCtrl := controller.NewStatusController(service, config)
	app.MountStatusController(service, statusCtrl)

	log.Logger().Infoln("Git Commit SHA: ", controller.Commit)
	log.Logger().Infoln("UTC Build Time: ", controller.BuildTime)
	log.Logger().Infoln("UTC Start Time: ", controller.StartTime)
	log.Logger().Infoln("Dev mode:       ", config.IsDeveloperModeEnabled())

	http.Handle("/api/", service.Mux)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// Start http
	if err := http.ListenAndServe(config.GetHTTPAddress(), nil); err != nil {
		log.Error(nil, map[string]interface{}{
			"addr": config.GetHTTPAddress(),
			"err":  err,
		}, "unable to connect to server")
		service.LogError("startup", "err", err)
	}

}
