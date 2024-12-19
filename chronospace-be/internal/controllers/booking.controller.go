package controllers

import (
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"chronospace-be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService services.BookingService
}

func NewBookingController(bookingService services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
	}
}

// @Summary Create booking
// @Description Create a new booking
// @Tags Booking
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param booking body models.CreateBookingParams true "Booking details"
// @Success 201 {object} models.Booking
// @Failure 400 {object} models.ErrorResponse
// @Router /bookings [post]
func (c *BookingController) CreateBooking(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var params models.CreateBookingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.UserID = userID

	booking, err := c.bookingService.CreateBooking(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, booking)
}

// @Summary Get booking
// @Description Get booking by ID
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} models.Booking
// @Failure 400,404 {object} models.ErrorResponse
// @Router /bookings/{id} [get]
func (c *BookingController) GetBooking(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	booking, err := c.bookingService.GetBooking(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

// @Summary List all bookings
// @Description Get all bookings
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {array} models.Booking
// @Failure 400 {object} models.ErrorResponse
// @Router /bookings [get]
func (c *BookingController) ListBookings(ctx *gin.Context) {
	bookings, err := c.bookingService.ListBookings(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}

// @Summary List user bookings
// @Description Get all bookings for the authenticated user
// @Tags Booking
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Booking
// @Failure 400,401 {object} models.ErrorResponse
// @Router /bookings/user [get]
func (c *BookingController) ListUserBookings(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookings, err := c.bookingService.ListBookingsByUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}

// @Summary Update booking
// @Description Update an existing booking
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Param booking body models.UpdateBookingParams true "Booking details"
// @Success 200 {object} models.Booking
// @Failure 400,404 {object} models.ErrorResponse
// @Router /bookings/{id} [put]
func (c *BookingController) UpdateBooking(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var params models.UpdateBookingParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.ID = id

	booking, err := c.bookingService.UpdateBooking(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

// @Summary Delete booking
// @Description Delete a booking
// @Tags Booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 204 "No Content"
// @Failure 400,404 {object} models.ErrorResponse
// @Router /bookings/{id} [delete]
func (c *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := utils.ParseUUID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := c.bookingService.DeleteBooking(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
