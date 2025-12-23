package rest

import (
	"fmt"
	"net/http"
	"time"

	"learngolang/src/dto"
	exception "learngolang/src/errors"
	"learngolang/src/preference"

	"github.com/gin-gonic/gin"
)

func (e *rest) httpRespSuccess(c *gin.Context, statusCode int, resp any, p *dto.Pagination) {
	meta := dto.Meta{
		Path:       c.Request.URL.Path,
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)),
		Error:      nil,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	HttpResp := &dto.HttpSuccessResp{
		Meta:       meta,
		Data:       any(resp),
		Pagination: p,
	}

	c.JSON(statusCode, HttpResp)
}

func (e *rest) httpRespError(c *gin.Context, appErr error) {
	lang := preference.LANG_ID

	appLangHeader := http.CanonicalHeaderKey(preference.APP_LANG)
	if c.Request.Header[appLangHeader] != nil && c.Request.Header[appLangHeader][0] == preference.LANG_EN {
		lang = preference.LANG_EN
	}

	statusCode, displayError := exception.Compile(exception.COMMON, appErr, lang, true)
	statusStr := http.StatusText(statusCode)

	jsonErrResp := &dto.HTTPErrorResp{
		Meta: dto.Meta{
			Path:       c.Request.URL.Path,
			StatusCode: statusCode,
			Status:     statusStr,
			Message:    fmt.Sprintf("%s %s [%d] %s", c.Request.Method, c.Request.RequestURI, statusCode, http.StatusText(statusCode)),
			Error:      &displayError,
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	}

	c.JSON(statusCode, jsonErrResp)
}
