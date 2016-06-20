package main

import (
	"flag"
	"testing"
)

func TestDefaultOptionVersion(t *testing.T) {
	if options.Version != false {
		t.Error("Expected --version defaults to false")
	}
}

func TestDefaultOptionHelp(t *testing.T) {
	if options.Help != false {
		t.Error("Expected --help defaults to false")
	}
}

func TestDefaultOptionVerbose(t *testing.T) {
	if options.Version != false {
		t.Error("Expected --verbose defaults to false")
	}
}

func TestDefaultOptionHttpBind(t *testing.T) {
	if options.HttpBind != ":12080" {
		t.Errorf("Expected --http-bind defaults to %q, was %q", ":12080", options.HttpBind)
	}
}

func TestDefaultOptionSmtpBind(t *testing.T) {
	if options.SmtpBind != ":12025" {
		t.Errorf("Expected --smtp-bind defaults to %q, was %q", ":12025", options.SmtpBind)
	}
}

func TestOptionVersion(t *testing.T) {
	flag.CommandLine.Parse([]string{"-version"})

	if options.Version != true {
		t.Error("Expected --version to be true")
	}
}

func TestOptionHelpLong(t *testing.T) {
	flag.CommandLine.Parse([]string{"-help"})

	if options.Help != true {
		t.Error("Expected --help to be true")
	}
}

func TestOptionHelpShort(t *testing.T) {
	flag.CommandLine.Parse([]string{"-h"})

	if options.Help != true {
		t.Error("Expected --help to be true")
	}
}
