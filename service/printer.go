package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"

	"github.com/ninogresenz/civar/gitlab"
)

type CiPrinter interface {
	Print(data gitlab.CiVariableList) string
}

func PrinterProvider(format string) (printer CiPrinter) {
	switch format {
	case jsonFormat:
		// print raw json
		return jsonPrinter{}
	case prettyFormat:
		// print as table
		return prettyPrinter{}
	case dotenvFormat:
		// print dotenv format
		return dotenvPrinter{}
	default:
		log.Fatalf("Not a valid format: %s", format)
	}
	return printer
}

// DotenvPrinter prints values as a dotenv file
type dotenvPrinter struct{}

func (p dotenvPrinter) Print(data gitlab.CiVariableList) string {
	var scopes []string
	scopeIds := []string{"*", "staging", "production", "prodtest"}
	splitVars := p.splitVarsByScope(data)
	for _, scopeId := range scopeIds {
		str, err := godotenv.Marshal(p.toMap(splitVars[scopeId]))
		if err != nil {
			log.Fatal("could not marshal data to env scopeId")
		}
		var b bytes.Buffer
		b.WriteString(fmt.Sprintf("# Scope: %v\n", scopeId))
		b.WriteString(str)
		scopes = append(scopes, b.String())
	}
	dotenvString := strings.Join(scopes, "\n\n")
	return dotenvString
}

func (p dotenvPrinter) splitVarsByScope(data gitlab.CiVariableList) map[string]gitlab.CiVariableList {
	var scopeMap = make(map[string]gitlab.CiVariableList)
	for _, variable := range data {
		scope := variable.EnvironmentScope
		scopeSlice, present := scopeMap[scope]
		if !present {
			scopeMap[scope] = gitlab.CiVariableList{}
			scopeSlice = scopeMap[scope]
		}
		variable.Key = strings.Replace(variable.Key, "K8S_SECRET_", "", 1)
		scopeSlice = append(scopeSlice, variable)
		scopeMap[scope] = scopeSlice
	}
	return scopeMap
}

func (p dotenvPrinter) toMap(data gitlab.CiVariableList) map[string]string {
	varMap := make(map[string]string)
	for _, variable := range data {
		varMap[variable.Key] = variable.Value
	}
	return varMap
}

// PrettyPrinter prints values as a table
type prettyPrinter struct{}

func (p prettyPrinter) Print(data gitlab.CiVariableList) string {
	var buf bytes.Buffer
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Scope", "Key", "Value", "Type", "Masked", "Protected"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	for _, v := range data {
		table.Append([]string{
			v.EnvironmentScope,
			v.Key,
			v.Value,
			v.VariableType,
			strconv.FormatBool(v.Masked),
			strconv.FormatBool(v.Protected),
		})
	}
	table.Render()
	return buf.String()
}

// JsonPrinter prints values as json
type jsonPrinter struct{}

func (p jsonPrinter) Print(data gitlab.CiVariableList) string {
	prettyJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(prettyJson)
}
