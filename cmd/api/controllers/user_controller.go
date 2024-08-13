package controllers

import (
	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"chat-system/cmd/api/security"
	"chat-system/cmd/api/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	RegisterUser(c *gin.Context)
	VerifyUserPassword(c *gin.Context)
	ListUsers(c *gin.Context)
}

type userControllerImpl struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userControllerImpl{userService: userService}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user by providing the necessary details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Failure 400 {object} errors.ErrAPI
// @Failure 401 {object} errors.ErrAPI
// @Failure 500 {object} errors.ErrAPI
// @Router /users/register [post]
func (u *userControllerImpl) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		apiErr := errApi.NewErrAPIBadRequest(errors.New("invalid request payload"))
		c.JSON(apiErr.GetHTTPStatusCode(), apiErr)
		return
	}

	err := u.userService.RegisterUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(err.GetHTTPStatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Message: "User registered successfully"})
}

// VerifyUserPassword godoc
// @Summary Verify user password
// @Description Verify the password of a user for login
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "Login Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.ErrAPI
// @Failure 401 {object} errors.ErrAPI
// @Failure 500 {object} errors.ErrAPI
// @Router /users/login [post]
func (u *userControllerImpl) VerifyUserPassword(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		apiErr := errApi.NewErrAPIBadRequest(errors.New("invalid request payload"))
		c.JSON(apiErr.GetHTTPStatusCode(), apiErr)
		return
	}

	idUser, loginSuccess, err := u.userService.VerifyUserPassword(c.Request.Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(err.GetHTTPStatusCode(), err)
		return
	}

	if loginSuccess {
		token, tokenErr := security.GenerateToken(idUser)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, errApi.NewErrAPIInternalServer(tokenErr))
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Message: "Login successful",
			Data:    token,
		})
	} else {
		apiErr := errApi.NewErrAPIUnauthorized(errors.New("login failed"))
		c.JSON(apiErr.GetHTTPStatusCode(), apiErr)
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieves a list of all registered users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.User} "List of users"
// @Failure 404 {object} errors.ErrAPI "No users found"
// @Failure 500 {object} errors.ErrAPI "Internal Server Error"
// @Router /users [get]
func (u *userControllerImpl) ListUsers(c *gin.Context) {
	users, err := u.userService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(err.GetHTTPStatusCode(), err)
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: users})
}
