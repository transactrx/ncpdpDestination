package pbmerrors

type UnableToParseClaimError struct {
}

func (e UnableToParseClaimError) Error() string {
	return ErrorCode.TRX02.Message
}
func (e UnableToParseClaimError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX02
}
