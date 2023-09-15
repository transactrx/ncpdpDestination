package pbmerrors

type EndpointUnavailableError struct{}

func (e EndpointUnavailableError) Error() string {
	return ErrorCode.TRX08.Message
}

func (e EndpointUnavailableError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX08
}
