package pbmerrors

type HostProcessingError struct{}

func (e HostProcessingError) Error() string {
	return ErrorCode.TRX9999.Message
}
func (e HostProcessingError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX9999
}
