package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

var envVarFilterRegex *regexp.Regexp

func init() {
	regexp.MustCompile("^(_|DISPLAY|MAIL|USER|TERM|HOME|LOGNAME|SHELL|SHLVL|PWD|SSH_.+)=")
}

// if we got an error, panic and log it. otherwise do nothing
func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

// same as check but do not panic
func checkNoPanic(e error) {
	if e != nil {
		log.Println("ERROR:", e)
	}
}

// filter the current environment variables according to the regex
func filteredEnvironmentVars() []string {
	filteredVars := []string{}
	for _, envVarLine := range os.Environ() {
		if !envVarFilterRegex.MatchString(envVarLine) {
			filteredVars = append(filteredVars, envVarLine)
		}
	}
	return filteredVars
}

// fetches an environment variable. if the variable is not set, it returns a default
func fetchEnvValue(key string, fallback string) string {
	value, isset := os.LookupEnv(key)
	if !isset {
		return fallback
	} else {
		return value
	}
}

func fetchEnvValueInt(key string, fallback int) int {
	value, isset := os.LookupEnv(key)
	if !isset {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
