package api

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"tenantMicroserviceAssignment/database"
	"tenantMicroserviceAssignment/handler"
	"tenantMicroserviceAssignment/models"
)

type TenantApi struct{}

// AddNewTenant RootHandler - Adds a new tenant to the database.
// @Summary This API can be used to create a new tenant.
// @Description Adds a new tenant in the database
// @Tags tenant
// @Accept  json
// @Produce  json
// @Param default body models.CreateTenantRequestBody true "create tenant request"
// @Success 200
// @Failure 500 {string} error response "when there's some error adding the tenant, can be a database failure or a tenant already exist with the id"
// @Router /tenant [post]
func (api TenantApi) AddNewTenant(writer http.ResponseWriter, request *http.Request, facade *cqrs.Facade) {
	// 1: create tenant
	newTenantRequest := models.CreateTenantRequestBody{}
	json.NewDecoder(request.Body).Decode(&newTenantRequest)
	requestId, _ := uuid.NewUUID()
	command := models.CreateTenantRequest{
		RequestId:        requestId.String(),
		Name:             newTenantRequest.Name,
		Username:         newTenantRequest.Username,
		Password:         newTenantRequest.Password,
		TotalUserCount:   newTenantRequest.TotalUserCount,
		IsEnabled:        newTenantRequest.IsEnabled,
		LicenseType:      newTenantRequest.LicenseType,
		LicenseStartDate: newTenantRequest.LicenseStartDate,
		LicenseEndDate:   newTenantRequest.LicenseEndDate,
	}
	if err := facade.CommandBus().Send(context.Background(), &command); err != nil {
		panic(err)
	}
	result := database.GetCommandsDbConn().Create(&command)
	if result.Error != nil {
		logrus.Errorf("Unable to save create tenant request, %s", result.Error)
	}
	writer.WriteHeader(Ok)
}

// GetAllTenants / RootHandler - Returns all tenants in the database.
// @Summary This API can be used to get all tenants.
// @Description Designed to be used with some sort of admin web panel, etc.
// @Tags tenant
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Tenant "when the operation is successful"
// @Failure 500 {string} error response "when failed reading the database"
// @Router /tenant [get]
func (api TenantApi) GetAllTenants(writer http.ResponseWriter, _ *http.Request) {
	tenants := handler.TenantServiceImpl{}.GetAllTenants(database.TenantDbServiceImpl{})
	writer.WriteHeader(Ok)
	json.NewEncoder(writer).Encode(tenants)
}

// GetTenantById / RootHandler - Returns a tenant matching their tenant id
// @Summary This API can be used to get a tenant by id.
// @Description Retrieves a tenant from the database by id.
// @Tags tenant
// @Accept  json
// @Produce  json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} models.Tenant "when the operation is successful"
// @Failure 404 {string} error response "when no tenant is found against the provided tenant id"
// @Router /tenant/{tenant_id} [get]
func (api TenantApi) GetTenantById(writer http.ResponseWriter, request *http.Request) {
	tenantId := chi.URLParam(request, "tenant_id")
	logrus.Println("getting tenant with id :" + tenantId)
	tenant := handler.TenantServiceImpl{}.GetTenantById(tenantId, database.TenantDbServiceImpl{})
	if tenant != nil {
		writer.WriteHeader(Ok)
		json.NewEncoder(writer).Encode(*tenant)
	} else {
		writer.WriteHeader(NotFound)
	}
}

// UpdateTenantCredentials RootHandler - Updates an existing tenant credentials
// @Summary This API can be used to update a tenant credentials.
// @Description Updates a tenant credentials in the database.
// @Tags tenant
// @Accept  json
// @Produce  json
// @Param default body models.UpdateTenantCredentialsRequestBody true "update tenant credentials request"
// @Success 200
// @Failure 406 {string} error response "when there's some error updating the tenant credentials, can be a database failure or no tenant found with the id"
// @Router /tenant/credentials [put]
func (api TenantApi) UpdateTenantCredentials(writer http.ResponseWriter, request *http.Request, facade *cqrs.Facade) {
	_, claims, _ := jwtauth.FromContext(request.Context())
	tenantId := claims["tenant_id"]
	updateRequest := models.UpdateTenantCredentialsRequestBody{}
	json.NewDecoder(request.Body).Decode(&updateRequest)
	requestId, _ := uuid.NewUUID()
	command := models.UpdateTenantCredentialsRequest{
		RequestId:   requestId.String(),
		TenantId:    tenantId.(string),
		NewPassword: updateRequest.NewPassword,
	}
	if err := facade.CommandBus().Send(context.Background(), &command); err != nil {
		panic(err)
	}
	result := database.GetCommandsDbConn().Create(&command)
	if result.Error != nil {
		logrus.Errorf("Unable to save update tenant request, %s", result.Error)
	}
	writer.WriteHeader(Ok)
}

// UpdateTenantStatus RootHandler - Updates an existing tenant status
// @Summary This API can be used to update a tenant status.
// @Description Updates a tenant status in the database.
// @Tags tenant
// @Accept  json
// @Produce  json
// @Param default body models.UpdateTenantStatusRequestBody true "update tenant status request"
// @Success 200
// @Failure 406 {string} error response "when there's some error updating the tenant status, can be a database failure or no tenant found with the id"
// @Router /tenant/status [put]
func (api TenantApi) UpdateTenantStatus(writer http.ResponseWriter, request *http.Request, facade *cqrs.Facade) {
	_, claims, _ := jwtauth.FromContext(request.Context())
	tenantId := claims["tenant_id"]
	updateRequest := models.UpdateTenantStatusRequestBody{}
	json.NewDecoder(request.Body).Decode(&updateRequest)
	requestId, _ := uuid.NewUUID()
	command := models.UpdateTenantStatusRequest{
		RequestId: requestId.String(),
		TenantId:  tenantId.(string),
		IsEnabled: updateRequest.IsEnabled,
	}
	if err := facade.CommandBus().Send(context.Background(), &command); err != nil {
		panic(err)
	}
	result := database.GetCommandsDbConn().Create(&command)
	if result.Error != nil {
		logrus.Errorf("Unable to save update tenant status request, %s", result.Error)
	}
	writer.WriteHeader(Ok)
}

// UpdateTenantLicense RootHandler - Updates an existing tenant license
// @Summary This API can be used to update a tenant license.
// @Description Updates a tenant license in the database.
// @Tags tenant
// @Accept  json
// @Produce  json
// @Param default body models.UpdateTenantLicenseRequestBody true "update tenant license request"
// @Success 200
// @Failure 406 {string} error response "when there's some error updating the tenant, can be a database failure or no tenant found with the id"
// @Router /tenant/license [put]
func (api TenantApi) UpdateTenantLicense(writer http.ResponseWriter, request *http.Request, facade *cqrs.Facade) {
	_, claims, _ := jwtauth.FromContext(request.Context())
	tenantId := claims["tenant_id"]
	updateRequest := models.UpdateTenantLicenseRequestBody{}
	json.NewDecoder(request.Body).Decode(&updateRequest)
	requestId, _ := uuid.NewUUID()
	command := models.UpdateTenantLicenseRequest{
		RequestId:           requestId.String(),
		TenantId:            tenantId.(string),
		NewLicenseType:      updateRequest.NewLicenseType,
		NewLicenseStartDate: updateRequest.NewLicenseStartDate,
		NewLicenseEndDate:   updateRequest.NewLicenseEndDate,
	}
	if err := facade.CommandBus().Send(context.Background(), &command); err != nil {
		panic(err)
	}
	result := database.GetCommandsDbConn().Create(&command)
	if result.Error != nil {
		logrus.Errorf("Unable to save update tenant license request, %s", result.Error)
	}
	writer.WriteHeader(Ok)
}
