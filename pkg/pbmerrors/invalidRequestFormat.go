package pbmerrors

type InvalidRequestFormat struct{}

func (e InvalidRequestFormat) Error() string {
	return ErrorCode.TRX04.Message
}
func (e InvalidRequestFormat) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX04
}
