package cqrshandler

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/sirupsen/logrus"
	"tenantMicroserviceAssignment/database"
	"tenantMicroserviceAssignment/handler"
	"tenantMicroserviceAssignment/models"
)

// UpdateTenantCredentialsHandler is a command handler, which handles UpdateTenantCredentials command and emits TenantCreated.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be retuerned.
type UpdateTenantCredentialsHandler struct {
	EventBus *cqrs.EventBus
}

func (b UpdateTenantCredentialsHandler) HandlerName() string {
	return "UpdateTenantCredentialsHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b UpdateTenantCredentialsHandler) NewCommand() interface{} {
	return &models.UpdateTenantCredentialsRequest{}
}

func (b UpdateTenantCredentialsHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*models.UpdateTenantCredentialsRequest)
	updatedCredentials := handler.TenantServiceImpl{}.UpdateTenantCredentials(*cmd, database.LoginCredentialsDbServiceImpl{})
	logrus.Infof("updated tenant creds")
	if updatedCredentials == nil {
		return errors.New("unable to update tenant credentials")
	}
	return nil
}
