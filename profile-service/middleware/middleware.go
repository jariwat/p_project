package middleware

import (
	"net/http"
	"strings"

	_middleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

func CreateOpenapiMiddleware(
	getSwaggers ...func() (*openapi3.T, error),
) (gin.HandlerFunc, error) {
	validators := make([]gin.HandlerFunc, 0, len(getSwaggers))

	for _, getSwagger := range getSwaggers {
		spec, err := getSwagger()
		if err != nil {
			return nil, err
		}

		validator := _middleware.OapiRequestValidatorWithOptions(spec, &_middleware.Options{
			Skipper: func(c _middleware.Context) bool {
				path := c.Request().URL.Path
				return path == "/" || path == "/health-check"
			},
			ErrorHandler: func(c _middleware.Context, err error) error {
				return &_middleware.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: err.Error(),
				}
			},
			SilenceServersWarning: true,
		})

		// Wrap Echo middleware as Gin middleware
		validators = append(validators, func(c *gin.Context) {
			var matched bool
			var epError error

			// สร้าง Echo context จำลองเพื่อให้ middleware ทำงานได้
			eCtx := _middleware.NewContextAdapter(c.Request, c.Writer)

			// ใช้ middleware ตรวจสอบ
			err := validator(func(ec _middleware.Context) error {
				matched = true
				return nil
			})(eCtx)

			epError = err

			if !matched && epError != nil {
				msg := epError.Error()
				if strings.Contains(msg, "no matching operation") {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "no matching operation was found"})
					return
				}
				if strings.Contains(msg, "parameter") || strings.Contains(msg, "request body") {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
					return
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
				return
			}

			// ตรวจสอบผ่าน
			c.Next()
		})
	}

	// รวม validator หลายตัวเป็น middleware เดียว
	return func(c *gin.Context) {
		for _, v := range validators {
			v(c)
			if c.IsAborted() {
				return
			}
		}
		c.Next()
	}, nil
}
