package user_management

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sonikq/gravitum_test_task/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteUser(ctx *gin.Context) {
	const source = "handler.DeleteUser"
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

	c, cancel := context.WithTimeout(ctx, h.config.CtxTimeOut)
	defer cancel()

	err = h.service.DeleteUser(c, int64(userID))
	if err != nil {
		var (
			statusCode int
			userMsg    string
			logMsg     = "failed to delete user"
		)

		switch {
		case errors.Is(err, models.ErrUserDoesNotExist):
			statusCode = http.StatusNoContent
			userMsg = models.ErrUserDoesNotExist.Error()

		case errors.Is(err, models.ErrUserIsGone):
			statusCode = http.StatusConflict
			userMsg, logMsg = models.ErrDeleteDeletedUser.Error(), "failed to delete deleted user"

		case errors.Is(err, models.ErrDeleteDeletedUser):
			statusCode = http.StatusConflict
			userMsg, logMsg = models.ErrDeleteDeletedUser.Error(), "failed to delete deleted user"
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
