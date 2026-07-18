package handler

import (
	"net"
	"net/http"

	"github.com/Vaha95/golang_pet/internal/config"
	"github.com/Vaha95/golang_pet/internal/service/stats"
	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

// GetStatsHandler returns an Echo handler that responds with internal stats.
func GetStatsHandler(cfg config.StorageConfig, l *zap.SugaredLogger) func(c *echo.Context) error {
	return func(c *echo.Context) error {
		if !isRequestTrusted(c, cfg.Config.TrustedSubnet) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
		}

		resp, err := stats.CollectFromConfig(cfg)
		if err != nil {
			l.Errorf("Failed to collect stats: %v", err)

			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func isRequestTrusted(c *echo.Context, trustedSubnet string) bool {
	if trustedSubnet == "" {
		return false
	}

	ip := net.ParseIP(c.Request().Header.Get("X-Real-IP"))
	if ip == nil {
		return false
	}

	return isInSubnet(ip, trustedSubnet)
}

func isInSubnet(ip net.IP, subnetCIDR string) bool {
	_, subnet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return false
	}

	return subnet.Contains(ip)
}
