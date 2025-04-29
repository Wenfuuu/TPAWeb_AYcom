package messaging

type VerificationEmailPayload struct {
    Email string `json:"email"`
    Code  string `json:"code"`
}

type EmailPayload struct {
	Email string `json:"email"`
	Word string `json:"word"`
}