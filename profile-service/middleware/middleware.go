package middleware

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
)

func CreateOpenapiMiddleware(
	getSwaggers ...func() (*openapi3.T, error),
) (gin.HandlerFunc, error) {
	// สร้าง router จากแต่ละ spec
	routersList := make([]routers.Router, 0, len(getSwaggers))

	for _, getSwagger := range getSwaggers {
		spec, err := getSwagger()
		if err != nil {
			return nil, err
		}
		// Validate และสร้าง router
		if err := spec.Validate(context.Background()); err != nil {
			return nil, err
		}

		r, err := legacy.NewRouter(spec)
		if err != nil {
			return nil, err
		}
		routersList = append(routersList, r)
	}

	return func(c *gin.Context) {
		var matched bool

		for _, r := range routersList {
			route, pathParams, err := r.FindRoute(c.Request)
			if err != nil {
				continue
			}

			// Validate request input
			reqValidation := &openapi3filter.RequestValidationInput{
				Request:    c.Request,
				PathParams: pathParams,
				Route:      route,
			}

			if err := openapi3filter.ValidateRequest(c.Request.Context(), reqValidation); err == nil {
				// ผ่าน
				matched = true
				break
			}
		}

		if !matched {
			c.JSON(http.StatusNotFound, gin.H{"error": "no matching operation was found"})
			c.Abort()
			return
		}

		c.Next()
	}, nil
}