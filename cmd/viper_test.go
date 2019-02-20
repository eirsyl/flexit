package cmd

import (
	"testing"
)

func TestUppercaseName(t *testing.T) {
	name := "redisServer"
	if uppercaseName(name) != "REDIS_SERVER" {
		t.Error("redisServer is not converted to REDIS_SERVER")
	}

	name = "redisDB"
	if uppercaseName(name) != "REDIS_DB" {
		t.Error("redisDB is not converted to REDIS_DB")
	}
}
