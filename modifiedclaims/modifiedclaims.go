package modifiedclaims

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/rs/xid"
	"github.com/y14636/itshome-claims/utilities"
)

const READY_STATUS = "Ready"

var (
	mList []ModifiedClaims
	mtx   sync.RWMutex
	once  sync.Once
	condb *sql.DB
)

func init() {
	once.Do(initializeList)
}

func initializeList() {
	mList = []ModifiedClaims{}
}

// Modified Claims data structure
type ModifiedClaims struct {
	ID                   string `json:"id"`
	ClaimType            string `json:"claimtype"`
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

func AddMultipleClaims(claimsData string) error {
	fmt.Println("Add Multiple Claims", claimsData)

	criteria := utilities.ParseParameters(claimsData)
	fmt.Println("Add Multiple Claims", criteria)
	// Split on comma.
	result := strings.Split(criteria, "&")
	var errdb error
	condb, errdb = utilities.GetSqlConnection()
	if errdb != nil {
		log.Fatal("Open connection failed:", errdb.Error())
	}
	fmt.Printf("Connected!\n")
	defer condb.Close()

	fmt.Println("insert criteria", result[1])

	claimIds := strings.Split(result[1], ",")
	fmt.Println("length of claimIds", len(claimIds))
	if len(claimIds) == 1 {
		claimIds[0] = utilities.TrimSuffix(claimIds[0], ";")
	}
	currentTime := time.Now()
	currentDateTime := currentTime.Format("2006-01-02 15:04:05")
	modifiedBy := "mprue"
	//data := "15,16,17"
	//claimIds := strings.Split(data, ",")

	// Display all elements.
	for i := range claimIds {
		fmt.Println("Claim Ids", claimIds[i])
		// Create claims
		createID, err := CreateModifiedClaims(condb, claimIds[i], modifiedBy, currentDateTime)
		if err != nil {
			log.Fatal("Creating Modified Claims failed: ", err.Error())
		}
		fmt.Printf("Inserted ID: %d successfully.\n", createID)
		id := int(createID)
		t := strconv.Itoa(id)

		newCriteria := UpdateSubscriberId(result[0])
		updateID, err := UpdateModifiedClaims(condb, newCriteria, t)
		if err != nil {
			log.Fatal("Updating Modified Claims failed: ", err.Error())
		}
		fmt.Printf("Updated ID: %d successfully.\n", updateID)

		// Create modified procedures
		createProceduresID, err := CreateModifiedProcedures(condb, claimIds[i], modifiedBy, currentDateTime)
		if err != nil {
			log.Fatal("Creating Modified Procedures failed: ", err.Error())
		}
		fmt.Printf("Inserted ID: %d successfully.\n", createProceduresID)
		procId := int(createProceduresID)
		pid := strconv.Itoa(procId)

		updateProceduresID, err := UpdateModifiedProcedures(condb, newCriteria, pid)
		if err != nil {
			log.Fatal("Updating Modified Procedures failed: ", err.Error())
		}
		fmt.Printf("Updated Procedures ID: %d successfully.\n", updateProceduresID)
	}

	fmt.Println("update criteria", result[0])
	return nil
}

func UpdateSubscriberId(criteria string) string {
	strCriteria := strings.Split(criteria, ";")
	var newCriteria string
	var subIdHolder string
	var suffixHolder string
	for i := range strCriteria {
		fmt.Println("strCriteria", strCriteria[i])
		parameter := strings.Split(strCriteria[i], "=")
		fmt.Println(len(parameter[0]))
		if len(parameter[0]) > 0 {
			if len(parameter[1]) > 0 && parameter[0] == "subscriberId" {
				subIdHolder = trimQuote(parameter[1])
			} else if len(parameter[1]) > 0 && parameter[0] == "suffix" {
				suffixHolder = trimQuote(parameter[1])
			} else {
				newCriteria += parameter[0] + "=" + parameter[1] + ";"
				fmt.Println("new criteria", newCriteria)
			}
		}
	}

	newCriteria = newCriteria + "subscriberId='" + subIdHolder + suffixHolder + "';"
	fmt.Println("newCriteria end", newCriteria)
	return newCriteria
}

func trimQuote(s string) string {
	s = s[1 : len(s)-1]
	return s
}

func CreateModifiedClaims(db *sql.DB, claimId string, userId string, currentDateTime string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO ITSHome.ModifiedClaims ( OriginalClaimID, FromDate, ToDate, DiagnosisCode, NetworkIndicator, SubscriberId, PatientAccountNumber, SCCFNumber, Claim, ModifiedDate, ModifiedBy, Status, CreatedDate, CREATE_DT, CREATED_BY 	) SELECT Id AS OriginalClaimID, FromDate, ToDate, DiagnosisCode, NetworkIndicator, SubscriberId, PatientAccountNumber, SCCFNumber,	Claim, '%s' AS ModifiedDate, '%s' AS ModifiedBy, '%s' AS Status, '%s' AS CreatedDate, '%s' as CREATE_DT, '%s' as CREATED_BY FROM ITSHome.OriginalClaims WHERE Id = '%s'",
		currentDateTime, userId, READY_STATUS, currentDateTime, currentDateTime, currentDateTime, claimId)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Println("Error inserting new row: " + err.Error())
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateModifiedClaims(db *sql.DB, criteria string, modifiedClaimId string) (int64, error) {
	fmt.Println("Inside UpdateModifiedClaims", criteria)
	fmt.Println("Inside UpdateModifiedClaims", modifiedClaimId)
	updateCriteria := strings.Split(criteria, ";")
	fmt.Println(len(updateCriteria))
	for i := range updateCriteria {
		fmt.Println("update criteria", updateCriteria[i])
		parameter := strings.Split(updateCriteria[i], "=")
		fmt.Println(len(parameter[0]))
		if len(parameter[0]) > 0 && parameter[0] != "modifier" && parameter[0] != "procedureCode" {
			tsql := fmt.Sprintf("UPDATE ITSHome.ModifiedClaims SET %s = %s WHERE Id = '%s'",
				strings.Title(parameter[0]), parameter[1], modifiedClaimId)
			result, err := db.Exec(tsql)
			if err != nil {
				fmt.Println("Error updating row: " + err.Error())
				return -1, err
			}
			fmt.Printf("Updated ID: %d successfully.\n", result)
		}
	}
	return 0, nil
}

func CreateModifiedProcedures(db *sql.DB, claimId string, userId string, currentDateTime string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO ITSHome.ModifiedProcedures ( ModifiedClaimID, LineIndex, ProcedureCode, RevenueCode, Modifier, DateOfService, DateOfServiceTo, CREATE_DT, CREATED_BY ) SELECT OriginalClaimID, LineIndex, ProcedureCode, RevenueCode, Modifier, DateOfService, DateOfServiceTo, '%s' AS CREATE_DT, '%s' AS CREATED_BY FROM ITSHome.OriginalProcedures WHERE OriginalClaimID = '%s'",
		currentDateTime, userId, claimId)
	result, err := db.Exec(tsql)
	if err != nil {
		fmt.Println("Error inserting new row: " + err.Error())
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateModifiedProcedures(db *sql.DB, criteria string, modifiedClaimId string) (int64, error) {
	fmt.Println("Inside UpdateModifiedProcedures", criteria)
	fmt.Println("Inside UpdateModifiedProcedures", modifiedClaimId)
	updateCriteria := strings.Split(criteria, ";")
	fmt.Println(len(updateCriteria))
	for i := range updateCriteria {
		fmt.Println("update criteria", updateCriteria[i])
		parameter := strings.Split(updateCriteria[i], "=")
		fmt.Println(len(parameter[0]))
		if len(parameter[0]) > 0 && (parameter[0] == "modifier" || parameter[0] == "procedureCode") {
			tsql := fmt.Sprintf("UPDATE ITSHome.ModifiedProcedures SET %s = %s WHERE Id = '%s'",
				strings.Title(parameter[0]), parameter[1], modifiedClaimId)
			result, err := db.Exec(tsql)
			if err != nil {
				fmt.Println("Error updating row: " + err.Error())
				return -1, err
			}
			fmt.Printf("Updated ID: %d successfully.\n", result)
		}
	}
	// return result.LastInsertId()
	return 0, nil
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
		ID:                   xid.New().String(),
		ClaimType:            claimType,
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
