package middleware

import (
	"net/http"
	"strings"

	"backend/internal/appcontext"
	"backend/internal/model"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, jwtSecret string, bl utils.BlacklistStore) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ---------- HEADER ----------
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		if bl != nil {
			blocked, _ := bl.IsBlacklisted(c.Request.Context(), tokenStr)
			if blocked {
				c.JSON(401, gin.H{"error": "token logged out"})
				c.Abort()
				return
			}
		}
		// ---------- PARSE TOKEN ----------
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := uint(claims["user_id"].(float64))
		tenantID := uint(claims["tenant_id"].(float64))
		role := claims["role"].(string)

		// ---------- LOAD USER ----------
		var user model.User
		if err := db.
			Where("id = ? AND tenant_id = ? AND is_active = true AND deleted_at IS NULL",
				userID, tenantID).
			First(&user).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			c.Abort()
			return
		}

		// ---------- CONTEXT ----------
		ctx := appcontext.WithAppContext(c.Request.Context(), &appcontext.AppContext{
			UserID:   userID,
			TenantID: tenantID,
			Role:     role,
		})
		c.Request = c.Request.WithContext(ctx)

		// ⭐ supaya /me jalan
		c.Set("user", user)

		c.Next()
	}
}
