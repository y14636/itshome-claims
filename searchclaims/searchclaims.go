package searchclaims

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"github.com/y14636/itshome-claims/model"
	"github.com/y14636/itshome-claims/utilities"
)

const SELECT_STATEMENT string = "SELECT TOP (?) orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, ''), orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName, orig.CREATE_DT, orig.CREATED_BY, COALESCE(pricing.SFMessageCode, '') AS PsfMessageCode, COALESCE(pricing.PricingMethod, ''), COALESCE(pricing.PricingRule, ''), COALESCE(orig_proc.ProcedureCode, ''), COALESCE(orig_proc.RevenueCode, ''), COALESCE(orig_proc.Modifier, ''), COALESCE(orig_proc.DateOfService, ''), COALESCE(orig_proc.DateOfServiceTo, ''), COALESCE(orig_proc.PlaceOfService, '') FROM ITSHome.OriginalClaims orig LEFT OUTER JOIN ITSHome.OriginalPricing pricing ON orig.Id = pricing.OriginalClaimID LEFT OUTER JOIN ITSHome.OriginalProcedures orig_proc ON pricing.OriginalClaimID = orig_proc.OriginalClaimID"

//const SELECT_STATEMENT_ALL string = "SELECT orig.Id, orig.ClaimType, COALESCE(orig.ServiceId, 'N/A') AS ServiceId, orig.ReceiptDate, orig.FromDate, orig.ToDate, orig.ProviderId, orig.ProviderType, orig.ProviderSpecialty, orig.DiagnosisCode, orig.NetworkIndicator, orig.SubscriberId, orig.PatientAccountNumber, orig.SCCFNumber, orig.BillType, orig.PlanCode, orig.SFMessageCode, orig.DeliveryMethod, orig.InputDate, orig.FileName FROM ITSHome.OriginalClaims orig"

var (
	logger               = logrus.New()
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
	file, err := os.OpenFile("searchclaims.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize log file %s", err)
		os.Exit(1)
	}

	logger = &logrus.Logger{
		Out:   file,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}
}

func initializeList() {
	list = []structs.Claims{}
}

func GetResults(search string) []structs.Claims {
	var rList = []structs.Claims{}
	criteria := utilities.ParseParameters(search)
	logger.Debug("Inside GetResults(), before CleanParameters()", criteria)
	criteria = GetClaimsThreshold(criteria)
	logger.Debug("Inside GetResults(), after GetClaimsThreshold()", criteria)
	newCriteria := strings.Split(criteria, "|")
	claimsThreshold := newCriteria[1]
	claimsThresholdValue := strings.SplitAfter(claimsThreshold, "=")
	thresholdLimit := utilities.TrimQuote(claimsThresholdValue[1])
	logger.Debug("thresholdLimit", thresholdLimit)
	criteria = utilities.CleanParameters(newCriteria[0])
	logger.Debug("Inside GetResults(), after CleanParameters()", criteria)

	condb, errdb := utilities.GetSqlConnection()
	if errdb != nil {
		logger.Fatal(" Error open db:", errdb.Error())
	}

	rList = FetchClaims(condb, criteria, thresholdLimit)

	defer condb.Close()

	return rList
}

func GetClaimsThreshold(criteria string) string {
	logger.Debug("Inside GetClaimsThreshold")
	newCriteria := criteria
	var paramName string
	var paramValue string
	result := strings.Split(newCriteria, ";")
	for j := range result {
		logger.Debug("result", result[j])
		field := result[j]
		fieldName := strings.Split(field, "=")
		logger.Debug("fieldName", fieldName[0])
		if len(fieldName[0]) > 0 && fieldName[0] == "claimsThreshold" {
			paramName = fieldName[0]
			logger.Debug("name", paramName)
			logger.Debug("value", fieldName[1])
			logger.Debug("string before removing", newCriteria)
			logger.Debug("string to be removed", paramName+"="+fieldName[1])
			paramValue = fieldName[1]
			logger.Debug("paramValue", paramValue)
			newCriteria = strings.Replace(newCriteria, paramName+"="+paramValue+";", "", -1)
			//fmt.Println("string after removing", criteria)
		}
	}
	return newCriteria + "|" + paramName + "=" + paramValue
}

func FetchClaims(condb *sql.DB, criteria string, thresholdLimit string) []structs.Claims {
	var rList = []structs.Claims{}
	strSelect := SELECT_STATEMENT
	if len(criteria) > 0 {
		strSelect += criteria
	}
	logger.Debug("select statement", strSelect)
	i, err := strconv.Atoi(thresholdLimit)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := condb.Query(strSelect, i)
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
