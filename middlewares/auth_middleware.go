package middleware

import (
	"fmt"
	"net/http"

	helper "appadming/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, err1 := c.Cookie("access_token")
		if err1 != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("uid", claims.Uid)
		c.Set("org_id", claims.Organization_id)

		c.Next()

	}
}
