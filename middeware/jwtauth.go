package middeware

import (
	"context"
	"errors"

	"github.com/Sna-ken/hellogo/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func JWTAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(consts.StatusUnauthorized, errors.New("Token is empty"))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(token) //验证token
		if err != nil {
			c.JSON(consts.StatusUnauthorized, errors.New("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next(ctx)
	}
}
