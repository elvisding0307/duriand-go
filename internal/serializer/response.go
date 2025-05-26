package serializer

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, msg string, data interface{}) Response {
	return Response{Code: code, Msg: msg, Data: data}
}

func NewSuccessResponse(data interface{}) Response {
	return NewResponse(0, "Success", data)
}

func NewErrorResponse(errorCode int, msg string) Response {
	return NewResponse(errorCode, msg, nil)
}
