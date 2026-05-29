package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
	middleware "github.com/labstack/echo/v5/middleware"
)

func corsMiddleware(allowedOrigins []string) echo.MiddlewareFunc {

	return middleware.CORSWithConfig(

		middleware.CORSConfig{

			AllowOrigins: allowedOrigins,

			AllowMethods: []string{

				http.MethodGet,

				http.MethodPost,

				http.MethodPut,

				http.MethodPatch,

				http.MethodDelete,

				http.MethodOptions,
			},

			AllowHeaders: []string{

				echo.HeaderOrigin,

				echo.HeaderContentType,

				echo.HeaderAccept,

				echo.HeaderAuthorization,
			},
		},
	)

}
