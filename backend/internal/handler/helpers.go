package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getFamiliaID(c *gin.Context) (string, bool) {
	v, ok := c.Get("familiaID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return "", false
	}
	fid, ok := v.(string)
	if !ok || fid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return "", false
	}
	return fid, true
}
