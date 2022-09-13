package payloads

type Authenticate struct {
	ID string `json:"identity"`

	Secret string `json:"secret"`
}
