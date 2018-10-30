package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func ParseParameters(parameters string) string {
	var criteria string

	fmt.Println("parameters", parameters)
	b := []byte(parameters)
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			if k == "instInputItems" || k == "profInputItems" || k == "hiddenInputItems" {
				for i, u := range vv {
					fmt.Println("Key:", i, "Value:", u)
					str := fmt.Sprintf("%v", u)
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
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	criteria = CleanParameters(criteria)
	return criteria
}

func CleanParameters(parameters string) string {
	var criteria string
	//adding hidden fields to array
	option := Option{Name: "Claim Type", Value: 23, Type: "ClaimType"}
	options.Options = append(options.Options, option)
	criteria = parameters
	for i := 0; i < len(options.Options); i++ {
		fmt.Println(options.Options[i].Type)
		parameter := regexp.MustCompile(options.Options[i].Type)

		matches := parameter.FindAllStringIndex(parameters, -1)
		fmt.Println("parameter="+options.Options[i].Type+", occurrences=", len(matches))
		if len(matches) > 1 {
			var removedValues []string
			var removeFieldName string

			// Split on comma.
			result := strings.Split(parameters, ";")

			for j := range result {
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
			// criteria = strings.Replace(criteria, ";", "';", -1)
			// criteria = strings.Replace(criteria, "=", "='", -1)
			// fmt.Println("string after adding single quotes", criteria)

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
				fmt.Println("strInClause", strInClause+strInClauseValues+")")
				criteria = criteria + strInClause + strInClauseValues + ");"
			}
		}

	}
	criteria = " AND " + strings.Replace(criteria, ";", " AND ", -1)
	return TrimSuffix(criteria, " AND ")
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}