package main

import (
	"github.com/twinj/uuid"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	uuid.Init()
	os.Exit(m.Run())
}
