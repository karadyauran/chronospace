package controllers

import (
	"chronospace-be/internal/models"
	"chronospace-be/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary Register new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserParams true "User registration details"
// @Success 201 {object} models.UserCreatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/api/users/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var params models.CreateUserParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := c.userService.RegisterUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary User login
// @Description Authenticate user and return access token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /v1/api/users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	response, err := c.userService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary User logout
// @Description Logout user and invalidate their token
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /v1/api/users/logout [post]
func (c *UserController) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "unauthorized"})
		return
	}

	if err := c.userService.LogoutUser(ctx, userID.(pgtype.UUID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{Message: "logged out successfully"})
}

// @Summary Get user profile
// @Description Get user profile by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /v1/api/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user ID"})
		return
	}

	user, err := c.userService.GetUser(ctx, uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Update user profile
// @Description Update user profile information
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUserParams true "User update information"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/api/users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user ID"})
		return
	}

	var params models.UpdateUserParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := c.userService.UpdateUser(ctx, uuid, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete user account
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/api/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var uuid pgtype.UUID
	if err := uuid.Scan(userID); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid user ID"})
		return
	}

	if err := c.userService.DeleteUser(ctx, uuid); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{Message: "user deleted successfully"})
}

// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit number of users"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/api/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	var params models.ListUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	users, err := c.userService.ListUsers(ctx, params.ToDBParams())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
