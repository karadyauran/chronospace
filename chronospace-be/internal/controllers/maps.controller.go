package controllers

import (
	"chronospace-be/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MapsController struct {
	mapsService *services.MapsService
}

func NewMapsController(mapsService *services.MapsService) *MapsController {
	return &MapsController{
		mapsService: mapsService,
	}
}

// @Summary Search places
// @Description Search for places using Google Maps API
// @Tags Maps
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {object} interface{}
// @Router /v1/api/maps/search [get]
func (c *MapsController) SearchPlaces(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	results, err := c.mapsService.SearchPlaces(query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}
