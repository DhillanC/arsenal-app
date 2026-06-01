package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams extrae limit y offset de query params con defaults seguros.
// Default: limit=20, offset=0. Max limit=100 para evitar cargas masivas.
func PaginationParams(c *gin.Context) (limit, offset int) {
	limit = 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			if parsed > 100 {
				parsed = 100
			}
			limit = parsed
		}
	}

	offset = 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	return
}
