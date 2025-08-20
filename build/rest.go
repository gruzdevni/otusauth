package build

import (
	"context"
	"fmt"
	"net/http"

	"otusauth/internal/restapi"
	"otusauth/internal/restapi/operations"
	"otusauth/internal/restapi/operations/other"
	"otusauth/internal/restapi/operations/user_c_r_u_d"
	"otusauth/internal/service/api/auth"
	"otusauth/internal/service/api/user"

	httpMW "otusauth/pkg/http"

	"github.com/go-openapi/loads"
	mdlwr "github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
)

func (b *Builder) buildAPI() (*operations.RestServerAPI, *loads.Document, error) {
	swaggerSpec, err := loads.Spec("api/swagger/file.yaml")
	if err != nil {
		return nil, nil, fmt.Errorf("load swagger specs: %w", err)
	}

	api := operations.NewRestServerAPI(swaggerSpec)

	psql, err := b.PostgresClient()
	if err != nil {
		return nil, nil, fmt.Errorf("creating postgres client: %w", err)
	}

	repo := b.NewRepo(psql.DB)

	userSrv := user.NewService(repo)
	authSrv := auth.NewService(repo)

	handler := restapi.NewHandler(userSrv, authSrv)

	api.OtherGetHealthHandler = other.GetHealthHandlerFunc(
		handler.GetHealth,
	)

	api.UsercrudGetUserGUIDHandler = user_c_r_u_d.GetUserGUIDHandlerFunc(
		handler.GetUser,
	)
	api.UsercrudDeleteUserGUIDHandler = user_c_r_u_d.DeleteUserGUIDHandlerFunc(
		handler.DeleteUser,
	)
	api.UsercrudPatchUserGUIDHandler = user_c_r_u_d.PatchUserGUIDHandlerFunc(
		handler.UpdateUser,
	)
	api.UsercrudPostUserHandler = user_c_r_u_d.PostUserHandlerFunc(
		handler.CreateUser,
	)
	api.OtherGetAuthHandler = other.GetAuthHandlerFunc(
		handler.Auth,
	)
	api.OtherPostLoginHandler = other.PostLoginHandlerFunc(
		handler.Login,
	)
	api.OtherPostSignupHandler = other.PostSignupHandlerFunc(
		handler.Signup,
	)

	return api, swaggerSpec, nil
}

//nolint:funlen
func (b *Builder) RestAPIServer(ctx context.Context) (*http.Server, error) {
	server, err := b.HTTPServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating http server: %w", err)
	}

	router := b.httpRouter()

	api, swaggerSpec, err := b.buildAPI()
	if err != nil {
		return nil, errors.Wrap(err, "building API")
	}

	apiEndpoint := swaggerSpec.BasePath()
	apiRouter := router.Name("api").Subrouter()

	metricsMW, err := NewRouter(b.config.App.Name, b.prometheusRegistry, router, api)
	if err != nil {
		return nil, fmt.Errorf("creating metrics middleware: %w", err)
	}

	authMW := httpMW.NewAuthMiddleware()

	apiRouter.Use(metricsMW)
	apiRouter.Use(httpMW.HTTPRequestBodyLoggerWithContext(ctx))
	apiRouter.Use(authMW.UserAuthorizationMiddleware)

	swaggerUIOpts := mdlwr.SwaggerUIOpts{ //nolint:exhaustruct
		BasePath: apiEndpoint,
		SpecURL:  fmt.Sprintf("%s/swagger.json", apiEndpoint),
	}

	apiRouter.PathPrefix(apiEndpoint).Handler(
		func() http.Handler {
			api.Init()

			return mdlwr.Spec(
				apiEndpoint,
				swaggerSpec.Raw(),
				mdlwr.SwaggerUI(
					swaggerUIOpts,
					api.Context().RoutesHandler(metricsMW),
				),
			)
		}(),
	)

	return server, nil
}
