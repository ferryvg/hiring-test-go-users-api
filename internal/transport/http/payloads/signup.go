package payloads

type SignUp struct {
	ID string `json:"identity"`

	Secret string `json:"secret"`
}
