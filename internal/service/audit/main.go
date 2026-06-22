package audit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/model/DTO"
	"github.com/labstack/gommon/log"
)

// PushToAudit sends an audit event to both a local file and a remote HTTP endpoint.
func PushToAudit(cfg config.Config, dto DTO.BaseAuditItem) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		log.Errorf("Fail to push audit: %s", err.Error())

		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go pushToFile(cfg, string(jsonData), &wg)
	go pushByURL(cfg, jsonData, &wg)

	wg.Wait()
}

func pushToFile(cfg config.Config, data string, wg *sync.WaitGroup) {
	defer wg.Done()
	if cfg.AuditFilePath == "" {
		return
	}

	f, err := os.OpenFile(cfg.AuditFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)

		return
	}
	defer f.Close()

	_, err = f.WriteString(data + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func pushByURL(cfg config.Config, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	if cfg.AuditURL == "" {
		return
	}

	_, err := http.Post(cfg.AuditURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
}
