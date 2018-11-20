package claims

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/rs/xid"
	"github.com/y14636/itshome-claims/model"
	"github.com/y14636/itshome-claims/utilities"
)

const SELECT_STATEMENT string = "SELECT orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, ''), orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName, orig.CREATE_DT, orig.CREATED_BY, COALESCE(pricing.SFMessageCode, '') AS PsfMessageCode, COALESCE(pricing.PricingMethod, ''), COALESCE(pricing.PricingRule, ''), COALESCE(orig_proc.ProcedureCode, ''), COALESCE(orig_proc.RevenueCode, ''), COALESCE(orig_proc.Modifier, ''), COALESCE(orig_proc.DateOfService, ''), COALESCE(orig_proc.DateOfServiceTo, ''), COALESCE(orig_proc.PlaceOfService, '') FROM ITSHome.OriginalClaims orig LEFT OUTER JOIN ITSHome.OriginalPricing pricing ON orig.Id = pricing.OriginalClaimID LEFT OUTER JOIN ITSHome.OriginalProcedures orig_proc ON pricing.OriginalClaimID = orig_proc.OriginalClaimID"
const SELECT_STATEMENT_ALL string = "SELECT orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, 'N/A') AS ServiceId, orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName FROM ITSHome.OriginalClaims orig"

var (
	list                 = []structs.Claims{}
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
	list = []structs.Claims{}
	//dummy Institutional claims
	// Add("11", "N/A", "20180110", "20180101", "20180110", "NA", "000000001000", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098001", "99999999", "30120180400000000",
	// 	"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile1.dat")
}

// Get retrieves all elements from the structs.Claims list
func Get() []structs.Claims {
	var list []structs.Claims

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

func GetListById(claimId string) []structs.Claims {
	fmt.Println("inside GetListById", claimId)
	var list []structs.Claims

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

// Delete will remove a claim from the structs.Claims list
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
	modifier string, dosFrom string, dosTo string, placeOfService string) structs.Claims {
	return structs.Claims{
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
	pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) structs.Claims {
	return structs.Claims{
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
