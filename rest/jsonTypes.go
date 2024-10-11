package rest

type Health struct {
	Status   string   `json:"status"`
	Messages []string `json:"messages"`
}

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
