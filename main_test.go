package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestParseArgs_Correct(t *testing.T) {
	is := is.New(t)

	port := "8000"
	addrs := []string{"host1", "host2", "host3"}
	args := append([]string{port}, addrs...)
	appArgs, err := parseArgs(args)

	is.NoErr(err)
	is.True(appArgs != nil)

	is.Equal(appArgs.Port, port)
	is.Equal(appArgs.Addrs, addrs)
}

func TestParseArgs_NoArgs(t *testing.T) {
	is := is.New(t)

	appArgs, err := parseArgs(nil)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}

func TestParseArgs_BadPort(t *testing.T) {
	is := is.New(t)

	args := []string{"Bad port", "Host1"}
	appArgs, err := parseArgs(args)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}

func TestParseArgs_PortOutOfRange(t *testing.T) {
	is := is.New(t)

	args := []string{"-1", "Host1"}
	appArgs, err := parseArgs(args)

	is.True(err != nil)
	is.Equal(appArgs, nil)
}
