package rest

import (
	"net/http"

	"learngolang/src/dto"
	exception "learngolang/src/errors"
	"learngolang/src/preference"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (r *rest) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Invalid request body")
		r.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPUnmarshal, "invalid_request_body"))
		return
	}

	user, err := r.svc.User.CreateUser(ctx, req)
	if err != nil {
		r.httpRespError(c, err)
		return
	}

	r.httpRespSuccess(c, http.StatusCreated, user, nil)
}

func (r *rest) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_user_id")
		r.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPBadRequest, "invalid_user_id"))
		return
	}

	user, err := r.svc.User.GetUser(ctx, id.String())
	if err != nil {
		r.httpRespError(c, err)
		return
	}

	r.httpRespSuccess(c, http.StatusOK, user, nil)
}

func (e *rest) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var (
		filter       dto.UserFilter
		cacheControl dto.CacheControl
	)

	if err := c.ShouldBindQuery(&filter); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_query_parameters")
		e.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPBadRequest, "invalid_query_parameters"))
		return
	}

	if c.Request.Header[http.CanonicalHeaderKey(preference.CacheControl)] != nil && c.Request.Header[http.CanonicalHeaderKey(preference.CacheControl)][0] == preference.CacheMustRevalidate {
		cacheControl.MustRevalidate = true
	}

	users, pagination, err := e.svc.User.ListUsers(ctx, cacheControl, filter)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, users, &pagination)
}

func (e *rest) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_user_id")
		e.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPBadRequest, "Invalid user ID"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_request_body")
		e.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPUnmarshal, "Invalid request body"))
		return
	}

	user, err := e.svc.User.UpdateUser(ctx, id.String(), req)
	if err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, user, nil)
}

func (e *rest) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("invalid_user_id")
		e.httpRespError(c, exception.WrapWithCode(err, exception.CodeHTTPBadRequest, "Invalid user ID"))
		return
	}

	if err := e.svc.User.DeleteUser(ctx, id.String()); err != nil {
		e.httpRespError(c, err)
		return
	}

	e.httpRespSuccess(c, http.StatusOK, nil, nil)
}
