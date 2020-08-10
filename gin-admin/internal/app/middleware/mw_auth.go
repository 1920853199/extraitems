package middleware

import (
	"github.com/1920853199/extraitems/gin-admin/internal/app/config"
	"github.com/1920853199/extraitems/gin-admin/internal/app/ginplus"
	"github.com/1920853199/extraitems/gin-admin/internal/app/icontext"
	"github.com/1920853199/extraitems/gin-admin/pkg/auth"
	"github.com/1920853199/extraitems/gin-admin/pkg/errors"
	"github.com/1920853199/extraitems/gin-admin/pkg/logger"
	"github.com/gin-gonic/gin"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	ginplus.SetUserID(c, userID)
	ctx := icontext.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), ginplus.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, config.C.Root.UserName)
					c.Next()
					return
				}
				ginplus.ResError(c, errors.ErrInvalidToken)
				return
			}
			ginplus.ResError(c, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
