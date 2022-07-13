package monitoring

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Monitoring is an instance of STUNner monitoring
type MonitoringServer struct {
	httpServer *http.Server
	Endpoint   string
}

// NewMonitoring initiates the monitoring subsystem
func NewMonitoringServer(endpoint string) (*MonitoringServer, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to parse: %s", endpoint))
	}

	addr := u.Hostname()
	port := u.Port()
	if port != "" {
		addr = addr + ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: promhttp.Handler(),
	}

	m := &MonitoringServer{
		httpServer: server,
		Endpoint:   endpoint,
	}

	return m, nil
}

func (m *MonitoringServer) Start() { // specify config, create new server; move init here?
	if m.Endpoint == "" {
		return
	}
	// serve Prometheus metrics over HTTP
	go func() {
		m.httpServer.ListenAndServe()
	}()
}

func (m *MonitoringServer) Stop() {
	if m.httpServer == nil {
		return
	}
	m.httpServer.Shutdown(context.Background())
}
