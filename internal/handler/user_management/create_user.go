package user_management

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/sonikq/gravitum_test_task/pkg/reader"
	"net/http"
)

func (h *Handler) CreateUser(ctx *gin.Context) {
	const source = "handler.CreateUser"

	if ctx.GetHeader(contentTypeHeaderKey) != contentTypeJSON {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{models.ErrMsgKey: "Invalid type of content"})
		h.logger.Error().
			Str("error", "invalid content type").
			Str("source", source).
			Send()
		return
	}

	bodyBytes, err := reader.GetBody(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{models.ErrMsgKey: "Error in reading request body"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to read request body")
		return
	}

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	var request models.UserInfo
	if err = json.Unmarshal(bodyBytes, &request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{models.ErrMsgKey: "Error in parsing request body"})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg("failed to unmarshal request body")
		return
	}

	id, err := h.service.CreateUser(c, request)
	if err != nil {
		var (
			statusCode int
			userMsg    string
			logMsg     = "failed to create user"
		)

		switch {
		case errors.Is(err, models.ErrUsernameIsAlreadyTaken):
			statusCode = http.StatusConflict
			userMsg = models.ErrUsernameIsAlreadyTaken.Error()

		case errors.Is(err, models.ErrInvalidEmail):
			statusCode = http.StatusBadRequest
			userMsg, logMsg = models.ErrInvalidEmail.Error(), "failed to validate email"

		case errors.Is(err, models.ErrInvalidAge):
			statusCode = http.StatusBadRequest
			userMsg, logMsg = models.ErrInvalidEmail.Error(), "failed to validate age"

		case errors.Is(err, models.ErrInvalidGender):
			statusCode = http.StatusBadRequest
			userMsg, logMsg = models.ErrInvalidEmail.Error(), "failed to validate gender"

		default:
			statusCode = http.StatusInternalServerError
			userMsg = "internal server error, something went wrong"
		}

		ctx.AbortWithStatusJSON(statusCode, gin.H{models.ErrMsgKey: userMsg})
		h.logger.Error().
			Err(err).
			Str("source", source).
			Msg(logMsg)
		return
	}
	ctx.Data(http.StatusCreated, contentTypeTextPlain, []byte(id))
}
