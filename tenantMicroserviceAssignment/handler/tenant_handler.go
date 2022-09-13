package handler

import (
	"crypto/sha512"
	"encoding/base64"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"math/rand"
	"tenantMicroserviceAssignment/database"
	"tenantMicroserviceAssignment/models"
)

type TenantServiceImpl struct{}

func (s TenantServiceImpl) AddNewTenant(request models.CreateTenantRequest, tenantDbService database.TenantDbService, loginCredentialsDbService database.LoginCredentialDbService) *models.TenantProfileResponse {
	// 1: save tenant
	newTenant := models.NewTenant(request)
	// 2: save login credentials
	salt := generateRandomSalt(16)
	hashedPassword := hashPassword(request.Password, salt)
	credentials := models.LoginCredentials{
		ID:           newTenant.ID,
		Username:     request.Username,
		PasswordHash: hashedPassword,
		Salt:         salt,
	}
	loginCredentialsDbService.AddNewCredentials(credentials)
	savedTenant := tenantDbService.AddNewTenant(*newTenant)
	token := generateUserToken(savedTenant.ID)
	return &models.TenantProfileResponse{
		Profile: *savedTenant,
		Token:   token,
	}
}

func (s TenantServiceImpl) GetAllTenants(tenantDbService database.TenantDbService) *[]models.Tenant {
	return tenantDbService.GetAllTenants()
}

func (s TenantServiceImpl) GetTenantById(tenantId string, tenantDbService database.TenantDbService) *models.Tenant {
	return tenantDbService.GetTenant(tenantId)
}

func (s TenantServiceImpl) UpdateTenantCredentials(request models.UpdateTenantCredentialsRequest, loginCredentialsDbService database.LoginCredentialDbService) *models.LoginCredentials {
	return loginCredentialsDbService.UpdateTenantCredentials(request.TenantId, request.NewPassword)
}

func (s TenantServiceImpl) UpdateTenantLicense(request models.UpdateTenantLicenseRequest, tenantDbService database.TenantDbServiceImpl) *models.Tenant {
	return tenantDbService.UpdateTenantLicense(request)
}

func (s TenantServiceImpl) UpdateTenantStatus(request models.UpdateTenantStatusRequest, tenantDbService database.TenantDbServiceImpl) *models.Tenant {
	return tenantDbService.UpdateTenantStatus(request)
}

func hashPassword(password string, salt string) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	// Create sha-512 hasher
	var sha512Hasher = sha512.New()
	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)
	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)
	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	// Convert the hashed password to a base64 encoded string
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)
	return base64EncodedPasswordHash
}

func generateRandomSalt(saltSize int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, saltSize)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateUserToken(tenantId string) string {
	logrus.Infof("Generaring token for tenantId: %s", tenantId)
	var tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"tenant_id": tenantId,
	})
	return tokenString
}
