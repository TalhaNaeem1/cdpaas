package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

func ParseDBError(err error, ph string) (int, string) {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case pgerrcode.ForeignKeyViolation:
			return http.StatusBadRequest, fmt.Sprintf("Failed to create %s", ph)
		case pgerrcode.UniqueViolation:
			return http.StatusBadRequest, fmt.Sprintf("%s with this Info already exists", ph)
		case pgerrcode.CaseNotFound:
			return http.StatusNotFound, fmt.Sprintf("%s doesn't exist", ph)
		case pgerrcode.UndefinedTable:
			return http.StatusInternalServerError, fmt.Sprintf("%s table doesn't exist", ph)
		default:
			return http.StatusBadRequest, fmt.Sprintf("Something went wrong with %s", ph)
		}
	}

	return http.StatusBadRequest, "Something went wrong"
}

func GetUserAndWorkspaceIDFromContext(ctx *gin.Context) (int, int, string) {
	userID := ctx.Request.Header.Get("userID")
	workspaceID := ctx.Request.Header.Get("workspaceID")
	airbyteWorkspaceId := ctx.Request.Header.Get("airbyteWorkspaceID")

	userId, _ := strconv.Atoi(userID)
	workspaceId, _ := strconv.Atoi(workspaceID)

	return userId, workspaceId, airbyteWorkspaceId
}
