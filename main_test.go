package main

import (
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestParseArgs_Correct(t *testing.T) {
	is := is.New(t)

	testPort := 8000
	testHosts := []string{"host1", "host2", "host3"}
	testArgs := append([]string{strconv.Itoa(testPort)}, testHosts...)
	appArgs, err := parseArgs(testArgs)

	is.NoErr(err)
	is.True(appArgs != nil)

	is.Equal(appArgs.port, testPort)
	is.Equal(appArgs.hosts, testHosts)
}

func TestParseArgs_NoArgs(t *testing.T) {
	is := is.New(t)

	appArgs, err := parseArgs(nil)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}

func TestParseArgs_BadPort(t *testing.T) {
	is := is.New(t)

	testArgs := []string{"Bad port", "Host1"}
	appArgs, err := parseArgs(testArgs)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}

func TestParseArgs_PortOutOfRange(t *testing.T) {
	is := is.New(t)

	testArgs := []string{"-1", "Host1"}
	appArgs, err := parseArgs(testArgs)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}
