package services

import (
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/cts3njitedu/healthful-heart/repositories/mysqlrepo"
	"time"
)

type LogoutService struct {
	workflowService IWorkflowService
	tokenRepository mysqlrepo.ITokenRepository

}

func NewLogoutService(workflowService IWorkflowService,
	tokenRepository mysqlrepo.ITokenRepository) *LogoutService {
	return &LogoutService{workflowService, tokenRepository}
}

func (logout *LogoutService) LogoutService(credentials models.Credentials) (error) {

	err := logout.tokenRepository.SaveRefreshToken("", time.Time{}, credentials.UserId)

	if err != nil {
		return err;
	}
	return nil
}