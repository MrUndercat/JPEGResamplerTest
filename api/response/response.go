package response

type ResponseOK struct {
	Code int `json:"code"`
	Body struct {
		Time   int  `json:"time"`
		Cached bool `json:"cached"`
	} `json:"body"`
}

type ResponseError struct {
	Code int `json:"code"`
	Body struct {
		Error string `json:"error"`
	} `json:"body"`
}
