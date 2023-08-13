package middleware

import (
	"code.smartsheep.studio/atom/matrix/pkg/server/datasource/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"

	"code.smartsheep.studio/atom/bedrock/pkg/kit/subapps"
	tmodels "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
)

type AuthHandler func(force bool, scope []string, perms []string) fiber.Handler

type AuthMiddleware struct {
	db   *gorm.DB
	conn *subapps.HeLiCoPtErConnection

	Fn AuthHandler
}

type AuthConfig struct {
	Next        func(c *fiber.Ctx) bool
	LookupToken string
}

func NewAuth(db *gorm.DB, c *subapps.HeLiCoPtErConnection) *AuthMiddleware {
	cfg := AuthConfig{
		Next:        nil,
		LookupToken: "header: Authorization, query: token, cookie: authorization",
	}

	inst := &AuthMiddleware{db, c, nil}
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

					var account *models.Account
					if err := inst.db.Where("user_id = ?", c.Locals("principal").(tmodels.User).ID).First(&account).Error; err != nil {
						if errors.Is(gorm.ErrRecordNotFound, err) {
							account = &models.Account{
								Nickname: c.Locals("principal").(tmodels.User).Nickname,
								UserID:   c.Locals("principal").(tmodels.User).ID,
							}

							if err := db.Save(&account).Error; err != nil {
								return fiber.NewError(fiber.StatusInternalServerError, err.Error())
							}
						} else {
							return fiber.NewError(fiber.StatusInternalServerError, err.Error())
						}
					}

					c.Locals("matrix-id", account)
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
