package messaging

type VerificationEmailPayload struct {
    Email string `json:"email"`
    Code  string `json:"code"`
}