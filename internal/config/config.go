package config

import (
	"flag"
	"os"
	"strconv"
)

// Config holds application configuration read from environment variables or flags.
type Config struct {
	ListenHost    string
	URLHost       string
	FilePath      string
	DbDSN         string
	AuditFilePath string
	AuditURL      string
	EnableHttps   bool
}

func GetMainConfig() Config {
	listenHostENV := os.Getenv("SERVER_ADDRESS")
	urlHostENV := os.Getenv("BASE_URL")
	filePathENV := os.Getenv("FILE_STORAGE_PATH")
	dsnENV := os.Getenv("DATABASE_DSN")
	auditFileENV := os.Getenv("AUDIT_FILE")
	auditURLENV := os.Getenv("AUDIT_URL")
	enableHttpsENV := os.Getenv("ENABLE_HTTPS")

	listenHostFlag := flag.String("a", `localhost:8080`, "Host for app")
	urlHostFlag := flag.String("b", `http://localhost:8080`, "Host for url")
	filePathFlag := flag.String("f", ``, "Storage file path")
	dsnFlag := flag.String("d", ``, "Database dsn")
	auditFileFlag := flag.String("audit-file", ``, "Audit file path")
	auditURLFlag := flag.String("audit-url", ``, "Audit URL path")
	enableHttpsFlag := flag.Bool("s", false, "Enable Https")
	flag.Parse()

	listenHost := listenHostENV
	if listenHost == "" {
		listenHost = *listenHostFlag
	}

	urlHost := urlHostENV
	if urlHost == "" {
		urlHost = *urlHostFlag
	}

	filePath := filePathENV
	if filePath == "" {
		filePath = *filePathFlag
	}

	dsn := dsnENV
	if dsn == "" {
		dsn = *dsnFlag
	}

	auditFile := auditFileENV
	if auditFile == "" {
		auditFile = *auditFileFlag
	}

	auditURL := auditURLENV
	if auditURL == "" {
		auditURL = *auditURLFlag
	}

	enableHttps, err := strconv.ParseBool(enableHttpsENV)
	if err != nil {
		enableHttps = *enableHttpsFlag
	}

	return Config{
		ListenHost:    listenHost,
		URLHost:       urlHost,
		FilePath:      filePath,
		DbDSN:         dsn,
		AuditFilePath: auditFile,
		AuditURL:      auditURL,
		EnableHttps:   enableHttps,
	}
}
