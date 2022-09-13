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

// UpdateTenantLicenseHandler is a command handler, which handles CreateTenant command and emits TenantCreated.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be returned.
type UpdateTenantLicenseHandler struct {
	EventBus *cqrs.EventBus
}

func (b UpdateTenantLicenseHandler) HandlerName() string {
	return "UpdateTenantLicenseHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b UpdateTenantLicenseHandler) NewCommand() interface{} {
	return &models.UpdateTenantLicenseRequest{}
}

func (b UpdateTenantLicenseHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*models.UpdateTenantLicenseRequest)
	updatedTenant := handler.TenantServiceImpl{}.UpdateTenantLicense(*cmd, database.TenantDbServiceImpl{})
	logrus.Infof("updated tenant lisence, %s", updatedTenant)
	if updatedTenant == nil {
		return errors.New("unable to update tenant license")
	}
	return nil
}
