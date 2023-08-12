package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"

	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	tmodels "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
)

type AuthHandler func(force bool, scope []string, perms []string) fiber.Handler

type AuthMiddleware struct {
	conn *subapps.HeLiCoPtErConnection

	Fn AuthHandler
}

type AuthConfig struct {
	Next        func(c *fiber.Ctx) bool
	LookupToken string
}

func NewAuth(c *subapps.HeLiCoPtErConnection) *AuthMiddleware {
	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: authorization",
	}

	inst := &AuthMiddleware{conn: c}
	inst.Fn = func(force bool, scope []string, perms []string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			if cfg.Next != nil && cfg.Next(c) {
				return c.Next()
			}

			err := inst.LookupAuthToken(c, strings.Split(cfg.LookupToken, ","))
			if err != nil && force {
				return fiber.NewError(fiber.StatusUnauthorized, err.Error())
			} else {
				if err == nil {
					if err := c.Locals("principal-session").(tmodels.UserSession).HasScope(scope...); err != nil {
						return fiber.NewError(fiber.StatusForbidden, err.Error())
					} else if err := c.Locals("principal").(tmodels.User).HasPermissions(perms...); err != nil {
						return fiber.NewError(fiber.StatusForbidden, err.Error())
					}
				}

				c.Locals("principal-ok", err == nil)

			}

			return c.Next()
		}
	}

	return inst
}

func (v *AuthMiddleware) LookupAuthToken(c *fiber.Ctx, args []string) error {
	var str string
	for _, arg := range args {
		parts := strings.Split(strings.TrimSpace(arg), ":")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "header":
			if len(c.GetReqHeaders()[v]) > 0 {
				str = strings.TrimSpace(strings.ReplaceAll(c.GetReqHeaders()[v], "Bearer", ""))
			}
		case "query":
			if len(c.Query(v)) > 0 {
				str = c.Query(v)
			}
		case "cookie":
			if len(c.Cookies(v)) > 0 {
				str = c.Cookies(v)
			}
		}
	}

	if len(str) == 0 {
		return fmt.Errorf("missing token in request")
	}

	resp, err := v.conn.GetAccount(str)
	if err != nil {
		return fmt.Errorf("failed to read details: %q", err)
	}

	c.Locals("principal", resp.User)
	c.Locals("principal-claims", resp.Claims)
	c.Locals("principal-session", resp.Session)

	return nil
}
