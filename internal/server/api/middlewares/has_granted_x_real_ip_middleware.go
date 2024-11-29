package middlewares

import (
	"fmt"
	"net"
	"net/http"

	"github.com/AndrXxX/go-metrics-collector/internal/services/logger"
)

const headerXRealIP = "X-Real-Ip"

type hasGrantedXRealIPMiddleware struct {
	TrustedSubnet string
}

// Handler возвращает http.HandlerFunc
func (m *hasGrantedXRealIPMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.check(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *hasGrantedXRealIPMiddleware) check(r *http.Request) bool {
	if m.TrustedSubnet == "" {
		return true
	}
	ip := r.Header.Get(headerXRealIP)
	if ip == "" {
		return true
	}
	_, ipNet, err := net.ParseCIDR(m.TrustedSubnet)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("TrustedSubnet %s is not valid", m.TrustedSubnet))
		return true
	}
	if !ipNet.Contains(net.ParseIP(ip)) {
		return false
	}
	return true
}

// HasGrantedXRealIPOr403 возвращает middleware, которая проверяет что переданный IP-адрес
// в заголовке запроса X-Real-IP входит в доверенную подсеть
func HasGrantedXRealIPOr403(trustedSubnet string) *hasGrantedXRealIPMiddleware {
	return &hasGrantedXRealIPMiddleware{TrustedSubnet: trustedSubnet}
}
