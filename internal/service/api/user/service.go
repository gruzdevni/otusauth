package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	"otusgruz/internal/models"
	query "otusgruz/internal/repo"
	"otusgruz/pkg/http"
)

var ErrNoPermission = errors.New("No permission to perform action")

type repo interface {
	GetUser(ctx context.Context, guid uuid.UUID) (query.User, error)
	DeleteUser(ctx context.Context, guid uuid.UUID) error
	InsertUser(ctx context.Context, arg query.InsertUserParams) error
	UpdateUser(ctx context.Context, arg query.UpdateUserParams) error
}

type service struct {
	repo repo
}

type Service interface {
	GetUser(ctx context.Context, guid uuid.UUID) (*models.UserData, error)
	DeleteUser(ctx context.Context, guid uuid.UUID) (*models.DefaultStatusResponse, error)
	UpdateUser(ctx context.Context, guid uuid.UUID, info *models.UserCreateParams) (*models.DefaultStatusResponse, error)
	CreateUser(ctx context.Context, info *models.UserCreateParams) (*models.DefaultStatusResponse, error)
}

func NewService(repo repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUser(ctx context.Context, guid uuid.UUID) (*models.UserData, error) {
	ctxUserGUID := http.UserGUIDFromContext(ctx)

	if ctxUserGUID == uuid.Nil || ctxUserGUID != guid {
		return nil, ErrNoPermission
	}

	res, err := s.repo.GetUser(ctx, guid)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &models.UserData{
		GUID:       strfmt.UUID(res.Guid.String()),
		IsDeleted:  res.IsDeleted,
		Name:       res.Name,
		Occupation: res.Occupation,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, guid uuid.UUID) (*models.DefaultStatusResponse, error) {
	ctxUserGUID := http.UserGUIDFromContext(ctx)

	if ctxUserGUID == uuid.Nil || ctxUserGUID != guid {
		return nil, ErrNoPermission
	}

	err := s.repo.DeleteUser(ctx, guid)
	if err != nil {
		return nil, fmt.Errorf("deleting user: %w", err)
	}

	return &models.DefaultStatusResponse{
		Code:    "01",
		Message: "Successfully deleted",
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, guid uuid.UUID, info *models.UserCreateParams) (*models.DefaultStatusResponse, error) {
	ctxUserGUID := http.UserGUIDFromContext(ctx)

	if ctxUserGUID == uuid.Nil || ctxUserGUID != guid {
		return nil, ErrNoPermission
	}

	err := s.repo.UpdateUser(ctx, query.UpdateUserParams{
		Guid:       guid,
		Occupation: info.Occupation,
		Name:       info.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return &models.DefaultStatusResponse{
		Code:    "01",
		Message: "Successfully updated",
	}, nil
}

func (s *service) CreateUser(ctx context.Context, info *models.UserCreateParams) (*models.DefaultStatusResponse, error) {
	ctxUserGUID := http.UserGUIDFromContext(ctx)

	if ctxUserGUID == uuid.Nil {
		return nil, ErrNoPermission
	}

	err := s.repo.InsertUser(ctx, query.InsertUserParams{
		Guid:       uuid.New(),
		Occupation: info.Occupation,
		Name:       info.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return &models.DefaultStatusResponse{
		Code:    "01",
		Message: "Successfully created",
	}, nil
}
