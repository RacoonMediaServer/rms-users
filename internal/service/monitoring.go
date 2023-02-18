package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	deviceRequestsCounter        *prometheus.CounterVec
	unknownDeviceRequestsCounter prometheus.Counter
)

func init() {
	deviceRequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "rms",
		Name:      "device_requests",
		Help:      "Amount of device's requests",
	}, []string{"device"})

	unknownDeviceRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "rms",
		Name:      "unknown_device_requests",
		Help:      "Amount of unknown device's requests",
	})
}
