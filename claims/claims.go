package claims

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/rs/xid"
	"github.com/y14636/itshome-claims/utilities"
)

const SELECT_STATEMENT string = "SELECT orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, ''), orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName, orig.CREATE_DT, orig.CREATED_BY, COALESCE(pricing.SFMessageCode, '') AS PsfMessageCode, COALESCE(pricing.PricingMethod, ''), COALESCE(pricing.PricingRule, ''), COALESCE(orig_proc.ProcedureCode, ''), COALESCE(orig_proc.RevenueCode, ''), COALESCE(orig_proc.Modifier, ''), COALESCE(orig_proc.DateOfService, ''), COALESCE(orig_proc.DateOfServiceTo, ''), COALESCE(orig_proc.PlaceOfService, '') FROM ITSHome.OriginalClaims orig LEFT OUTER JOIN ITSHome.OriginalPricing pricing ON orig.Id = pricing.OriginalClaimID LEFT OUTER JOIN ITSHome.OriginalProcedures orig_proc ON pricing.OriginalClaimID = orig_proc.OriginalClaimID"
const SELECT_STATEMENT_ALL string = "SELECT orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, 'N/A') AS ServiceId, orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName FROM ITSHome.OriginalClaims orig"

var (
	list                 = []Claims{}
	mtx                  sync.RWMutex
	once                 sync.Once
	id                   int
	claimType            string
	serviceId            string
	receiptDate          string
	fromDate             string
	toDate               string
	providerId           string
	providerType         string
	providerSpecialty    string
	diagnosisCode        string
	networkIndicator     string
	subscriberId         string
	patientAccountNumber string
	sccfNumber           string
	billType             string
	planCode             string
	sfMessageCode        string
	deliveryMethod       string
	inputDate            string
	fileName             string
	createDate           string
	createdBy            string
	pSfMessageCode       string
	pricingMethod        string
	pricingRule          string
	procedureCode        string
	revenueCode          string
	modifier             string
	dosFrom              string
	dosTo                string
	placeOfService       string
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Claims{}

	//dummy Institutional claims
	// Add("11", "N/A", "20180110", "20180101", "20180110", "NA", "000000001000", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098001", "99999999", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile1.dat")
	// Add("12", "N/A", "20180710", "20180701", "20180710", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "22222222", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20180210", "20180201", "20180210", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "33333333", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20180310", "20180301", "20180310", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "44444444", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20180410", "20180401", "20180410", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "55555555", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("12", "N/A", "20180510", "20180501", "20180510", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "66666666", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("12", "N/A", "20180610", "20180601", "20180610", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "77777777", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20180810", "20180801", "20180810", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "88888888", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20180910", "20180901", "20180910", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "99999991", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("11", "N/A", "20181010", "20181001", "20181010", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "99999992", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// Add("12", "N/A", "20180107", "20180111", "20180112", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098004", "99999993", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	// //dummy Professional claims
	// Add("20", "600", "20180110", "20180105", "20180110", "30", "000000001001", "A4", "111N00000X  ", "99212", "M5134", "3", "PAY20089098001", "11111111", "30120180400000001",
	// 	"NA", "111", "50", "302", "P302", "40", "9", "A", "20180110", "testfile2.dat")
}

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

// Get retrieves all elements from the claims list
func Get() []Claims {
	var list []Claims

	condb, errdb := utilities.GetSqlConnection()
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}

	rows, err := condb.Query(SELECT_STATEMENT)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &claimType, &serviceId, &receiptDate, &fromDate, &toDate, &providerId,
			&providerType, &providerSpecialty, &diagnosisCode, &networkIndicator, &subscriberId,
			&patientAccountNumber, &sccfNumber, &billType, &planCode, &sfMessageCode, &deliveryMethod,
			&inputDate, &fileName, &createDate, &createdBy, &pSfMessageCode, &pricingMethod, &pricingRule,
			&procedureCode, &revenueCode, &modifier, &dosFrom, &dosTo, &placeOfService)
		if err != nil {
			log.Fatal(err)
		}
		strId := strconv.Itoa(id)
		t := newResult(strId, claimType, serviceId, receiptDate, fromDate, toDate, providerId,
			providerType, providerSpecialty, diagnosisCode, networkIndicator, subscriberId,
			patientAccountNumber, sccfNumber, billType, planCode, sfMessageCode, deliveryMethod,
			inputDate, fileName, createDate, createdBy, pSfMessageCode, pricingMethod, pricingRule,
			procedureCode, revenueCode, modifier, dosFrom, dosTo, placeOfService)
		list = append(list, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return list
}

func GetListById(claimId string) []Claims {
	fmt.Println("inside GetListById", claimId)
	var list []Claims

	condb, errdb := utilities.GetSqlConnection()
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}

	rows, err := condb.Query(SELECT_STATEMENT + " WHERE orig.Id=" + claimId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &claimType, &serviceId, &receiptDate, &fromDate, &toDate, &providerId,
			&providerType, &providerSpecialty, &diagnosisCode, &networkIndicator, &subscriberId,
			&patientAccountNumber, &sccfNumber, &billType, &planCode, &sfMessageCode, &deliveryMethod,
			&inputDate, &fileName, &createDate, &createdBy, &pSfMessageCode, &pricingMethod, &pricingRule,
			&procedureCode, &revenueCode, &modifier, &dosFrom, &dosTo, &placeOfService)
		if err != nil {
			log.Fatal(err)
		}
		strId := strconv.Itoa(id)
		t := newResult(strId, claimType, serviceId, receiptDate, fromDate, toDate, providerId,
			providerType, providerSpecialty, diagnosisCode, networkIndicator, subscriberId,
			patientAccountNumber, sccfNumber, billType, planCode, sfMessageCode, deliveryMethod,
			inputDate, fileName, createDate, createdBy, pSfMessageCode, pricingMethod, pricingRule,
			procedureCode, revenueCode, modifier, dosFrom, dosTo, placeOfService)
		list = append(list, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return list

}

// Add will add a new claim
func Add(claimType string, serviceId string, receiptDate string, fromDate string, toDate string, placeOfService string, providerId string,
	providerType string, providerSpecialty string, procedureCode string, diagnosisCode string,
	networkIndicator string, subscriberId string, patientAccountNumber string, sccfNumber string,
	revenueCode string, billType string, modifier string, planCode string, sfMessageCode string,
	pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) string {
	t := newClaim(claimType, serviceId, receiptDate, fromDate, toDate, placeOfService, providerId,
		providerType, providerSpecialty, procedureCode, diagnosisCode,
		networkIndicator, subscriberId, patientAccountNumber, sccfNumber,
		revenueCode, billType, modifier, planCode, sfMessageCode,
		pricingMethod, pricingRule, deliveryMethod, inputDate, fileName)
	mtx.Lock()
	list = append(list, t)
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

func newResult(id string, claimType string, serviceId string, receiptDate string, fromDate string, toDate string, providerId string,
	providerType string, providerSpecialty string, diagnosisCode string, networkIndicator string, subscriberId string, patientAccountNumber string,
	sccfNumber string, billType string, planCode string, sfMessageCode string, deliveryMethod string, inputDate string, fileName string,
	createDate string, createdBy string, pSfMessageCode string, pricingMethod string, pricingRule string, procedureCode string, revenueCode string,
	modifier string, dosFrom string, dosTo string, placeOfService string) Claims {
	return Claims{
		ID:                   id,
		ClaimType:            claimType,
		ServiceId:            serviceId,
		ReceiptDate:          receiptDate,
		FromDate:             fromDate,
		ToDate:               toDate,
		ProviderId:           providerId,
		ProviderType:         providerType,
		ProviderSpecialty:    providerSpecialty,
		DiagnosisCode:        diagnosisCode,
		NetworkIndicator:     networkIndicator,
		SubscriberId:         subscriberId,
		PatientAccountNumber: patientAccountNumber,
		SccfNumber:           sccfNumber,
		BillType:             billType,
		PlanCode:             planCode,
		SfMessageCode:        sfMessageCode,
		DeliveryMethod:       deliveryMethod,
		InputDate:            inputDate,
		FileName:             fileName,
		CreateDate:           createDate,
		CreatedBy:            createdBy,
		PsfMessageCode:       pSfMessageCode,
		PricingMethod:        pricingMethod,
		PricingRule:          pricingRule,
		ProcedureCode:        procedureCode,
		RevenueCode:          revenueCode,
		Modifier:             modifier,
		DosFrom:              dosFrom,
		DosTo:                dosTo,
		PlaceOfService:       placeOfService,
	}
}

func newClaim(claimType string, serviceId string, receiptDate string, fromDate string, toDate string, placeOfService string, providerId string,
	providerType string, providerSpecialty string, procedureCode string, diagnosisCode string,
	networkIndicator string, subscriberId string, patientAccountNumber string, sccfNumber string,
	revenueCode string, billType string, modifier string, planCode string, sfMessageCode string,
	pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) Claims {
	return Claims{
		ID:                   xid.New().String(),
		ClaimType:            claimType,
		ServiceId:            serviceId,
		ReceiptDate:          receiptDate,
		FromDate:             fromDate,
		ToDate:               toDate,
		PlaceOfService:       placeOfService,
		ProviderId:           providerId,
		ProviderType:         providerType,
		ProviderSpecialty:    providerSpecialty,
		ProcedureCode:        procedureCode,
		DiagnosisCode:        diagnosisCode,
		NetworkIndicator:     networkIndicator,
		SubscriberId:         subscriberId,
		PatientAccountNumber: patientAccountNumber,
		SccfNumber:           sccfNumber,
		RevenueCode:          revenueCode,
		BillType:             billType,
		Modifier:             modifier,
		PlanCode:             planCode,
		SfMessageCode:        sfMessageCode,
		PricingMethod:        pricingMethod,
		PricingRule:          pricingRule,
		DeliveryMethod:       deliveryMethod,
		InputDate:            inputDate,
		FileName:             fileName,
	}
}

func findClaimsLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		// if isMatchingID(t.ID, id) {
		// 	return i, nil
		// }
		fmt.Println("Subscriber ID=", t.SubscriberId)
		if isMatchingID(t.SubscriberId, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find claims based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
