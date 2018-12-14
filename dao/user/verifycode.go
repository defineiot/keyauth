package user

// VerifyCode code
type VerifyCode struct {
	ID           int64  `json:"id"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
	Code         int    `json:"code"`
	CreateAt     int64  `json:"create_at"`
	ExpireAt     int64  `json:"expire_at"`
	Status       int    `json:"status"`
	Type         int    `json:"type"`
}
