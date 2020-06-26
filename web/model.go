package main

type APIBody struct {
	URL string `json:"url"`
	Method string `json:"method"`
	ReqBody string `json:"req_body"`
}

type Err struct {
	Error string `json:"error"`
	ErrorCode string `json:"error_code"`
}


var (
	ErrorBadRequest = Err{
		Error:     "bad request for api",
		ErrorCode: "001",
	}
	ErrorReqBodyParseFailed = Err{
		Error:     "request body parse failed",
		ErrorCode: "002",
	}
	ErrorInternalError = Err{
		Error:     "Internal server error",
		ErrorCode: "003",
	}
)