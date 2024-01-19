package service_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"

	"github.com/ninogresenz/civar/gitlab"
	"github.com/ninogresenz/civar/service"
)

func TestGet(t *testing.T) {

}

func TestApplyScopeFilter(t *testing.T) {
	output := service.ApplyScopeFilter(getVars(), "staging")
	cupaloy.SnapshotT(t, output)
}

func TestParseDotEnv(t *testing.T) {
	input := []byte(`# Scope: *
TEST_KEY1="MY_VARIABLE1"
TEST_KEY2="MY_VARIABLE2"
TEST_KEY3="MY_VARIABLE3"

# Scope: staging
TEST_KEY1="MY_VARIABLE1"
TEST_KEY2="MY_VARIABLE2"
TEST_KEY3="MY_VARIABLE3"

# Scope: production
TEST_KEY1="MY_VARIABLE1"
TEST_KEY2="MY_VARIABLE2"
TEST_KEY3="MY_VARIABLE3"
`)

	actual := service.ParseDotEnv(input)

	assert.ElementsMatch(t, actual, []gitlab.CiVariable{
		{
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		},
	})
}

func TestAddPrefix(t *testing.T) {
	vars := []gitlab.CiVariable{
		{
			Key:              "TEST_VAR1",
			VariableType:     "env_var",
			Value:            "TEST_VAL1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		},
		{
			Key:              "TEST_VAR2",
			VariableType:     "env_var",
			Value:            "TEST_VAL2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		},
	}
	vars = service.AddPrefix(vars)
	cupaloy.SnapshotT(t, vars)
}

func getVars() gitlab.CiVariableList {
	return gitlab.CiVariableList{
		{
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "*",
		}, {
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "staging",
		}, {
			Key:              "TEST_KEY1",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE1",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		}, {
			Key:              "TEST_KEY2",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE2",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		}, {
			Key:              "TEST_KEY3",
			VariableType:     "env_var",
			Value:            "MY_VARIABLE3 with a very very very very very very very very very very very very very very very very very very very very very very very very very very very very long name",
			Protected:        false,
			Masked:           false,
			EnvironmentScope: "production",
		},
	}
}
