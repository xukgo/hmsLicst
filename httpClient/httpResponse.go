package httpClient

type HttpResponse struct {
	Error      error
	Data       []byte
	Header     map[string][]string
	StatusCode int    // e.g. 200
	Status     string // e.g. "200 OK"
}

func NewHttpResponse(err error, code int, state string, data []byte, Header map[string][]string) HttpResponse {
	model := HttpResponse{
		Error:      err,
		Data:       data,
		Header:     Header,
		StatusCode: code,
		Status:     state,
	}
	return model
}

func NewErrorHttpResponse(err error) HttpResponse {
	model := HttpResponse{
		Error: err,
	}
	return model
}
