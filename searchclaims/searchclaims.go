package searchclaims

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"

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

func GetResults(search string) []Claims {
	var rList = []Claims{}
	criteria := utilities.ParseParameters(search)
	criteria = utilities.CleanParameters(criteria)

	condb, errdb := utilities.GetSqlConnection()
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
		rList = append(rList, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return rList
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
