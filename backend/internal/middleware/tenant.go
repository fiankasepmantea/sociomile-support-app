package middleware

import (
	"net/http"
	"strings"

	"backend/internal/appcontext"
	"backend/internal/apperrors"
	"backend/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract tenant identifier (subdomain atau header)
		host := c.Request.Host
		parts := strings.Split(host, ".")

		var tenantIdentifier string
		if len(parts) > 1 && parts[0] != "www" && parts[0] != "api" && parts[0] != "localhost" {
			tenantIdentifier = parts[0] // subdomain: acme.support.com
		} else {
			tenantIdentifier = c.GetHeader("X-Tenant-ID")
		}

		if tenantIdentifier == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant identifier required (subdomain or X-Tenant-ID header)"})
			c.Abort()
			return
		}

		// Cari tenant aktif di database
		var tenant model.Tenant
		if err := db.Where("name = ? AND deleted_at IS NULL", tenantIdentifier).
			Or("id::text = ? AND deleted_at IS NULL", tenantIdentifier).
			First(&tenant).Error; err != nil {

			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid tenant"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		// Simpan tenant ke Gin context dan standard context
		c.Set("tenant", &tenant)
		ctx := appcontext.WithAppContext(c.Request.Context(), &appcontext.AppContext{
			TenantID: tenant.ID,
		})
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetTenantFromContext helper untuk ambil tenant dari Gin context
func GetTenant(c *gin.Context) (*model.Tenant, error) {
	tenant, exists := c.Get("tenant")
	if !exists {
		return nil, apperrors.ErrInternal("tenant not found in context")
	}
	if t, ok := tenant.(*model.Tenant); ok {
		return t, nil
	}
	return nil, apperrors.ErrInternal("invalid tenant type in context")
}