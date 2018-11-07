package utilities

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	options Options
)

func init() {
	loadOptions()
}

func loadOptions() {
	// Open our jsonFile
	jsonFile, err := os.Open("./ui/src/assets/options.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened options.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened file as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'options' which we defined above
	json.Unmarshal(byteValue, &options)

	// iterate through every option within our options array and
	// print out the option Type, the name, and the value
	for i := 0; i < len(options.Options); i++ {
		fmt.Println("Option Name: " + options.Options[i].Name)
		fmt.Println("Option Value: " + strconv.Itoa(options.Options[i].Value))
		fmt.Println("Option Type: " + options.Options[i].Type)
	}
}

// Options struct which contains
// an array of options
type Options struct {
	Options []Option `json:"options"`
}

type Option struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Type  string `json:"type"`
}

func GetSqlConnection() (*sql.DB, error) {
	condb, errdb := sql.Open("mssql", "server=SQLDEV34\\SQL_DEV34;user id=;password=;database=zdb63q_itshc_syst")
	if errdb != nil {
		fmt.Println(" Error open db:", errdb.Error())
	}

	return condb, nil
}

func ParseParameters(parameters string) string {
	var criteria string
	var inClauseCriteria string
	fmt.Println("parameters", parameters)
	b := []byte(parameters)
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
			criteria += k + "='" + vv + "';"
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			if k == "instInputItems" || k == "profInputItems" || k == "hiddenInputItems" || k == "selectedActiveInstitutionalClaimIds" {
				var strInClauseValues string
				for i, u := range vv {
					fmt.Println("Key:", i, "Value:", u)
					str := fmt.Sprintf("%v", u)
					if k == "selectedActiveInstitutionalClaimIds" || k == "selectedActiveProfessionalClaimIds" {
						strInClauseValues += str
						if len(vv)-i > 1 {
							strInClauseValues += ","
						}
						inClauseCriteria = strInClauseValues
						fmt.Println("inClauseCriteria", inClauseCriteria)
					} else {
						out := strings.TrimLeft(strings.TrimRight(str, "]"), "map[")
						inputArray := strings.Fields(out)
						fmt.Println("param1 before check", inputArray[0])
						fmt.Println("param2 before check", inputArray[1])
						startsWith := strings.HasPrefix(inputArray[0], "inputName")
						var param1 string
						var param2 string
						if startsWith {
							param1 = inputArray[0]
							param2 = inputArray[1]
						} else {
							param1 = inputArray[1]
							param2 = inputArray[0]
						}
						fmt.Println("param1 after check", inputArray[0])
						fmt.Println("param2 after check", inputArray[1])

						paramValue := strings.SplitAfter(param1, ":")
						paramName := strings.SplitAfter(param2, ":")
						criteria += paramName[1] + "=" + paramValue[1] + ";"
						fmt.Println("paramName", paramName[1])
						fmt.Println("paramValue", paramValue[1])
					}
				}
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	if len(inClauseCriteria) > 0 {
		criteria = criteria + "&" + inClauseCriteria + ";"
	}
	return criteria
}

func CleanParameters(parameters string) string {
	fmt.Println("entering CleanParameters()...")
	criteria := parameters
	for i := 0; i < len(options.Options); i++ {
		fmt.Println(options.Options[i].Type)
		parameter := regexp.MustCompile(options.Options[i].Type)

		matches := parameter.FindAllStringIndex(parameters, -1)
		fmt.Println("parameter="+options.Options[i].Type+", occurrences=", len(matches))
		fmt.Println("matches=", matches)
		if len(matches) > 1 {
			var removedValues []string
			var removeFieldName string

			// Split on comma.
			result := strings.Split(parameters, ";")
			fmt.Println("counter is at ", i)
			for j := range result {
				fmt.Println("second counter is at ", j)
				fmt.Println("result", result[j])
				field := result[j]
				fieldName := strings.Split(field, "=")
				fmt.Println("fieldName", fieldName[0])
				if len(fieldName[0]) > 0 && fieldName[0] == options.Options[i].Type {
					removeFieldName = fieldName[0]
					fmt.Println("removeField name", removeFieldName)
					fmt.Println("removeField value", fieldName[1])
					fmt.Println("string before removing", criteria)
					fmt.Println("string to be removed", removeFieldName+"="+fieldName[1])
					removedValues = append(removedValues, fieldName[1])
					fmt.Println("removedValues", removedValues)
					criteria = strings.Replace(criteria, removeFieldName+"="+fieldName[1]+";", "", -1)
					fmt.Println("string after removing", criteria)
				}
			}

			if len(removedValues) == len(matches) {
				strInClause := removeFieldName + " IN ("
				var strInClauseValues string
				//add criteria back in
				for k := 0; k < len(removedValues); k++ {
					strInClauseValues += "'" + removedValues[k] + "'"
					if len(removedValues)-k > 1 {
						strInClauseValues += ","
					}
				}
				criteria = criteria + strInClause + strInClauseValues + ");"
			}
		}

	}
	criteria2 := strings.Replace(criteria, ";", "';", -1)
	criteria2 = strings.Replace(criteria2, "=", "='", -1)
	criteria2 = strings.Replace(criteria2, ")'", ")", -1)
	//TrimSuffix(criteria2, "';")
	fmt.Println("string after adding single quotes", criteria2)
	fmt.Println("original criteria", criteria)
	criteria = " WHERE " + strings.Replace(criteria2, ";", " AND ", -1)
	return TrimSuffix(criteria, " AND ")
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
