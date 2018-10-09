package modifiedclaims

import (
    "errors"
    "sync"

    "github.com/rs/xid"
)

var (
    mList []ModifiedClaims
    mtx  sync.RWMutex
    once sync.Once
)

func init() {
    once.Do(initializeList)
}

func initializeList() {
	mList = []ModifiedClaims{}
}

// Modified Claims data structure
type ModifiedClaims struct {
    ID       				string `json:"id"`
	ClaimType 				string `json:"claimtype"`
	FromDate 				string `json:"fromDate"`
	ToDate 					string `json:"toDate"`
	PlaceOfService 			string `json:"placeOfService"`
	ProviderId				string `json:"providerId"`
	ProviderType			string `json:"providerType"`
	ProviderSpecialty 		string `json:"providerSpecialty"`
	ProcedureCode			string `json:"procedureCode"`
	DiagnosisCode			string `json:"diagnosisCode"`
	NetworkIndicator 		string `json:"networkIndicator"`
	SubscriberId			string `json:"subscriberId"`
	PatientAccountNumber 	string `json:"patientAccountNumber"`
	SccfNumber 				string `json:"sccfNumber"`
	RevenueCode 			string `json:"revenueCode"`
	BillType 				string `json:"billType"`
	Modifier 				string `json:"modifier"`
	PlanCode 				string `json:"planCode"`
	SfMessageCode 			string `json:"sfMessageCode"`
	PricingMethod 			string `json:"pricingMethod"`
	PricingRule 			string `json:"pricingRule"`
	DeliveryMethod 			string `json:"deliveryMethod"`
	InputDate 				string `json:"inputDate"`
	FileName 				string `json:"fileName"`
}

// Get retrieves all elements from the claims list
func GetModifiedClaims() []ModifiedClaims {
    return mList
}

// Add will add a new modified claim
func AddModifiedClaim(claimType string, fromDate string, toDate string, placeOfService string, providerId string, 
				providerType string, providerSpecialty string, procedureCode string, diagnosisCode string,
				networkIndicator string, subscriberId string, patientAccountNumber string, sccfNumber string,
				revenueCode string, billType string, modifier string, planCode string, sfMessageCode string,
				pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) string {
    t := newModifiedClaim(claimType, fromDate, toDate, placeOfService, providerId, 
				providerType, providerSpecialty, procedureCode, diagnosisCode,
				networkIndicator, subscriberId, patientAccountNumber, sccfNumber,
				revenueCode, billType, modifier, planCode, sfMessageCode,
				pricingMethod, pricingRule, deliveryMethod, inputDate, fileName)
    mtx.Lock()
    mList = append(mList, t)
    mtx.Unlock()
    return t.ID
}

// Delete will remove a claim from the Claims list
func Delete(id string) error {
    location, err := findClaimsLocation(id)
    if err != nil {
        return err
    }
    removeElementByLocation(location)
    return nil
}

func newModifiedClaim(claimType string, fromDate string, toDate string, placeOfService string, providerId string, 
				providerType string, providerSpecialty string, procedureCode string, diagnosisCode string,
				networkIndicator string, subscriberId string, patientAccountNumber string, sccfNumber string,
				revenueCode string, billType string, modifier string, planCode string, sfMessageCode string,
				pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) ModifiedClaims {
    return ModifiedClaims{
        ID:       xid.New().String(),
		ClaimType: claimType,
		FromDate: fromDate,
		ToDate: toDate,
		PlaceOfService: placeOfService,
		ProviderId: providerId,
		ProviderType: providerType,
		ProviderSpecialty: providerSpecialty,
		ProcedureCode: procedureCode,
		DiagnosisCode: diagnosisCode,
		NetworkIndicator: networkIndicator,
		SubscriberId: subscriberId,
		PatientAccountNumber: patientAccountNumber,
		SccfNumber: sccfNumber,
		RevenueCode: revenueCode,
		BillType: billType,
		Modifier: modifier,
		PlanCode: planCode,
		SfMessageCode: sfMessageCode,
		PricingMethod: pricingMethod,
		PricingRule: pricingRule,
		DeliveryMethod: deliveryMethod,
		InputDate: inputDate,
		FileName: fileName,
    }
}

func findClaimsLocation(id string) (int, error) {
    mtx.RLock()
    defer mtx.RUnlock()
    for i, t := range mList {
        if isMatchingID(t.ID, id) {
            return i, nil
        }
    }
    return 0, errors.New("could not find claims based on id")
}

func removeElementByLocation(i int) {
    mtx.Lock()
    mList = append(mList[:i], mList[i+1:]...)
    mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
    return a == b
}