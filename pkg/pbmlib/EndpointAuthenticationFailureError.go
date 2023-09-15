package pbmlib

type EndpointAuthenticationFailureError struct{}

func (e EndpointAuthenticationFailureError) Error() string {
	return ErrorCode.TRX03.Message
}
func (e EndpointAuthenticationFailureError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX03
}
