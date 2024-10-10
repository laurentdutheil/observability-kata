package rest

type Health struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages"`
}
