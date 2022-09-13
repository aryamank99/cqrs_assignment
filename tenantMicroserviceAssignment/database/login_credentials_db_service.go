package database

import "tenantMicroserviceAssignment/models"

type LoginCredentialDbService interface {
	AddNewCredentials(credentials models.LoginCredentials) *models.LoginCredentials
	GetCredentials(email string) *models.LoginCredentials
	UpdateTenantCredentials(string, string) *models.LoginCredentials
}
