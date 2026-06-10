package user_device

type RegisterDeviceRequest struct {
	FcmToken string `json:"fcm_token" validate:"required"`
}

type UnregisterDeviceRequest struct {
	FcmToken string `json:"fcm_token" validate:"required"`
}
