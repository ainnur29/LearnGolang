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

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		domain.CreateUserRequest	true	"User data"
//	@Success		201		{object}	domain.Response{data=domain.User}
//	@Failure		400		{object}	domain.Response
//	@Failure		500		{object}	domain.Response
//	@Router			/users [post]
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

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get a user by their ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	domain.Response{data=domain.User}
//	@Failure		404	{object}	domain.Response
//	@Failure		500	{object}	domain.Response
//	@Router			/users/{id} [get]
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

// ListUsers godoc
//
//	@Summary		List users
//	@Description	Get a paginated list of users with optional filters
//	@Tags			users
//	@Produce		json
//	@Param			name		query		string	false	"Filter by name"
//	@Param			email		query		string	false	"Filter by email"
//	@Param			min_age		query		int		false	"Minimum age"
//	@Param			max_age		query		int		false	"Maximum age"
//	@Param			page		query		int		false	"Page number"	default(1)
//	@Param			page_size	query		int		false	"Page size"		default(10)
//	@Param			sort_by		query		string	false	"Sort by field"
//	@Param			sort_dir	query		string	false	"Sort direction (asc/desc)"	default(asc)
//	@Success		200			{object}	domain.PaginatedResponse{data=[]domain.User}
//	@Failure		400			{object}	domain.Response
//	@Failure		500			{object}	domain.Response
//	@Router			/users [get]
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

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	Update an existing user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"User ID"
//	@Param			user	body		domain.UpdateUserRequest	true	"User data"
//	@Success		200		{object}	domain.Response{data=domain.User}
//	@Failure		400		{object}	domain.Response
//	@Failure		404		{object}	domain.Response
//	@Failure		500		{object}	domain.Response
//	@Router			/users/{id} [put]
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

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	domain.Response
//	@Failure		404	{object}	domain.Response
//	@Failure		500	{object}	domain.Response
//	@Router			/users/{id} [delete]
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
