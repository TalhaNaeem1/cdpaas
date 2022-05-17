//go:generate mockgen -destination=mocks/mock_authservice.go -package=mock_authservice . AuthServiceQuerier
package authService

import (
	"github.com/gin-gonic/gin"
	"pipelineService/models/v1"
)

type AuthServiceQuerier interface {
	GetUserByID(ownerID int) (models.UserDetails, error)
	ValidateSession(ctx *gin.Context)
}

var _ AuthServiceQuerier = (*RequestMaker)(nil)

