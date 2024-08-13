package services

import (
	errApi "chat-system/cmd/api/errors"
	"chat-system/cmd/api/models"
	"chat-system/cmd/api/repositories"
	"chat-system/cmd/api/security"
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type UserService interface {
	RegisterUser(ctx context.Context, user models.User) errApi.ErrAPI
	VerifyUserPassword(ctx context.Context, username string, password string) (idUser string, ok bool, errs errApi.ErrAPI)
	ListUsers(ctx context.Context) ([]models.User, errApi.ErrAPI)
}

type userServiceImpl struct {
	dbRepository repositories.DBRepository
}

func NewUserService(dbRepository repositories.DBRepository) UserService {
	return &userServiceImpl{dbRepository: dbRepository}
}

func (u *userServiceImpl) RegisterUser(ctx context.Context, user models.User) errApi.ErrAPI {
	exists, err := u.dbRepository.UserExists(ctx, user.Username)
	if err != nil {
		return errApi.NewErrAPIInternalServer(err)
	}
	if exists {
		return errApi.NewErrAPIBadRequest(errors.New("username already exists"))
	}

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		return errApi.NewErrAPIInternalServer(err)
	}

	user.Id = gocql.TimeUUID()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = u.dbRepository.SaveUser(ctx, user)
	if err != nil {
		return errApi.NewErrAPIInternalServer(err)
	}

	return nil
}

func (u *userServiceImpl) VerifyUserPassword(ctx context.Context, username string, password string) (string, bool, errApi.ErrAPI) {
	idUser, storedHash, err := u.dbRepository.GetUserPassword(ctx, username)
	if err != nil {
		return "", false, err
	}

	return idUser, security.CheckPasswordHash(password, storedHash), nil
}

func (u *userServiceImpl) ListUsers(ctx context.Context) ([]models.User, errApi.ErrAPI) {
	users, err := u.dbRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, errApi.NewErrAPIInternalServer(err)
	}

	if len(users) == 0 {
		return nil, errApi.NewErrAPINotFound(errors.New("no users found"))
	}

	return users, nil
}
