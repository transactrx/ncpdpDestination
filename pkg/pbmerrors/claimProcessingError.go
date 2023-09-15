package pbmerrors

type ClaimProcessingError struct{}

func (e ClaimProcessingError) Error() string {
	return ErrorCode.TRX07.Message
}

func (e ClaimProcessingError) GetErrorInfo() ErrorInfo {
	return ErrorCode.TRX07
}
