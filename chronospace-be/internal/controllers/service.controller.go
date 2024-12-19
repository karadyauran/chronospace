package controllers

import (
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"chronospace-be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	serviceService services.ServiceService
}

func NewServiceController(serviceService services.ServiceService) *ServiceController {
	return &ServiceController{
		serviceService: serviceService,
	}
}

// @Summary Create service
// @Description Create a new service
// @Tags Service
// @Accept json
// @Produce json
// @Param service body models.CreateServiceRequest true "Service details"
// @Success 201 {object} models.ServiceResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /services [post]
func (c *ServiceController) CreateService(ctx *gin.Context) {
	var req models.CreateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := c.serviceService.CreateService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, service)
}

// @Summary Get service
// @Description Get service by ID
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} models.ServiceResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Router /services/{id} [get]
func (c *ServiceController) GetService(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	service, err := c.serviceService.GetService(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

// @Summary Update service
// @Description Update an existing service
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body models.UpdateServiceRequest true "Service details"
// @Success 200 {object} models.ServiceResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Router /services/{id} [put]
func (c *ServiceController) UpdateService(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	var req models.UpdateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := c.serviceService.UpdateService(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, service)
}

// @Summary Delete service
// @Description Delete a service
// @Tags Service
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 204 "No Content"
// @Failure 400,404 {object} models.ErrorResponse
// @Router /services/{id} [delete]
func (c *ServiceController) DeleteService(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	if err := c.serviceService.DeleteService(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary List services
// @Description Get all services
// @Tags Service
// @Accept json
// @Produce json
// @Success 200 {array} models.ServiceResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /services [get]
func (c *ServiceController) ListServices(ctx *gin.Context) {
	services, err := c.serviceService.ListServices(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, services)
}