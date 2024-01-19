package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/ninogresenz/civar/gitlab"
)

const (
	K8sPrefix   = "K8S_SECRET_"
	ScopePrefix = "# Scope: "

	// scopes
	AllScope      = "*"
	StageScope    = "staging"
	ProdScope     = "production"
	ProdTestScope = "prodtest"

	// formats
	jsonFormat   = "json"
	prettyFormat = "pretty"
	dotenvFormat = "dotenv"
)

type Service interface {
	Search()
	Get(format string, scopeFilter string)
	Create(format string, k8s bool, fileFlag string)
	Update(format string, k8s bool, fileFlag string)
}

func NewService(api gitlab.Api, cmd *cobra.Command, args []string) Service {
	return &service{api, cmd, args}
}

type service struct {
	api  gitlab.Api
	cmd  *cobra.Command
	args []string
}

func (s *service) Search() {
	data, err := s.api.Search(s.args[0])
	if err != nil {
		log.Fatalf("could not get vars: %v", err)
	}
	for _, p := range data {
		fmt.Printf("%s  %s\n", p.PathWithNamespace, p.WebUrl)
	}
}

func (s *service) Get(format string, scopeFilter string) {
	data, err := s.api.GetProjectVars(s.args[0])
	if err != nil {
		log.Fatalf("could not get vars: %v", err)
	}
	// apply scope filter if applicable
	if len(scopeFilter) > 0 {
		data = ApplyScopeFilter(data, scopeFilter)
	}
	printer := PrinterProvider(format)
	fmt.Println(printer.Print(data))
}

func (s *service) Create(format string, k8s bool, fileFlag string) {
	if format != dotenvFormat && format != jsonFormat {
		log.Fatal("format must be one of [json | dotenv]")
	}
	input := getInput(fileFlag, s.cmd.InOrStdin())
	data := parseInput(format, input)
	if k8s {
		data = AddPrefix(data)
	}
	project := s.args[0]

	existingVars, err := s.api.GetProjectVars(project)
	if err != nil {
		log.Fatalf("could not get vars: %v", err)
	}
	notCreatedVars := make(gitlab.CiVariableList, 0)
	for _, variable := range data {
		if existingVars.Includes(variable) {
			notCreatedVars.Push(variable)
			continue
		}
		_, err := s.api.CreateVar(project, variable)
		if err != nil {
			log.Fatalf("could not create variable [%s]: %v", variable.Key, err)
		}
	}
	if len(notCreatedVars) > 0 {
		printer := PrinterProvider(format)
		fmt.Println(printer.Print(notCreatedVars))
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Duplicate variables skipped: %d/%d\n", len(notCreatedVars), len(data)))
	}
}

func (s *service) Update(format string, k8s bool, fileFlag string) {
	input := getInput(fileFlag, s.cmd.InOrStdin())
	data := parseInput(format, input)
	if k8s {
		data = AddPrefix(data)
	}
	project := s.args[0]

	existingVars, err := s.api.GetProjectVars(project)
	if err != nil {
		log.Fatalf("could not get vars: %v", err)
	}

	notUpdatedVars := make(gitlab.CiVariableList, 0)
	for _, variable := range data {
		if !existingVars.Includes(variable) {
			notUpdatedVars = append(notUpdatedVars, variable)
			continue
		}
		_, err := s.api.UpdateVar(project, variable)
		if err != nil {
			log.Fatalf("could not update variable [key: %s, scope:%s]: %v", variable.Key, variable.EnvironmentScope, err)
		}
	}
	if len(notUpdatedVars) > 0 {
		prettyJson, err := json.MarshalIndent(notUpdatedVars, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(prettyJson))
		_, _ = os.Stderr.WriteString(fmt.Sprintf("variables skipped because not existent: %d/%d\n", len(notUpdatedVars), len(data)))
	}
}

func getInput(file string, stdin io.Reader) []byte {
	if len(file) > 0 {
		return getFileContent(file)
	}
	input, err := ioutil.ReadAll(stdin)
	if err != nil {
		log.Fatalf("could not get input from stdin: %v", err)
	}
	return input
}

func parseInput(format string, input []byte) []gitlab.CiVariable {
	if format == dotenvFormat {
		return ParseDotEnv(input)
	}
	var data []gitlab.CiVariable
	err := json.Unmarshal(input, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func AddPrefix(data []gitlab.CiVariable) []gitlab.CiVariable {
	for i := range data {
		if strings.Contains(data[i].Key, K8sPrefix) {
			continue
		}
		data[i].Key = K8sPrefix + data[i].Key
	}
	return data
}

func getFileContent(filepath string) []byte {
	_, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("Could not receive file info from: %v\n", filepath)
	}
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Could not read file: %v\n", filepath)
	}
	return fileContent
}

func ApplyScopeFilter(data gitlab.CiVariableList, scopeFilter string) gitlab.CiVariableList {
	var filteredList []gitlab.CiVariable
	for _, envVar := range data {
		if envVar.EnvironmentScope != scopeFilter {
			continue
		}
		filteredList = append(filteredList, envVar)
	}
	return filteredList
}

func ParseDotEnv(input []byte) []gitlab.CiVariable {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	var allBuffer bytes.Buffer
	var stagingBuffer bytes.Buffer
	var productionBuffer bytes.Buffer
	var prodtestBuffer bytes.Buffer
	scope := "*"
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), ScopePrefix) {
			scope = strings.Replace(scanner.Text(), ScopePrefix, "", 1)
			scanner.Scan()
		}
		switch scope {
		case AllScope:
			allBuffer.Write(append(scanner.Bytes(), '\n'))
		case ProdScope:
			productionBuffer.Write(append(scanner.Bytes(), '\n'))
		case StageScope:
			stagingBuffer.Write(append(scanner.Bytes(), '\n'))
		case ProdTestScope:
			prodtestBuffer.Write(append(scanner.Bytes(), '\n'))
		}
	}
	return join(
		toStruct(toMap(allBuffer), AllScope),
		toStruct(toMap(stagingBuffer), StageScope),
		toStruct(toMap(productionBuffer), ProdScope),
		toStruct(toMap(prodtestBuffer), ProdTestScope),
	)
}

func join[Type interface{}](vars ...[]Type) (allVars []Type) {
	for _, varList := range vars {
		allVars = append(allVars, varList...)
	}
	return allVars
}

func toMap(buf bytes.Buffer) map[string]string {
	env, err := godotenv.Unmarshal(buf.String())
	if err != nil {
		log.Fatal("Could not unmarshal dotenv structure")
	}
	return env
}

func toStruct(envMap map[string]string, scope string) []gitlab.CiVariable {
	var variables []gitlab.CiVariable
	for key, value := range envMap {
		variable := gitlab.CiVariable{
			Key:              key,
			Value:            value,
			EnvironmentScope: scope,
			VariableType:     "env_var",
			Protected:        false,
			Masked:           false,
		}
		variables = append(variables, variable)
	}
	return variables
}
