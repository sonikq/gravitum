package user_management

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"github.com/sonikq/gravitum_test_task/pkg/reader"
	"net/http"
	"strconv"
)

func (h *Handler) UpdateUser(ctx *gin.Context) {
	const source = "handler.UpdateUser"

	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{models.ErrMsgKey: "Invalid type of user_id"})
		h.logger.Error().
			Str("error", "invalid type of user_id: "+userIDStr).
			Str("source", source).
			Send()
		return
	}

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
	request.ID = int64(userID)

	err = h.service.UpdateUser(c, request)
	if err != nil {
		var (
			statusCode int
			userMsg    string
			logMsg     = "failed to update user"
		)

		switch {
		case errors.Is(err, models.ErrUserDoesNotExist):
			statusCode = http.StatusNoContent
			userMsg = models.ErrUserDoesNotExist.Error()

		case errors.Is(err, models.ErrUserIsGone):
			statusCode = http.StatusGone
			userMsg, logMsg = models.ErrUserIsGone.Error(), "failed to update deleted user"

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
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
