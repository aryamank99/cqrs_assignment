package database

import "tenantMicroserviceAssignment/models"

type TenantDbService interface {
	AddNewTenant(models.Tenant) *models.Tenant
	GetTenant(string) *models.Tenant
	GetAllTenants() *[]models.Tenant
	UpdateTenantLicense(models.UpdateTenantLicenseRequest) *models.Tenant
	UpdateTenantStatus(models.UpdateTenantStatusRequest) *models.Tenant
}
