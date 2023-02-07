package handlers

import (
	"github.com/RunawayVPN/Runaway-Agent/tools/hub"
	"github.com/RunawayVPN/security"
	"github.com/gin-gonic/gin"
)

func AddConfig(c *gin.Context) {
	// Get payload from context

}

// Middleware for JWT authentication
func JWTAuth(hub_info hub.HubInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT from header
		token := c.Request.Header.Get("Authorization")
		// Verify JWT
		payload, err := security.VerifyToken(token, hub_info.PublicKey)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		// Add payload to context
		c.Set("payload", payload)
		// Continue
		c.Next()
	}
}
