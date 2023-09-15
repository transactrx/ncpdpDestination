package pbmerrors

type RequestAuthorizationFailureError struct{}

func (e RequestAuthorizationFailureError) Error() string {
	return ErrorCode.TRX09.Message
}
func (e RequestAuthorizationFailureError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX09
}
