package restapi

import (
	"errors"

	"otusgruz/internal/models"
	"otusgruz/internal/restapi/operations/other"
	"otusgruz/internal/restapi/operations/user_c_r_u_d"
	"otusgruz/internal/service/api/auth"
	"otusgruz/internal/service/api/user"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type Handler struct {
	userSrv user.Service
	authSrv auth.Service
}

func NewHandler(userSrv user.Service, authSrv auth.Service) *Handler {
	return &Handler{
		userSrv: userSrv,
		authSrv: authSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) Auth(params other.GetAuthParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	if params.XUser == nil {
		return other.NewGetAuthUnauthorized()
	}

	userGUID, err := uuid.Parse(*params.XUser)
	if err != nil {
		return other.NewGetAuthUnauthorized()
	}

	authorizedGUID, err := h.authSrv.Auth(ctx, userGUID)
	if err != nil {
		return other.NewGetAuthInternalServerError()
	}

	if authorizedGUID == uuid.Nil {
		return other.NewGetAuthUnauthorized()
	}

	return other.NewGetAuthOK().WithXUser(strfmt.UUID(authorizedGUID.String()))
}

func (h *Handler) Login(params other.PostLoginParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	loginParams := params.Request

	authorizedGUID, err := h.authSrv.Login(ctx, loginParams.Email.String(), loginParams.Password)
	if err != nil {
		if errors.Is(err, auth.ErrNotCorrectData) {
			return other.NewPostLoginUnauthorized().WithPayload(&models.DefaultStatusResponse{Message: auth.ErrNotCorrectData.Error()})
		}

		return other.NewPostLoginInternalServerError()
	}

	return other.NewPostLoginOK().WithXUser(strfmt.UUID(authorizedGUID.String()))
}

func (h *Handler) Signup(params other.PostSignupParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	regParams := params.Request

	err := h.authSrv.Singup(ctx, regParams.Email.String(), regParams.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyUsed) {
			return other.NewPostSignupOK().WithPayload(&models.DefaultStatusResponse{Message: auth.ErrEmailAlreadyUsed.Error()})
		}

		return other.NewPostSignupInternalServerError()
	}

	return other.NewPostSignupOK()
}

func (h *Handler) GetUser(params user_c_r_u_d.GetUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewGetUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.GetUser(ctx, userGUID)
	if err != nil {
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewGetUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

		errText = err.Error()
		return user_c_r_u_d.NewGetUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewGetUserGUIDOK().WithPayload(res)
}

func (h *Handler) CreateUser(params user_c_r_u_d.PostUserParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	res, err := h.userSrv.CreateUser(ctx, params.Request)
	if err != nil {
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewPostUserForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

		errText = err.Error()
		return user_c_r_u_d.NewPostUserInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewPostUserOK().WithPayload(res)
}

func (h *Handler) UpdateUser(params user_c_r_u_d.PatchUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewPatchUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

		errText = err.Error()
		return user_c_r_u_d.NewPatchUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.UpdateUser(ctx, userGUID, params.Request)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewPatchUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewPatchUserGUIDOK().WithPayload(res)
}

func (h *Handler) DeleteUser(params user_c_r_u_d.DeleteUserGUIDParams) middleware.Responder {
	var errText string
	ctx := params.HTTPRequest.Context()

	userGUID, err := uuid.Parse(params.GUID.String())
	if err != nil {
		if errors.Is(err, user.ErrNoPermission) {
			return user_c_r_u_d.NewDeleteUserGUIDForbidden().WithPayload(&models.DefaultStatusResponse{Message: user.ErrNoPermission.Error()})
		}

		errText = err.Error()
		return user_c_r_u_d.NewDeleteUserGUIDBadRequest().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	res, err := h.userSrv.DeleteUser(ctx, userGUID)
	if err != nil {
		errText = err.Error()
		return user_c_r_u_d.NewDeleteUserGUIDInternalServerError().WithPayload(&models.Error{Code: 0o3, Message: &errText})
	}

	return user_c_r_u_d.NewDeleteUserGUIDOK().WithPayload(res)
}
