package controllers

import (
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ScheduleController struct {
	scheduleService *services.ScheduleService
}

func NewScheduleController(scheduleService services.ScheduleService) *ScheduleController {
	return &ScheduleController{
		scheduleService: &scheduleService,
	}
}

// CreateSchedule godoc
// @Summary Create a new schedule
// @Description Create a new schedule with the given time range
// @Tags schedules
// @Accept json
// @Produce json
// @Param schedule body models.CreateScheduleRequest true "Schedule details"
// @Success 201 {object} models.ScheduleResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/api/schedules [post]
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var req models.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	schedule, err := c.scheduleService.CreateSchedule(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, schedule)
}

// GetSchedule godoc
// @Summary Get a schedule by ID
// @Description Get schedule details by its ID
// @Tags schedules
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} models.ScheduleResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/api/schedules/{id} [get]
func (c *ScheduleController) GetSchedule(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid schedule ID"})
		return
	}

	schedule, err := c.scheduleService.GetSchedule(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schedule)
}

// ListSchedules godoc
// @Summary List all schedules
// @Description Get a list of all schedules
// @Tags schedules
// @Produce json
// @Success 200 {array} models.ScheduleResponse
// @Router /v1/api/schedules [get]
func (c *ScheduleController) ListSchedules(ctx *gin.Context) {
	schedules, err := c.scheduleService.ListSchedules(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schedules)
}

// UpdateSchedule godoc
// @Summary Update a schedule
// @Description Update a schedule's details by its ID
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param schedule body models.UpdateScheduleRequest true "Schedule details"
// @Success 200 {object} models.ScheduleResponse
// @Failure 400,404 {object} models.ErrorResponse
// @Router /v1/api/schedules/{id} [put]
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid schedule ID"})
		return
	}

	var req models.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	schedule, err := c.scheduleService.UpdateSchedule(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schedule)
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Description Delete a schedule by its ID
// @Tags schedules
// @Param id path string true "Schedule ID"
// @Success 204 "No Content"
// @Failure 400,404 {object} models.ErrorResponse
// @Router /schedules/{id} [delete]
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	var id pgtype.UUID
	if err := id.Scan(ctx.Param("id")); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid schedule ID"})
		return
	}

	if err := c.scheduleService.DeleteSchedule(ctx, id); err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}