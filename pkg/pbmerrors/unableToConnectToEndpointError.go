package pbmerrors

type UnableToConnectToEndpoint struct{}

func (e UnableToConnectToEndpoint) Error() string {
	return ErrorCode.TRX02.Message
}
func (e UnableToConnectToEndpoint) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX02
}
