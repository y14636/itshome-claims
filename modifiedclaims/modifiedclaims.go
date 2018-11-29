package modifiedclaims

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	log "github.com/sirupsen/logrus"
	"github.com/y14636/itshome-claims/utilities"
)

const READY_STATUS = "Ready"
const SELECT_STATEMENT = "SELECT mc.Id, mc.SubscriberId, mc.OriginalClaimID, mc.SCCFNumber, COALESCE(mp.ProcedureCode, ''), mc.DiagnosisCode, COALESCE(mp.Modifier, ''), mc.PatientAccountNumber, COALESCE(mc.NetworkIndicator, 'N/A'), COALESCE(mc.FromDate, ''), COALESCE(mc.ToDate, ''), mc.Status, COALESCE(mp.DateOfService, ''), COALESCE(mp.DateOfServiceTo, ''), mc.CREATE_DT, COALESCE(mc.CREATED_BY, '') FROM ITSHome.ModifiedClaims mc LEFT OUTER JOIN ITSHome.ModifiedProcedures mp ON mc.Id = mp.ModifiedClaimID WHERE mc.Status = 'Ready' ORDER BY CREATE_DT DESC"

var (
	mtx   sync.RWMutex
	once  sync.Once
	condb *sql.DB

	id                   int
	subscriberId         string
	originalClaimID      int
	sccfNumber           string
	procedureCode        string
	diagnosisCode        string
	modifier             string
	patientAccountNumber string
	networkIndicator     string
	fromDate             string
	toDate               string
	status               string
	dosFrom              string
	dosTo                string
	createDate           string
	createdBy            string
)

// Modified Claims data structure
type ModifiedClaims struct {
	ID                   string `json:"id"`
	SubscriberId         string `json:"subscriberId"`
	OriginalClaimID      string `json:"originalClaimID"`
	SccfNumber           string `json:"sccfNumber"`
	ProcedureCode        string `json:"procedureCode"`
	DiagnosisCode        string `json:"diagnosisCode"`
	Modifier             string `json:"modifier"`
	PatientAccountNumber string `json:"patientAccountNumber"`
	NetworkIndicator     string `json:"networkIndicator"`
	FromDate             string `json:"fromDate"`
	ToDate               string `json:"toDate"`
	Status               string `json:"status"`
	DosFrom              string `json:"dosFrom"`
	DosTo                string `json:"dosTo"`
	CreateDate           string `json:"createDate"`
	CreatedBy            string `json:"createdBy"`
}

// Get retrieves all elements from the claims list
func GetModifiedClaims() []ModifiedClaims {
	var mList []ModifiedClaims

	var errdb error
	condb, errdb = utilities.GetSqlConnection()
	if errdb != nil {
		log.Fatal("Open connection failed:", errdb.Error())
	}
	log.Println("Connected!")
	defer condb.Close()

	rows, err := condb.Query(SELECT_STATEMENT)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &subscriberId, &originalClaimID, &sccfNumber, &procedureCode, &diagnosisCode, &modifier, &patientAccountNumber, &networkIndicator, &fromDate, &toDate, &status, &dosFrom, &dosTo, &createDate, &createdBy)
		if err != nil {
			log.Fatal(err)
		}
		strClaimId := strconv.Itoa(originalClaimID)
		strId := strconv.Itoa(id)
		t := newResultAll(strId, subscriberId, strClaimId, sccfNumber, procedureCode, diagnosisCode, modifier, patientAccountNumber, networkIndicator, fromDate, toDate, status, dosFrom, dosTo, createDate, createdBy)
		mList = append(mList, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return mList
}

func AddClaims(claimsData string) error {
	log.Println("Add Claims", claimsData)

	criteria := utilities.ParseParameters(claimsData)
	log.Println("Add Claims", criteria)

	result := strings.Split(criteria, "&")
	var errdb error
	condb, errdb = utilities.GetSqlConnection()
	if errdb != nil {
		log.Fatal("Open connection failed:", errdb.Error())
	}
	log.Println("Connected!")
	defer condb.Close()

	log.Println("insert criteria", result[1])

	claimIds := strings.Split(result[1], ",")
	log.Println("length of claimIds", len(claimIds))

	currentTime := time.Now()
	currentDateTime := currentTime.Format("2006-01-02 15:04:05")
	modifiedBy := "mprue"
	for i := range claimIds {
		claimIds[i] = utilities.TrimSuffix(claimIds[i], ";")
		log.Println("Claim Ids", claimIds[i])
		// Create claims
		createID, err := CreateModifiedClaims(condb, claimIds[i], modifiedBy, currentDateTime)
		if err != nil {
			log.Fatal("Creating Modified Claims failed: ", err.Error())
		}
		log.Printf("Inserted ID: %d successfully.", createID)
		id := int(createID)
		t := strconv.Itoa(id)

		newCriteria := UpdateSubscriberId(result[0], id)
		updateID, err := UpdateModifiedClaims(condb, newCriteria, t)
		if err != nil {
			log.Fatal("Updating Modified Claims failed: ", err.Error())
		}
		log.Printf("Updated ID: %d successfully.", updateID)

		count, err := CheckIfOriginalProceduresExist(updateID)
		if err != nil {
			log.Fatal("ReadEmployees failed:", err.Error())
		}

		if count > 0 {
			// Create modified procedures
			createProceduresID, err := CreateModifiedProcedures(condb, claimIds[i], modifiedBy, currentDateTime)
			if err != nil {
				log.Fatal("Creating Modified Procedures failed: ", err.Error())
			}
			log.Printf("Inserted createProceduresID: %d successfully.", createProceduresID)
			procId := int(createProceduresID)
			pid := strconv.Itoa(procId)

			updateProceduresID, err := UpdateModifiedProcedures(condb, newCriteria, pid)
			if err != nil {
				log.Fatal("Updating Modified Procedures failed: ", err.Error())
			}
			log.Printf("Updated Procedures ID: %d successfully.", updateProceduresID)
		}
	}

	log.Println("update criteria", result[0])
	return nil
}

func UpdateSubscriberId(criteria string, id int) string {
	strCriteria := strings.Split(criteria, ";")
	var newCriteria string
	var subIdHolder string
	var suffixHolder string
	for i := range strCriteria {
		log.Println("strCriteria", strCriteria[i])
		parameter := strings.Split(strCriteria[i], "=")
		log.Println(len(parameter[0]))
		if len(parameter[0]) > 0 {
			if len(parameter[1]) > 0 && parameter[0] == "subscriberId" {
				subIdHolder = utilities.TrimQuote(parameter[1])
			} else if len(parameter[1]) > 0 && parameter[0] == "suffix" {
				suffixHolder = utilities.TrimQuote(parameter[1])
			} else {
				newCriteria += parameter[0] + "=" + parameter[1] + ";"
				log.Println("new criteria", newCriteria)
			}
		}
	}

	if len(subIdHolder) == 0 && len(suffixHolder) > 0 {
		subId := strings.TrimSpace(GetSubscriberId(id))
		log.Println("Inside UpdateSubscriberId - subscriberId = ", subId)
		if utilities.IsNumeric(subId) && len(subId) > 9 {
			subIdHolder = utilities.TrimSuffix(subId, subId[len(subId)-2:])
		} else if !utilities.IsNumeric(subId) && len(subId) > 12 {
			subIdHolder = utilities.TrimSuffix(subId, subId[len(subId)-2:])
		} else {
			subIdHolder = subId
		}
	}

	newCriteria = newCriteria + "subscriberId='" + subIdHolder + suffixHolder + "';"
	log.Println("newCriteria end", newCriteria)
	return newCriteria
}

func CreateModifiedClaims(db *sql.DB, claimId string, userId string, currentDateTime string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO ITSHome.ModifiedClaims ( OriginalClaimID, FromDate, ToDate, DiagnosisCode, NetworkIndicator, SubscriberId, PatientAccountNumber, SCCFNumber, Claim, ModifiedDate, ModifiedBy, Status, CreatedDate, CREATE_DT, CREATED_BY 	) SELECT Id AS OriginalClaimID, FromDate, ToDate, DiagnosisCode, NetworkIndicator, SubscriberId, PatientAccountNumber, SCCFNumber,	Claim, '%s' AS ModifiedDate, '%s' AS ModifiedBy, '%s' AS Status, '%s' AS CreatedDate, '%s' as CREATE_DT, '%s' as CREATED_BY FROM ITSHome.OriginalClaims WHERE Id = '%s'",
		currentDateTime, userId, READY_STATUS, currentDateTime, currentDateTime, userId, claimId)
	result, err := db.Exec(tsql)
	if err != nil {
		log.Println("Error inserting new row: " + err.Error())
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateModifiedClaims(db *sql.DB, criteria string, modifiedClaimId string) (int64, error) {
	log.Println("Inside UpdateModifiedClaims", criteria)
	log.Println("Inside UpdateModifiedClaims", modifiedClaimId)
	updateCriteria := strings.Split(criteria, ";")
	log.Println(len(updateCriteria))
	for i := range updateCriteria {
		log.Println("update criteria", updateCriteria[i])
		parameter := strings.Split(updateCriteria[i], "=")
		log.Println(len(parameter[0]))
		if len(parameter[0]) > 0 && parameter[0] != "modifier" && parameter[0] != "procedureCode" {
			tsql := fmt.Sprintf("UPDATE ITSHome.ModifiedClaims SET %s = %s WHERE Id = '%s'",
				strings.Title(parameter[0]), parameter[1], modifiedClaimId)
			result, err := db.Exec(tsql)
			if err != nil {
				log.Println("Error updating row: " + err.Error())
				return -1, err
			}
			log.Printf("Updated ID: %d successfully.", result)
		}
	}
	return 0, nil
}

func CreateModifiedProcedures(db *sql.DB, claimId string, userId string, currentDateTime string) (int64, error) {
	tsql := fmt.Sprintf("INSERT INTO ITSHome.ModifiedProcedures ( ModifiedClaimID, LineIndex, ProcedureCode, RevenueCode, Modifier, DateOfService, DateOfServiceTo, CREATE_DT, CREATED_BY ) SELECT OriginalClaimID, LineIndex, ProcedureCode, RevenueCode, Modifier, DateOfService, DateOfServiceTo, '%s' AS CREATE_DT, '%s' AS CREATED_BY FROM ITSHome.OriginalProcedures WHERE OriginalClaimID = '%s'",
		currentDateTime, userId, claimId)
	result, err := db.Exec(tsql)
	if err != nil {
		log.Println("Error inserting new row: " + err.Error())
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateModifiedProcedures(db *sql.DB, criteria string, modifiedClaimId string) (int64, error) {
	log.Println("Inside UpdateModifiedProcedures", criteria)
	log.Println("Inside UpdateModifiedProcedures", modifiedClaimId)
	updateCriteria := strings.Split(criteria, ";")
	log.Println(len(updateCriteria))
	for i := range updateCriteria {
		log.Println("update criteria", updateCriteria[i])
		parameter := strings.Split(updateCriteria[i], "=")
		log.Println(len(parameter[0]))
		if len(parameter[0]) > 0 && (parameter[0] == "modifier" || parameter[0] == "procedureCode") {
			tsql := fmt.Sprintf("UPDATE ITSHome.ModifiedProcedures SET %s = %s WHERE Id = '%s'",
				strings.Title(parameter[0]), parameter[1], modifiedClaimId)
			result, err := db.Exec(tsql)
			if err != nil {
				log.Println("Error updating row: " + err.Error())
				return -1, err
			}
			log.Printf("Updated ID: %d successfully.", result)
		}
	}
	return 0, nil
}

func GetSubscriberId(id int) string {
	log.Println("Inside GetSubscriberId")
	var subscriberId string
	var errdb error
	condb, errdb = utilities.GetSqlConnection()
	if errdb != nil {
		log.Fatal("Open connection failed:", errdb.Error())
	}
	log.Println("Connected!")
	//defer condb.Close()

	rows, err := condb.Query("SELECT SubscriberId FROM ITSHome.ModifiedClaims WHERE Id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&subscriberId)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return subscriberId
}

func CheckIfOriginalProceduresExist(id int64) (int, error) {
	rows, err := condb.Query("SELECT OriginalClaimID FROM ITSHome.OriginalProcedures WHERE OriginalClaimID = ?", id)
	if err != nil {
		log.Println("Error reading rows: " + err.Error())
		return -1, err
	}
	//defer rows.Close()
	count := 0
	for rows.Next() {
		var name, location string
		var id int
		err := rows.Scan(&id, &name, &location)
		if err != nil {
			log.Println("Error reading rows: " + err.Error())
			return -1, err
		}
		log.Printf("ID: %d", id)
		count++
	}
	return count, nil
}

//Delete will remove a claim from the Claims list
func Delete(id string) error {
	log.Println("ID to delete", id)

	var errdb error
	condb, errdb = utilities.GetSqlConnection()
	if errdb != nil {
		log.Fatal("Open connection failed:", errdb.Error())
	}
	log.Println("Connected!")
	defer condb.Close()

	tsql := fmt.Sprintf("DELETE FROM ITSHome.ModifiedProcedures WHERE ModifiedClaimID ='%s';", id)
	result, err := condb.Exec(tsql)
	if err != nil {
		log.Println("Error deleting ModifiedProcedures row: " + err.Error())
		return err
	}

	count, err2 := result.RowsAffected()
	if err2 != nil {
		log.Println("Error getting ModifiedProcedures row affected: " + err2.Error())
		return err2
	} else {
		log.Println("Deleted ModifiedProcedures row: " + strconv.FormatInt(count, 10))
	}

	tsql = fmt.Sprintf("DELETE FROM ITSHome.ModifiedClaims WHERE Id ='%s';", id)
	result, err = condb.Exec(tsql)
	if err != nil {
		log.Println("Error deleting ModifiedClaims row: " + err.Error())
		return err
	}

	count2, err3 := result.RowsAffected()
	if err3 != nil {
		log.Println("Error deleting ModifiedClaims row: " + err3.Error())
		return err3
	} else {
		log.Println("Deleted ModifiedClaims row: " + strconv.FormatInt(count2, 10))
	}

	return nil
}

func newResultAll(id string, subscriberId string, originalClaimID string, sccfNumber string, procedureCode string, diagnosisCode string,
	modifier string, patientAccountNumber string, networkIndicator string, fromDate string, toDate string, status string,
	dosFrom string, dosTo string, createDate string, createdBy string) ModifiedClaims {
	return ModifiedClaims{
		ID:                   id,
		SubscriberId:         subscriberId,
		OriginalClaimID:      originalClaimID,
		SccfNumber:           sccfNumber,
		ProcedureCode:        procedureCode,
		DiagnosisCode:        diagnosisCode,
		Modifier:             modifier,
		PatientAccountNumber: patientAccountNumber,
		NetworkIndicator:     networkIndicator,
		FromDate:             fromDate,
		ToDate:               toDate,
		Status:               status,
		DosFrom:              dosFrom,
		DosTo:                dosTo,
		CreateDate:           createDate,
		CreatedBy:            createdBy,
	}
}
