package pbmerrors

type TimeoutWaitingForResponseError struct{}

func (e TimeoutWaitingForResponseError) Error() string {
	return ErrorCode.TRX05.Message
}
func (e TimeoutWaitingForResponseError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX05
}
