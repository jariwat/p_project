package middleware

import (
	"net/http"

	oapigin "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

func CreateOpenapiMiddleware(
	getSwaggers ...func() (*openapi3.T, error),
) (gin.HandlerFunc, error) {
	// สร้าง slice ของ middleware ที่ validate ตามแต่ละ spec
	validators := make([]gin.HandlerFunc, 0, len(getSwaggers))

	for _, getSwagger := range getSwaggers {
		spec, err := getSwagger()
		if err != nil {
			return nil, err
		}

		validator := oapigin.OapiRequestValidatorWithOptions(spec, &oapigin.Options{})

		validators = append(validators, validator)
	}

	// return MiddlewareFunc
	return func(c *gin.Context) {
		var matched bool

		for _, v := range validators {
			// ใช้ middleware ตรวจแบบ dummy
			copyCtx := c.Copy()
			v(copyCtx)

			// ตรวจว่า request ผ่าน spec ไหน
			if !copyCtx.IsAborted() {
				matched = true
				break
			}

			// gin ไม่มี error return จาก middleware แบบ echo
			// ดังนั้นควรใช้ StatusCode หรือ Aborted Flag ตรวจเอา
		}

		if !matched {
			c.JSON(http.StatusNotFound, gin.H{"error": "no matching operation was found"})
			c.Abort()
			return
		}

		c.Next()
	}, nil
}
