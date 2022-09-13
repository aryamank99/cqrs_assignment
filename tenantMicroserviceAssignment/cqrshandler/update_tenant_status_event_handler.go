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

// UpdateTenantStatusHandler is a command handler, which handles CreateTenant command and emits TenantCreated.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be returned.
type UpdateTenantStatusHandler struct {
	EventBus *cqrs.EventBus
}

func (b UpdateTenantStatusHandler) HandlerName() string {
	return "UpdateTenantStatusHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b UpdateTenantStatusHandler) NewCommand() interface{} {
	return &models.UpdateTenantStatusRequest{}
}

func (b UpdateTenantStatusHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*models.UpdateTenantStatusRequest)
	updatedTenant := handler.TenantServiceImpl{}.UpdateTenantStatus(*cmd, database.TenantDbServiceImpl{})
	logrus.Infof("updated tenant status, %s", updatedTenant)
	if updatedTenant == nil {
		return errors.New("unable to update tenant status")
	}
	return nil
}
