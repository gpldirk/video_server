package model

type Err struct {
	Error string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrResponse struct {
	HttpSC int
	Error Err
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{
		HttpSC: 400,
		Error: Err{
			Error: "Request Body Parse Failed",
			ErrorCode: "001",
		}}

	ErrorUserAuthFailed = ErrResponse{
		HttpSC: 401,
		Error: Err{
			Error: "User Authentication Failed",
			ErrorCode: "002",
		}}

	ErrorDBError = ErrResponse{
		HttpSC: 500,
		Error: Err{
			Error:     "DB operation error",
			ErrorCode: "003",
	}}

	 ErrorInternalError = ErrResponse{
	 	HttpSC: 500,
	 	Error: Err{
			Error:     "Internal server error",
			ErrorCode: "004",
		},
	 }
)
