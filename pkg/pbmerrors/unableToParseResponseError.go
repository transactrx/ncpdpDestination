package pbmerrors

type UnableToParseResponseError struct{}

func (e UnableToParseResponseError) Error() string {
	return ErrorCode.TRX06.Message
}
func (e UnableToParseResponseError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX06
}
