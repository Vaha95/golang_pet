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

var wgPool = &sync.Pool{New: func() interface{} {
	return new(sync.WaitGroup)
}}

var fileMu sync.Mutex

// AuditWorker receives audit events on a channel and persists them
// to a local file and/or a remote HTTP endpoint configured in cfg.
type AuditWorker struct {
	cfg     config.Config
	AuditCh chan DTO.BaseAuditItem
}

// CreateAuditWorker returns a new AuditWorker with the given config
// and a buffered channel for receiving audit events.
func CreateAuditWorker(cfg config.Config) AuditWorker {
	return AuditWorker{
		cfg:     cfg,
		AuditCh: make(chan DTO.BaseAuditItem),
	}
}

// Run starts the background goroutine that drains AuditCh and
// forwards each event to PushToAudit.
func (aw *AuditWorker) Run() {
	go aw.listen()
}

func (aw *AuditWorker) listen() {
	for {
		msg := <-aw.AuditCh
		PushToAudit(aw.cfg, msg)
	}

}

// PushToAudit sends an audit event to both a local file and a remote HTTP endpoint.
func PushToAudit(cfg config.Config, dto DTO.BaseAuditItem) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		log.Errorf("Fail to push audit: %s", err.Error())

		return
	}

	wg := wgPool.Get().(*sync.WaitGroup)
	wg.Add(2)

	go pushToFile(cfg, string(jsonData), wg)
	go pushByURL(cfg, jsonData, wg)

	wg.Wait()
	wgPool.Put(wg)
}

func pushToFile(cfg config.Config, data string, wg *sync.WaitGroup) {
	defer wg.Done()
	if cfg.AuditFilePath == "" {
		return
	}

	f, err := os.OpenFile(cfg.AuditFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)

		return
	}
	defer f.Close()

	fileMu.Lock()
	_, err = f.WriteString(data + "\n")
	fileMu.Unlock()
	if err != nil {
		log.Error(err)
	}
}

func pushByURL(cfg config.Config, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	if cfg.AuditURL == "" {
		return
	}

	_, err := http.Post(cfg.AuditURL, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error(err)
	}
}
