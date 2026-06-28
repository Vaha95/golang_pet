package handler

import (
	"log"
	"os"
	"testing"

	"github.com/Vaha95/golang_pet/internal/repository"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	clearFileCache()

	os.Exit(exitCode)
}

func clearFileCache() {
	dirPath := repository.MainDir

	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Fatalf("Error deleting directory: %v", err)
	}
}
