package database

import (
	"github.com/sirupsen/logrus"
	"tenantMicroserviceAssignment/models"
)

type LoginCredentialsDbServiceImpl struct{}

func (s LoginCredentialsDbServiceImpl) AddNewCredentials(credentials models.LoginCredentials) *models.LoginCredentials {
	result := GetStoreDbConn().Create(&credentials)
	if result.Error != nil {
		logrus.Errorf("Unable to save credentials, %s", result.Error)
		return nil
	}
	return s.GetCredentials(credentials.Username)
}

func (s LoginCredentialsDbServiceImpl) GetCredentials(username string) *models.LoginCredentials {
	var credentials models.LoginCredentials
	result := GetStoreDbConn().First(&credentials, "username = ?", username)
	if result.Error != nil {
		logrus.Errorf("Unable to get login credentials, %s", result.Error)
		return nil
	}
	return &credentials
}

func (s LoginCredentialsDbServiceImpl) GetCredentialsById(id string) *models.LoginCredentials {
	var credentials models.LoginCredentials
	result := GetStoreDbConn().First(&credentials, "id = ?", id)
	if result.Error != nil {
		logrus.Errorf("Unable to get login credentials, %s", result.Error)
		return nil
	}
	return &credentials
}

func (s LoginCredentialsDbServiceImpl) UpdateTenantCredentials(tenantId string, newPassword string) *models.LoginCredentials {
	// 1: get the tenant first
	tenant := s.GetCredentialsById(tenantId)
	if tenant == nil {
		logrus.Errorf("unable to get tenant")
		return nil
	}
	// 2: update tenant
	result := GetStoreDbConn().Model(&tenant).Updates(models.LoginCredentials{
		PasswordHash: newPassword,
	})
	if result.Error != nil {
		logrus.Errorf("Unable to update tenant credentials, %s", result.Error)
		return nil
	}
	// 3: return updated tenant
	return s.GetCredentialsById(tenantId)
}
