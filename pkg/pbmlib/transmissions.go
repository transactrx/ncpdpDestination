package pbmlib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	SEGMENT_SEPARATOR = "\u001E"
	GROUP_SEPARATOR   = "\u001D"
	FIELD_SEPARATOR   = "\u001C"
	ETX               = "\u0003"
)

type Claim struct {
	TransmissionId string `json:"transmissionId"`
	InstanceId     string `json:"instanceId"`
	TenantId       string `json:"tenantId"`
	UnitName       string `json:"unitName"`
	ClientId       string `json:"clientId"`

	Elapsed         string        `json:"elapsed"`
	From            string        `json:"from"`
	RouteAddress    string        `json:"routeAddress"`
	RouteName       string        `json:"routeName"`
	Created         string        `json:"created"`
	RxNumbers       []interface{} `json:"rxNumbers"`
	Ncpdp           string        `json:"ncpdp"`
	CallOrigin      int           `json:"callOrigin"`
	TransactionData struct {
		NcpdpData string `json:"ncpdpData"`
	} `json:"transactionData"`
	ReplyTo  string    `json:"replyTo"`
	TimeRcvd time.Time `json:"timeRcvd"`
}

type Response struct {
	TransmissionId  string        `json:"transmissionId"`
	InstanceId      string        `json:"instanceId"`
	TenantId        string        `json:"tenantId"`
	UnitName        string        `json:"unitName"`
	ClientId        string        `json:"clientId"`
	Elapsed         string        `json:"elapsed"`
	From            string        `json:"from"`
	RouteAddress    string        `json:"routeAddress"`
	RouteName       string        `json:"routeName"`
	Created         string        `json:"created"`
	RxNumbers       []interface{} `json:"rxNumbers"`
	Ncpdp           string        `json:"ncpdp"`
	CallOrigin      int           `json:"callOrigin"`
	ReturnCode      string        `json:"returnCode"`
	StatusCode      string        `json:"statusCode"`
	Reader          string        `json:"reader"`
	ContextData     string        `json:"contextData"`
	TransactionData struct {
		NcpdpData string `json:"ncpdpData"`
	} `json:"transactionData"`
	ErrorData *ErrorInfo `json:"errorData"`
}

func (response *Response) BuildResponseSuccess(claim Claim, startTime time.Time, ncpdpResponse []byte) {

	response.TransmissionId = claim.TransmissionId
	response.InstanceId = claim.InstanceId
	response.TenantId = claim.TenantId
	response.Elapsed = fmt.Sprintf("%f", (time.Now().Sub(startTime).Seconds()))
	now := time.Now().UTC()
	formattedTime := now.Format("2006-01-02T15:04:05.9999999Z")
	response.Created = formattedTime
	response.RouteAddress = claim.RouteAddress
	response.RouteName = claim.RouteName
	response.RxNumbers = claim.RxNumbers
	response.Ncpdp = claim.Ncpdp
	response.CallOrigin = claim.CallOrigin
	response.From, _ = os.Hostname()
	response.ReturnCode = "1" // ok status
	response.StatusCode = "200"
	response.TransactionData.NcpdpData = string(ncpdpResponse)
}

func (response *Response) BuildResponseError(claim Claim, errorCode ErrorInfo, startTime time.Time) {

	response.TransmissionId = claim.TransmissionId
	response.InstanceId = claim.InstanceId
	response.Elapsed = fmt.Sprintf("%f", (time.Now().Sub(startTime).Seconds()))
	now := time.Now().UTC()
	formattedTime := now.Format("2006-01-02T15:04:05.9999999Z")
	response.Created = formattedTime
	response.RouteAddress = claim.RouteAddress
	response.RouteName = claim.RouteName
	response.RxNumbers = claim.RxNumbers
	response.Ncpdp = claim.Ncpdp
	response.CallOrigin = claim.CallOrigin
	response.From, _ = os.Hostname()
	response.ReturnCode = "3" // failure
	response.StatusCode = "200"
	response.TransactionData.NcpdpData = GenerateError([]byte(claim.TransactionData.NcpdpData), errorCode)
}

func (response *Response) ToJSON() []byte {

	respJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("What the heck!!")
	}
	return respJson
}

// func GenerateError(responseBuffer []byte,errorCode string) string {
func GenerateError(responseBuffer []byte, errorCode ErrorInfo) string {

	ncpdpClaim := string(responseBuffer)
	errorMessage := errorCode.Code + ": " + errorCode.Message
	response := "D0" + ncpdpClaim[8:8+2] + ncpdpClaim[20:20+1] + "R" + ncpdpClaim[21:21+2] + ncpdpClaim[23:23+15] + time.Now().Format("20060102")
	//Build Response Message Segment AM20
	response += SEGMENT_SEPARATOR
	response += FIELD_SEPARATOR + "AM20"
	response += FIELD_SEPARATOR + "F4" + errorMessage
	tmp := ncpdpClaim[20 : 20+1]
	numberOfClaims, err := strconv.Atoi(tmp)
	//Build Response Status Segment for each claim AM21
	if err != nil {
		numberOfClaims = 1
	}
	for i := 0; i < numberOfClaims; i++ {
		response += GROUP_SEPARATOR
		response += SEGMENT_SEPARATOR
		response += FIELD_SEPARATOR + "AM21"
		response += FIELD_SEPARATOR + "ANR"
		response += FIELD_SEPARATOR + "FA2"
		response += FIELD_SEPARATOR + "FB85"
		response += FIELD_SEPARATOR + "FBNN"
	}
	response += ETX
	return response
}
