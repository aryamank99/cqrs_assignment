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

// CreateTenantHandler is a command handler, which handles CreateTenant command and emits TenantCreated.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be returned.
type CreateTenantHandler struct {
	EventBus *cqrs.EventBus
}

func (b CreateTenantHandler) HandlerName() string {
	return "CreateTenantHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b CreateTenantHandler) NewCommand() interface{} {
	return &models.CreateTenantRequest{}
}

func (b CreateTenantHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*models.CreateTenantRequest)
	newTenant := handler.TenantServiceImpl{}.AddNewTenant(*cmd, database.TenantDbServiceImpl{}, database.LoginCredentialsDbServiceImpl{})
	logrus.Infof("Added a new tenant, %s", newTenant)
	if newTenant == nil {
		return errors.New("unable to save tenant")
	}
	return nil
}
