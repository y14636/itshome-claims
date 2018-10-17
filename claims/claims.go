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
)

var (
	list                      []Claims
	qList                     []Results
	mtx                       sync.RWMutex
	once                      sync.Once
	USER_SECURITY_QUESTION_ID int
	QUESTION                  string
	ANSWER                    string
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	list = []Claims{}
	qList = []Results{}

	condb, errdb := sql.Open("mssql", "server=SQLMODL29\\SQL_MODL29;user id=;password=;database=idb03q_qual")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}
	rows, err := condb.Query("SELECT USER_SECURITY_QUESTION_ID, QUESTION, ANSWER FROM dbo.USER_SECURITY_QUESTION WHERE REG_USER_ID = 'BR0000010402'	AND USER_SECURITY_QUESTION_ID = ?", 4123)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&USER_SECURITY_QUESTION_ID, &QUESTION, &ANSWER)
		if err != nil {
			log.Fatal(err)
		}
		t := strconv.Itoa(USER_SECURITY_QUESTION_ID)
		AddToQlist(t, QUESTION, ANSWER)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	defer condb.Close()

	//dummy Institutional claims
	Add("11", "N/A", "20180110", "20180101", "20180110", "NA", "000000001000", "0001", "314000000X", "52260", "N3010", "NA", "PAY20089098001", "99999999", "30120180400000000",
		"0450", "111", "59", "302", "P302", "30", "012", "2", "20180110", "testfile1.dat")
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

func GetResults() []Results {
	return qList
}

// Get retrieves all elements from the claims list
func Get() []Claims {
	return list
}

func AddToQlist(id string, question string, answer string) string {
	t := newResult(id, question, answer)
	mtx.Lock()
	qList = append(qList, t)
	mtx.Unlock()
	return t.ID
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

func newResult(id string, question string, answer string) Results {
	return Results{
		ID:       id,
		Question: question,
		Answer:   answer,
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
		if isMatchingID(t.ID, id) {
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
