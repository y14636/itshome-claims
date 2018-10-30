package claims

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/rs/xid"
	"github.com/y14636/itshome-claims/utilities"
)

const SELECT_STATEMENT string = "SELECT orig.Id, orig.ClaimType, orig.ServiceId, orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.Claim, orig.InputDate, orig.FileName, orig.CREATE_DT, orig.CREATED_BY, pricing.SFMessageCode, pricing.PricingMethod, pricing.PricingRule, orig_proc.ProcedureCode, orig_proc.RevenueCode, orig_proc.Modifier, orig_proc.DateOfService, orig_proc.DateOfServiceTo, orig_proc.PlaceOfService FROM ITSHome.OriginalClaims orig, ITSHome.OriginalPricing pricing, ITSHome.OriginalProcedures orig_proc WHERE orig.Id = pricing.OriginalClaimID AND pricing.OriginalClaimID = orig_proc.OriginalClaimID"

var (
	list                 []Claims
	mtx                  sync.RWMutex
	once                 sync.Once
	id                   int
	claimType            string
	serviceId            string
	receiptDate          string
	fromDate             string
	toDate               string
	placeOfService       string
	providerId           string
	providerType         string
	providerSpecialty    string
	procedureCode        string
	diagnosisCode        string
	networkIndicator     string
	subscriberId         string
	patientAccountNumber string
	sccfNumber           string
	revenueCode          string
	billType             string
	modifier             string
	planCode             string
	sfMessageCode        string
	pricingMethod        string
	pricingRule          string
	deliveryMethod       string
	inputDate            string
	fileName             string
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Claims{}

	//dummy Institutional claims
	Add("11", "N/A", "20180110", "20180101", "20180110", "NA", "000000001000", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098001", "99999999", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile1.dat")
	Add("12", "N/A", "20180710", "20180701", "20180710", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "22222222", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20180210", "20180201", "20180210", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "33333333", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20180310", "20180301", "20180310", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "44444444", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20180410", "20180401", "20180410", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "55555555", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("12", "N/A", "20180510", "20180501", "20180510", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "66666666", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("12", "N/A", "20180610", "20180601", "20180610", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "77777777", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20180810", "20180801", "20180810", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "88888888", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20180910", "20180901", "20180910", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "99999991", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("11", "N/A", "20181010", "20181001", "20181010", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098002", "99999992", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	Add("12", "N/A", "20180107", "20180111", "20180112", "NA", "000000001001", "0001", "314000000X", "52260", "N3010", "NA", "INT20089098004", "99999993", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile3.dat")
	//dummy Professional claims
	Add("20", "600", "20180110", "20180105", "20180110", "30", "000000001001", "A4", "111N00000X  ", "99212", "M5134", "3", "PAY20089098001", "11111111", "30120180400000001",
		"NA", "111", "50", "302", "P302", "40", "9", "A", "20180110", "testfile2.dat")
}

// Claims data structure
type Claims struct {
	ID                   string `json:"id"`
	ClaimType            string `json:"claimtype"`
	ServiceId            string `json:"serviceId"`
	ReceiptDate          string `json:"receiptDate"`
	FromDate             string `json:"fromDate"`
	ToDate               string `json:"toDate"`
	PlaceOfService       string `json:"placeOfService"`
	ProviderId           string `json:"providerId"`
	ProviderType         string `json:"providerType"`
	ProviderSpecialty    string `json:"providerSpecialty"`
	ProcedureCode        string `json:"procedureCode"`
	DiagnosisCode        string `json:"diagnosisCode"`
	NetworkIndicator     string `json:"networkIndicator"`
	SubscriberId         string `json:"subscriberId"`
	PatientAccountNumber string `json:"patientAccountNumber"`
	SccfNumber           string `json:"sccfNumber"`
	RevenueCode          string `json:"revenueCode"`
	BillType             string `json:"billType"`
	Modifier             string `json:"modifier"`
	PlanCode             string `json:"planCode"`
	SfMessageCode        string `json:"sfMessageCode"`
	PricingMethod        string `json:"pricingMethod"`
	PricingRule          string `json:"pricingRule"`
	DeliveryMethod       string `json:"deliveryMethod"`
	InputDate            string `json:"inputDate"`
	FileName             string `json:"fileName"`
}

type Results struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func GetResults(search string) []Claims {
	var rList = []Claims{}
	//fmt.Println("search string=", search)
	criteria := utilities.ParseParameters(search)

	condb, errdb := sql.Open("mssql", "server=SQLDEV34\\SQL_DEV34;user id=;password=;database=zdb63q_itshc_syst")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}

	rList = FetchClaims(condb, criteria)

	defer condb.Close()

	return rList
}

func FetchClaims(condb *sql.DB, criteria string) []Claims {
	var rList = []Claims{}
	strSelect := SELECT_STATEMENT
	if len(criteria) > 0 {
		strSelect += criteria
	}
	fmt.Println("select statement", strSelect)
	rows, err := condb.Query(strSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &claimType, &serviceId, &receiptDate, &fromDate, &toDate, &placeOfService, &providerId, &providerType, &providerSpecialty, &procedureCode, &diagnosisCode, &networkIndicator, &subscriberId, &patientAccountNumber, &sccfNumber, &revenueCode, &billType, &modifier, &planCode, &sfMessageCode, &pricingMethod, &pricingRule, &deliveryMethod, &inputDate, &fileName)
		if err != nil {
			log.Fatal(err)
		}
		strId := strconv.Itoa(id)
		t := newResult(strId, claimType, serviceId, receiptDate, fromDate, toDate, placeOfService, providerId, providerType, providerSpecialty, procedureCode, diagnosisCode, networkIndicator, subscriberId, patientAccountNumber, sccfNumber, revenueCode, billType, modifier, planCode, sfMessageCode, pricingMethod, pricingRule, deliveryMethod, inputDate, fileName)
		rList = append(rList, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return rList
}

// Get retrieves all elements from the claims list
func Get() []Claims {
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

func newResult(id string, claimType string, serviceId string, receiptDate string, fromDate string, toDate string, placeOfService string, providerId string,
	providerType string, providerSpecialty string, procedureCode string, diagnosisCode string,
	networkIndicator string, subscriberId string, patientAccountNumber string, sccfNumber string,
	revenueCode string, billType string, modifier string, planCode string, sfMessageCode string,
	pricingMethod string, pricingRule string, deliveryMethod string, inputDate string, fileName string) Claims {
	return Claims{
		ID:                   id,
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
