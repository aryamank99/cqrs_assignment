package database

import (
	"github.com/sirupsen/logrus"
	"tenantMicroserviceAssignment/models"
)

type TenantDbServiceImpl struct{}

func (s TenantDbServiceImpl) AddNewTenant(newTenant models.Tenant) *models.Tenant {
	result := GetStoreDbConn().Create(&newTenant)
	if result.Error != nil {
		logrus.Errorf("Unable to save tenant, %s", result.Error)
		return nil
	}
	return s.GetTenant(newTenant.ID)
}

func (s TenantDbServiceImpl) GetTenant(tenantId string) *models.Tenant {
	var tenant models.Tenant
	result := GetStoreDbConn().First(&tenant, "id = ?", tenantId)
	if result.Error != nil {
		logrus.Errorf("Unable to get tenant, %s", result.Error)
		return nil
	}
	return &tenant
}

func (s TenantDbServiceImpl) GetAllTenants() *[]models.Tenant {
	var tenants []models.Tenant
	result := GetStoreDbConn().Find(&tenants)
	if result.Error != nil {
		logrus.Errorf("Unable to get all tenants, %s", result.Error)
		return nil
	}
	return &tenants
}

func (s TenantDbServiceImpl) UpdateTenantLicense(request models.UpdateTenantLicenseRequest) *models.Tenant {
	// 1: get the tenant first
	tenant := s.GetTenant(request.TenantId)
	if tenant == nil {
		logrus.Errorf("unable to get tenant")
		return nil
	}
	// 2: update tenant
	result := GetStoreDbConn().Model(&tenant).Updates(models.Tenant{
		LicenseType:      request.NewLicenseType,
		LicenseStartDate: request.NewLicenseStartDate,
		LicenseEndDate:   request.NewLicenseEndDate,
	})
	if result.Error != nil {
		logrus.Errorf("Unable to update tenant, %s", result.Error)
		return nil
	}
	// 3: return updated tenant
	return s.GetTenant(request.TenantId)
}

func (s TenantDbServiceImpl) UpdateTenantStatus(request models.UpdateTenantStatusRequest) *models.Tenant {
	// 1: get the tenant first
	tenant := s.GetTenant(request.TenantId)
	if tenant == nil {
		logrus.Errorf("unable to get tenant")
		return nil
	}
	// 2: update tenant
	result := GetStoreDbConn().Model(&tenant).Updates(models.Tenant{
		IsEnabled: request.IsEnabled,
	})
	if result.Error != nil {
		logrus.Errorf("Unable to update tenant status, %s", result.Error)
		return nil
	}
	// 3: return updated tenant
	return s.GetTenant(request.TenantId)
}
