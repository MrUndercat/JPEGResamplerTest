package response

type OK struct {
	Code int `json:"code"`
	Body struct {
		Time   int64 `json:"time"`
		Cached bool  `json:"cached"`
	} `json:"body"`
}

type Error struct {
	Code int `json:"code"`
	Body struct {
		Error string `json:"error"`
	} `json:"body"`
}
