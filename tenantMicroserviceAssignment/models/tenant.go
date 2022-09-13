package models

import "github.com/gofrs/uuid"

type Tenant struct {
	ID               string `json:"_id"`
	Name             string `json:"first_name"`
	TotalUserCount   string `json:"total_user_count"`
	IsEnabled        string `json:"is_enabled"`
	LicenseType      string `json:"license_type"`
	LicenseStartDate string `json:"license_start_date"`
	LicenseEndDate   string `json:"license_end_date"`
}

type UpdateTenantCredentialsRequestBody struct {
	NewPassword string `json:"new_password"`
}

type UpdateTenantCredentialsRequest struct {
	RequestId   string `gorm:"primaryKey"`
	TenantId    string `json:"tenant_id"`
	NewPassword string `json:"new_password"`
}

type UpdateTenantLicenseRequestBody struct {
	RequestId           string `gorm:"primaryKey"`
	NewLicenseType      string `json:"new_license_type"`
	NewLicenseStartDate string `json:"new_license_start_date"`
	NewLicenseEndDate   string `json:"new_license_end_date"`
}

type UpdateTenantLicenseRequest struct {
	RequestId           string `gorm:"primaryKey"`
	TenantId            string `json:"tenant_id"`
	NewLicenseType      string `json:"new_license_type"`
	NewLicenseStartDate string `json:"new_license_start_date"`
	NewLicenseEndDate   string `json:"new_license_end_date"`
}
type UpdateTenantStatusRequestBody struct {
	IsEnabled string `json:"is_enabled"`
}

type UpdateTenantStatusRequest struct {
	RequestId string `gorm:"primaryKey"`
	TenantId  string `json:"tenant_id"`
	IsEnabled string `json:"is_enabled"`
}

type CreateTenantRequest struct {
	RequestId        string `gorm:"primaryKey"`
	Name             string `json:"first_name"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	TotalUserCount   string `json:"total_user_count"`
	IsEnabled        string `json:"is_enabled"`
	LicenseType      string `json:"license_type"`
	LicenseStartDate string `json:"license_start_date"`
	LicenseEndDate   string `json:"license_end_date"`
}

type CreateTenantRequestBody struct {
	Name             string `json:"first_name"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	TotalUserCount   string `json:"total_user_count"`
	IsEnabled        string `json:"is_enabled"`
	LicenseType      string `json:"license_type"`
	LicenseStartDate string `json:"license_start_date"`
	LicenseEndDate   string `json:"license_end_date"`
}

func NewTenant(request CreateTenantRequest) *Tenant {
	tenantId, err := uuid.NewV4()
	if err != nil {
		// return bad request in case of error
		panic("error generating tenant uuid")
	}
	return &Tenant{
		ID:               tenantId.String(),
		Name:             request.Name,
		TotalUserCount:   request.TotalUserCount,
		IsEnabled:        request.IsEnabled,
		LicenseType:      request.LicenseType,
		LicenseStartDate: request.LicenseStartDate,
		LicenseEndDate:   request.LicenseEndDate,
	}
}

type TenantProfileResponse struct {
	Profile Tenant `json:"profile"`
	Token   string `json:"token"`
}
