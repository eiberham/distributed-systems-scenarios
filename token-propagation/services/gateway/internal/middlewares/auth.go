package middlewares

import (
	"io"
	"strings"
	"token-propagation/internal/config"

	"github.com/labstack/echo/v5"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		// Inject client_id and client_secret into the request body before proxying to Keycloak
		c.Request().ParseForm()

		url := c.Request().URL.Path

		if strings.HasSuffix(url, "/login") {
			c.Request().Form.Set("grant_type", "password")
		} else if strings.HasSuffix(url, "/refresh") {
			c.Request().Form.Set("grant_type", "refresh_token")
		} else {
			return echo.NewHTTPError(400, "Invalid auth endpoint")
		}

		c.Request().Form.Set("client_id", config.Cfg.KeycloakClientID)
		c.Request().Form.Set("client_secret", config.Cfg.KeycloakSecretKey)

		// Update the request body with the modified form data
		body := c.Request().Form.Encode()
		c.Request().Body = io.NopCloser(strings.NewReader(body))
		c.Request().Header.Set("Content-Type", "application/x-www-form-urlencoded")
		c.Request().ContentLength = int64(len(body))
		return next(c)
	}
}
