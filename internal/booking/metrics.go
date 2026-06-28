package booking

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	bookingRequests prometheus.Counter
	bookingErrors   prometheus.Counter
}

func (s *Service) InitMetrics(p *prometheus.Registry) {
	s.metrics = NewBookingMetrics(p)
}

func NewBookingMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		bookingRequests: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "movies_booking_requests_total",
			Help: "The total number of booking requests",
		}),
		bookingErrors: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "movies_booking_errors_total",
			Help: "The total number of booking errors",
		}),
	}
	return m
}
