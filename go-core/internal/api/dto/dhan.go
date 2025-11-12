package dto

// DhanRenewTokenRequest represents the request to renew Dhan access token
type DhanRenewTokenRequest struct {
	AccessToken  string `json:"access_token" validate:"required"`
	DhanClientID string `json:"dhan_client_id" validate:"required"`
}

// DhanRenewTokenResponse represents the response for renewing token
type DhanRenewTokenResponse struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	ExpiryTime  string `json:"expiry_time"`
}

// DhanGenerateConsentRequest represents the request to generate consent
type DhanGenerateConsentRequest struct {
	DhanClientID string `json:"dhan_client_id" validate:"required"`
}

// DhanGenerateConsentResponse represents the response for generating consent
type DhanGenerateConsentResponse struct {
	ConsentAppID     string `json:"consent_app_id"`
	ConsentAppStatus string `json:"consent_app_status"`
	Status           string `json:"status"`
	LoginURL         string `json:"login_url"`
	CallbackURL      string `json:"callback_url"` // Suggested callback URL for Dhan redirect configuration
}

// DhanConsumeConsentRequest represents the request to consume consent
type DhanConsumeConsentRequest struct {
	TokenID string `json:"token_id" validate:"required"`
}

// DhanConsumeConsentResponse represents the response for consuming consent
type DhanConsumeConsentResponse struct {
	DhanClientID         string `json:"dhan_client_id"`
	DhanClientName       string `json:"dhan_client_name"`
	DhanClientUcc        string `json:"dhan_client_ucc"`
	GivenPowerOfAttorney bool   `json:"given_power_of_attorney"`
	AccessToken          string `json:"access_token"`
	ExpiryTime           string `json:"expiry_time"`
}

// DhanBrokerConfigResponse represents the broker configuration for a user
type DhanBrokerConfigResponse struct {
	Configured     bool   `json:"configured"`
	HasCredentials bool   `json:"has_credentials"` // Whether API key/secret are configured
	DhanClientID   string `json:"dhan_client_id,omitempty"`
	DhanClientName string `json:"dhan_client_name,omitempty"`
	ExpiryTime     string `json:"expiry_time,omitempty"`
}

// DhanSaveTokenRequest represents the request to save access token directly
type DhanSaveTokenRequest struct {
	AccessToken  string `json:"access_token" validate:"required"`
	DhanClientID string `json:"dhan_client_id" validate:"required"`
}

// DhanSaveCredentialsRequest represents the request to save API key and secret
type DhanSaveCredentialsRequest struct {
	APIKey       string `json:"api_key" validate:"required"`
	APISecret    string `json:"api_secret" validate:"required"`
	DhanClientID string `json:"dhan_client_id" validate:"required"` // Required for generate-consent API
}
