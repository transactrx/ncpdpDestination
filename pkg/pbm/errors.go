package pbm

type ErrorInfo struct {
	Message     string
	HttpCode    string
	Code        string
	Description string
	Causes      string
}

type ErrorCodes struct {
	TRX00 ErrorInfo
	TRX01 ErrorInfo
	TRX02 ErrorInfo
	TRX03 ErrorInfo
	TRX04 ErrorInfo
	TRX05 ErrorInfo
	TRX06 ErrorInfo
	TRX07 ErrorInfo
	TRX08 ErrorInfo
	TRX09 ErrorInfo
	TRX10 ErrorInfo
	// add more codes here
	TRX9999 ErrorInfo
}

var ErrorCode = ErrorCodes{

	TRX00: ErrorInfo{
		Message:     "No PBMError",
		HttpCode:    "200",
		Code:        "TRX00",
		Description: "NA",
		Causes:      "NA",
	},

	TRX01: ErrorInfo{
		Message:     "Unable to Parse Claim",
		HttpCode:    "400",
		Code:        "TRX01",
		Description: "This error occurs when the system cannot parse the claim data provided in the request.",
		Causes:      "Possible Causes: Invalid or improperly formatted claim data, missing required fields, or a mismatch between the data format and the expected format.",
	},

	TRX02: ErrorInfo{
		Message:     "Unable to Connect to Endpoint",
		HttpCode:    "500",
		Code:        "TRX02",
		Description: "This error indicates that the system was unable to establish a connection to the endpoint (URL) specified for claim processing.",
		Causes:      "Possible Causes: Network issues, endpoint URL misconfiguration, or the endpoint is temporarily unavailable.",
	},

	TRX03: ErrorInfo{
		Message:     "Endpoint Authentication Failure",
		HttpCode:    "401",
		Code:        "TRX03",
		Description: "This error occurs when the system fails to authenticate with the endpoint due to invalid credentials or authentication method.",
		Causes:      "Possible Causes: Incorrect authentication credentials, expired tokens, or misconfigured authentication settings.",
	},
	TRX04: ErrorInfo{
		Message:     "Invalid Request Format",
		HttpCode:    "400",
		Code:        "TRX04",
		Description: "This error is triggered when the request sent to the endpoint is in an incorrect format that the endpoint cannot process.",
		Causes:      "Possible Causes: Incorrect request headers, unsupported content type, or missing required request parameters.",
	},
	TRX05: ErrorInfo{
		Message:     "Time-out Waiting for Response",
		HttpCode:    "504",
		Code:        "TRX05",
		Description: "This error occurs when the system does not receive a response from the endpoint within the expected time frame.",
		Causes:      "Possible Causes: Slow network, endpoint server overload, or an unresponsive endpoint.",
	},
	TRX06: ErrorInfo{
		Message:     "Unable to Parse Response",
		HttpCode:    "500",
		Code:        "TRX06",
		Description: "This error indicates that the system could not parse the response received from the endpoint after sending the request.",
		Causes:      "Possible Causes: Incorrect response format, unexpected data structure, or a problem with the endpoint's response.",
	},
	TRX07: ErrorInfo{
		Message:     "Claim Processing PBMError",
		HttpCode:    "500",
		Code:        "TRX07",
		Description: "This error code can be used to represent any generic error that occurs during the claim processing at the endpoint.",
		Causes:      "Possible Causes: Issues specific to the claim processing logic at the endpoint, such as business rule violations or data inconsistencies.",
	},
	TRX08: ErrorInfo{
		Message:     "Endpoint Unavailable",
		HttpCode:    "503",
		Code:        "TRX08",
		Description: "This error is used when the endpoint is temporarily or permanently unavailable.",
		Causes:      "Possible Causes: Endpoint maintenance, downtime, or the endpoint URL no longer exists.",
	},
	TRX09: ErrorInfo{
		Message:     "Request Authorization Failure",
		HttpCode:    "403",
		Code:        "TRX09",
		Description: "This error occurs when the request lacks the necessary authorization to access the endpoint.",
		Causes:      "Possible Causes: Missing or insufficient authorization headers, tokens, or permissions.",
	},
	TRX10: ErrorInfo{
		Message:     "PBMError While Sending POST Request",
		HttpCode:    "403",
		Code:        "TRX10",
		Description: "This error occurs when there was an issue while sending a POST request to the endpoint.",
		Causes:      "Possible Causes: This error can be caused by missing or insufficient authorization headers, tokens, or permissions. It may also indicate a problem on the server's side.",
	},
	TRX9999: ErrorInfo{
		Message:     "Host Processing PBMError",
		HttpCode:    "500",
		Code:        "TRX9999",
		Description: "Unknown error",
		Causes:      "Possible Causes: ?",
	},
}
