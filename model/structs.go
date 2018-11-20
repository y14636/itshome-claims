package structs

// Claims data structure
type Claims struct {
	ID                   string `json:"id"`
	ClaimType            string `json:"claimtype"`
	ServiceId            string `json:"serviceId"`
	ReceiptDate          string `json:"receiptDate"`
	FromDate             string `json:"fromDate"`
	ToDate               string `json:"toDate"`
	ProviderId           string `json:"providerId"`
	ProviderType         string `json:"providerType"`
	ProviderSpecialty    string `json:"providerSpecialty"`
	DiagnosisCode        string `json:"diagnosisCode"`
	NetworkIndicator     string `json:"networkIndicator"`
	SubscriberId         string `json:"subscriberId"`
	PatientAccountNumber string `json:"patientAccountNumber"`
	SccfNumber           string `json:"sccfNumber"`
	BillType             string `json:"billType"`
	PlanCode             string `json:"planCode"`
	SfMessageCode        string `json:"sfMessageCode"`
	DeliveryMethod       string `json:"deliveryMethod"`
	InputDate            string `json:"inputDate"`
	FileName             string `json:"fileName"`
	CreateDate           string `json:"createDate"`
	CreatedBy            string `json:"createdBy"`
	PsfMessageCode       string `json:"pSfMessageCode"`
	PricingMethod        string `json:"pricingMethod"`
	PricingRule          string `json:"pricingRule"`
	ProcedureCode        string `json:"procedureCode"`
	RevenueCode          string `json:"revenueCode"`
	Modifier             string `json:"modifier"`
	DosFrom              string `json:"dosFrom"`
	DosTo                string `json:"dosTo"`
	PlaceOfService       string `json:"placeOfService"`
}
