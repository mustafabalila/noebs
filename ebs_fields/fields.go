package ebs_fields

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
	"time"
)

type IsAliveFields struct {
	CommonFields
}
type WorkingKeyFields struct {
	CommonFields
}

type BalanceFields struct {
	CommonFields
	CardInfoFields
}
type MiniStatementFields struct {
	CommonFields
	CardInfoFields
}
type ChangePINFields struct {
	CommonFields
	CardInfoFields
	NewPIN string `json:"newPIN" binding:"required"`
}

type CardTransferFields struct {
	CommonFields
	CardInfoFields
	AmountFields
	ToCard string `json:"toCard" binding:"required"`
}

type PurchaseFields struct {
	WorkingKeyFields
	CardInfoFields
	AmountFields
}

type BillPaymentFields struct {
	CommonFields
	CardInfoFields
	AmountFields
	BillerFields
}

type CashInFields struct {
	PurchaseFields
}
type CashOutFields struct {
	PurchaseFields
}
type RefundFields struct {
	PurchaseFields
	OriginalSTAN int `json:"originalSystemTraceAuditNumber" binding:"required"`
}

type PurchaseWithCashBackFields struct {
	PurchaseFields
}
type ReverseFields struct {
	PurchaseFields
}

type BillInquiryFields struct {
	CommonFields
	CardInfoFields
	AmountFields
	BillerFields
}

type CommonFields struct {
	SystemTraceAuditNumber int    `json:"systemTraceAuditNumber,omitempty" binding:"required" form:"systemTraceAuditNumber"`
	TranDateTime           string `json:"tranDateTime,omitempty" binding:"required" form:"tranDateTime"`
	TerminalID             string `json:"terminalId,omitempty" binding:"required,len=8" form:"terminalId"`
	ClientID               string `json:"clientId,omitempty" binding:"required" form:"clientId"`
}

func (c *CommonFields) sendRequest(f []byte) {
	panic("implement me!")
}

type CardInfoFields struct {
	Pan     string `json:"PAN" binding:"required" form:"PAN"`
	Pin     string `json:"PIN" binding:"required" form:"PIN"`
	Expdate string `json:"expDate" binding:"required" form:"expDate"`
}

type AmountFields struct {
	TranAmount       float32 `json:"tranAmount" binding:"required" form:"tranAmount"`
	TranCurrencyCode string  `json:"tranCurrencyCode" form:"tranCurrencyCode"`
}

type BillerFields struct {
	PersonalPaymentInfo string `json:"personalPaymentInfo" binding:"required" form:"personalPaymentInfo"`
	PayeeID             string `json:"payeeId" binding:"required" form:"payeeId"`
}

type PayeesListFields struct {
	CommonFields
}

func iso8601(fl validator.FieldLevel) bool {

	dateLayout := time.RFC3339
	_, err := time.Parse(dateLayout, fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

type GenericEBSResponseFields struct {
	TerminalID             string `json:"terminalId,omitempty"`
	SystemTraceAuditNumber int    `json:"systemTraceAuditNumber,omitempty"`
	ClientID               string `json:"clientId,omitempty"`
	PAN                    string `json:"PAN,omitempty"`

	ServiceID        string  `json:"serviceId,omitempty"`
	TranAmount       float32 `json:"tranAmount,omitempty"`
	PhoneNumber      string  `json:"phoneNumber,omitempty"`
	FromAccount      string  `json:"fromAccount,omitempty"`
	ToAccount        string  `json:"toAccount,omitempty"`
	FromCard         string  `json:"fromCard,omitempty"`
	ToCard           string  `json:"toCard,omitempty"`
	OTP              string  `json:"otp,omitempty"`
	OTPID            string  `json:"otpId,omitempty"`
	TranCurrencyCode string  `json:"tranCurrencyCode,omitempty"`
	EBSServiceName   string  `json:"-,omitempty"`
	WorkingKey       string  `json:"workingKey,omitempty" gorm:"-"`
	PayeeID          string  `json:"payeeId,omitempty"`

	// Consumer fields
	PubKeyValue string `json:"pubKeyValue,omitempty" form:"pubKeyValue"`
	UUID        string `json:"UUID,omitempty" form:"UUID"`

	ResponseMessage      string   `json:"responseMessage,omitempty"`
	ResponseStatus       string   `json:"responseStatus,omitempty"`
	ResponseCode         int      `json:"responseCode"`
	ReferenceNumber      string   `json:"referenceNumber,omitempty"`
	ApprovalCode         string   `json:"approvalCode,omitempty"`
	VoucherNumber        int      `json:"voucherNumber,omitempty"`
	MiniStatementRecords string   `json:"miniStatementRecords,omitempty"`
	DisputeRRN           string   `json:"DisputeRRN,omitempty"`
	AdditionalData       string   `json:"additionalData,omitempty"`
	TranDateTime         string   `json:"tranDateTime,omitempty"`
	TranFee              *float32 `json:"tranFee,omitempty"`

	AdditionalAmount *float32 `json:"additionalAmount,omitempty"`
	AcqTranFee       *float32 `json:"acqTranFee,omitempty"`
	IssTranFee       *float32 `json:"issuerTranFee,omitempty"`
	TranCurrency     string   `json:"tranCurrency,omitempty"`

	// QR payment fields
	MerchantID  string `json:"merchantID,omitempty"`
	GeneratedQR string `json:"generatedQR,omitempty"`
}

type ImportantEBSFields struct {
}

// you have to update this to account for the non-db-able fields
type EBSParserFields struct {
	EBSMapFields
	GenericEBSResponseFields
}

// To allow Redis to use this struct directly in marshaling
func (p *EBSParserFields) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)

}

// To allow Redis to use this struct directly in marshaling
func (p *EBSParserFields) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

// special case to handle ebs non-DB-able fields e.g., hashmaps and other complex types
type EBSMapFields struct {
	// these
	Balance     map[string]interface{} `json:"balance,omitempty"`
	PaymentInfo string                 `json:"paymentInfo,omitempty"`
	BillInfo    map[string]interface{} `json:"billInfo,omitempty"`
}

type ConsumerSpecificFields struct {
	UUID            string  `json:"UUID" form:"UUID" binding:"required,len=36"`
	Mbr             string  `json:"mbr,omitempty" form:"mbr"`
	Ipin            string  `json:"IPIN" form:"IPIN" binding:"required"`
	FromAccountType string  `json:"fromAccountType" form:"fromAccountType"`
	ToAccountType   string  `json:"toAccountType" form:"toAccountType"`
	AccountCurrency string  `json:"accountCurrency" form:"accountCurrency"`
	AcqTranFee      float32 `json:"acqTranFee" form:"acqTranFee"`
	IssuerTranFee   float32 `json:"issuerTranFee" form:"issuerTranFee"`

	// billers
	BillInfo string `json:"billInfo" form:"billInfo"`
	Payees   string `json:"payees" form:"payees"`

	// tran time
	OriginalTranUUID     string `json:"originalTranUUID" form:"originalTranUUID"`
	OriginalTranDateTime string `json:"originalTranDateTime" form:"originalTranDateTime"`

	// User settings
	Username     string `json:"userName" Form:"userName"`
	UserPassword string `json:"userPassword" form:"userPassword"`

	// Entities
	EntityType  string `json:"entityType" form:"entityType"`
	EntityId    string `json:"entityId" form:"entityId"`
	EntityGroup string `json:"entityGroup" form:"entityGroup"`
	PubKeyValue string `json:"pubKeyValue" form:"pubKeyValue"`
	Email       string `json:"email" form:"email"`
	ExtraInfo   string `json:"extraInfo" form:"extraInfo"`
	//.. omitted fields

	PhoneNo                string `json:"phoneNo" form:"phoneNo"`
	NewIpin                string `json:"newIPIN" form:"newIPIN"`
	NewUserPassword        string `json:"newUserPassword" form:"newUserPassword"`
	SecurityQuestion       string `json:"securityQuestion" form:"securityQuestion"`
	SecurityQuestionAnswer string `json:"securityQuestionAnswer" form:"securityQuestionAnswer"`
	AdminUserName          string `json:"adminUserName" form:"adminUserName"`

	// other fields
	OriginalTransaction map[string]interface{} `json:"originalTransaction" form:"originalTransaction"`
	OriginalTranType    string                 `json:"originalTranType" form:"originalTranType"`
}

// MaskPAN returns the last 4 digit of the PAN. We shouldn't care about the first 6
func (res *GenericEBSResponseFields) MaskPAN() {
	if res.PAN != "" {
		length := len(res.PAN)
		res.PAN = res.PAN[length-4 : length]
	}

	if res.ToCard != "" {
		length := len(res.ToCard)
		res.ToCard = res.ToCard[length-4 : length]
	}

	if res.FromCard != "" {
		length := len(res.FromCard)
		res.FromCard = res.FromCard[length-4 : length]
	}
}

type ConsumerCommonFields struct {
	ApplicationId string `json:"applicationId" form:"applicationId" binding:"required"`
	TranDateTime  string `json:"tranDateTime" form:"tranDateTime" binding:"required"`
	UUID          string `json:"UUID" form:"UUID" binding:"required"`
}

type ConsumerBillInquiryFields struct {
	ConsumerCommonFields
	ConsumersBillersFields
	ConsumerCardHolderFields
}

type ConsumerCardHolderFields struct {
	Pan     string `json:"PAN" form:"PAN" binding:"required"`
	Ipin    string `json:"IPIN" form:"IPIN" binding:"required"`
	ExpDate string `json:"expDate" form:"expDate" binding:"required"`
}

type ConsumerIsAliveFields struct {
	ConsumerCommonFields
}

type ConsumerBalanceFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
}
type ConsumersBillersFields struct {
	PayeeId     string `json:"payeeId" form:"payeeId" binding:"required"`
	PaymentInfo string `json:"paymentInfo" form:"paymentInfo" binding:"required"`
}

type ConsumerPurchaseFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	AmountFields
	ServiceProviderId string `json:"serviceProviderId" binding:"required"`
}

type ConsumerQRPaymentFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	AmountFields
	MerchantID string `json:"merchantID" binding:"required"`
}

type ConsumerQRRefundFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	OriginalTranUUID string `json:"originalTranUUID" binding:"required"`
}

type MerchantRegistrationFields struct {
	ConsumerCommonFields
	Merchant
	//allowed fields are CARD only for now. CF ebs document
	MerchantAccountType string `json:"merchantAccountType" binding:"required"`
	// this is the pan
	MerchantAccountReference string `json:"merchantAccountReference" binding:"required"`
	ExpDate                  string `json:"expDate" binding:"required"`
}

type Merchant struct {
	MerchantName         string `json:"merchantName" binding:"required"`
	MerchantCity         string `json:"merchantCity" binding:"required"`
	MerchantMobileNumber string `json:"mobileNo" binding:"required"`
	IDType               string `json:"idType" binding:"required"`
	IDNo                 string `json:"idNo" binding:"required"`
	TerminalID           string `json:"-"`
}
type mLabel struct {
	Value string
	Label string
	Help  string
}

func (m *Merchant) Details() []mLabel {
	res := []mLabel{
		{"merchantName", "Merchant Name", "Enter merchant name"},
		{"mobileNo", "Mobile Number", "Enter mobile number"},
		{"merchantCity", "Merchant City", "Enter merchant city"},
		{"idType", "ID Type (National ID, Driving License", "Enter ID Type"},
		{"idNo", "ID number", "Enter id no"},
	}
	return res
}

func (m *Merchant) ToMap() map[string]interface{} {
	a := make(map[string]interface{})
	a["merchantName"] = m.MerchantName
	a["mobileNo"] = m.MerchantMobileNumber
	a["merchantCity"] = m.MerchantCity
	a["idType"] = m.IDType
	a["idNo"] = m.IDNo
	return a
}

func (m *Merchant) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type ConsumerBillPaymentFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	AmountFields
	ConsumersBillersFields
}

type ConsumerWorkingKeyFields struct {
	ConsumerCommonFields
}

type ConsumerIPinFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	NewIPIN string `json:"newIPIN" binding:"required"`
}

type ConsumerCardTransferAndMobileFields struct {
	ConsumerCardTransferFields
	Mobile string `json:"mobile_number"`
}
type ConsumerCardTransferFields struct {
	ConsumerCommonFields
	ConsumerCardHolderFields
	AmountFields
	ToCard string `json:"toCard" binding:"required"`
}

type ConsumerStatusFields struct {
	ConsumerCommonFields
	OriginalTranUUID string `json:"originalTranUUID" binding:"required"`
}

type ConsumerGenerateIPin struct {
	ConsumerCommonFields
	Pan          string `json:"pan"`
	MobileNumber string `json:"mobileNumber"`
	Expdate      string `json:"expDate"`
}

type ConsumerGenerateIPinCompletion struct {
	ConsumerCommonFields
	Pan          string `json:"pan"`
	Expdate      string `json:"expDate"`
	MobileNumber string `json:"mobileNumber"`
	Otp          string `json:"otp"`
	Ipin         string `json:"ipin"`
}

type DisputeFields struct {
	Time    string  `json:"time,omitempty"`
	Service string  `json:"service,omitempty"`
	UUID    string  `json:"uuid,omitempty"`
	STAN    int     `json:"stan,omitempty"`
	Amount  float32 `json:"amount,omitempty"`
}

func (d *DisputeFields) New(f EBSParserFields) *DisputeFields {
	d.Amount = f.TranAmount
	d.Service = f.EBSServiceName
	d.UUID = f.UUID
	d.Time = f.TranDateTime
	return d
}

type CardsRedis struct {
	ID      int    `json:"id,omitempty"`
	PAN     string `json:"pan" binding:"required"`
	Expdate string `json:"exp_date" binding:"required,len=4"`
	IsMain  bool   `json:"is_main"`
	Name    string `json:"name"`
}

//
//func (c *CardsRedis) AddCard(username string) error {
//	buf, err := json.Marshal(c)
//
//	rC := utils.GetRedis()
//	if err != nil {
//		return err
//	}
//
//	z := &redis.Z{
//		Member: buf,
//	}
//	rC.ZAdd(username+":cards", z)
//	if c.IsMain {
//		rC.HSet(username, "main_card", buf)
//	}
//	return nil
//}

//func (c CardsRedis) RmCard(username string, id int) {
//	rC := utils.GetRedis()
//	if c.IsMain {
//		rC.HDel(username+":cards", "main_card")
//	} else {
//		rC.ZRemRangeByRank(username+":cards", int64(id-1), int64(id-1))
//	}
//}

type MobileRedis struct {
	Mobile   string `json:"mobile" binding:"required"`
	Provider string `json:"provider"`
	IsMain   bool   `json:"is_main"`
}

type ItemID struct {
	ID     int  `json:"id,omitempty" binding:"required"`
	IsMain bool `json:"is_main"`
}

func isEBS(pan string) bool {
	/*
		Bank Code        Bank Card PREFIX        Bank Short Name        Bank Full name
		2                    639186                      FISB                                 Faisal Islamic Bank
		4                    639256                      BAKH                                  Bank of Khartoum
		16                    639184                       RAKA                                  Al Baraka Sudanese Bank
		30                    639330                       ALSA                                  Al Salam Bank
	*/

	re := regexp.MustCompile(`(^639186|^639256|^639184|^639330)`)
	return re.Match([]byte(pan))
}
