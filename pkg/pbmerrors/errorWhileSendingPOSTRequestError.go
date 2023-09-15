package pbmerrors

type ErrorWhileSendingPOSTRequestError struct{}

func (e ErrorWhileSendingPOSTRequestError) Error() string {
	return ErrorCode.TRX10.Message
}

func (e ErrorWhileSendingPOSTRequestError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX10
}
