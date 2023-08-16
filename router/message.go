package router

import (
	"c6m/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleGetHistory(c *gin.Context) {
	uid := c.MustGet("uid").(string)
	id := c.Query("id")

	msgStorage, err := db.GetMsgStorage(uid, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, msgStorage)
}
