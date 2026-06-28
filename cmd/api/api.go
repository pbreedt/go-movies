package main

import (
	"net/http"

	"github.com/pbreedt/go-movies/internal/booking"
	"github.com/pbreedt/go-movies/internal/catalog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo/v5"
)

type BookingRequest struct {
	MovieId          int    `json:"movie_id"`
	ShowTimeId       int    `json:"show_time_id"`
	Seats            int    `json:"seats"`
	CreditCardNumber string `json:"cc_number,omitempty"`
	CreditCardExpiry string `json:"cc_expiry,omitempty"`
	CreditCardCVV    string `json:"cc_cvv,omitempty"`
	PayPalEmail      string `json:"paypal_email,omitempty"`
	PayPalPassword   string `json:"paypal_password,omitempty"`
}

func main() {
	reg := prometheus.NewRegistry()

	movieRepo := catalog.NewInMemoryRepository()
	bookingRepo := booking.NewInMemoryRepository()
	bookingService := booking.NewService(bookingRepo, movieRepo)
	bookingService.InitMetrics(reg)

	go func() { http.ListenAndServe(":2112", nil) }()

	e := echo.New()

	secGroup := e.Group("/api/v1")
	secGroup.Use(AuthMiddleware)
	secGroup.Use(SecureHeaderMiddleware)

	// e.GET("/metrics", echo.WrapHandler(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	e.GET("/api/v1/movies", func(c *echo.Context) error {
		movies, err := movieRepo.GetAllMovies()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch movies"})
		}
		return c.JSON(http.StatusOK, movies)
	})

	secGroup.POST("/book", func(c *echo.Context) error {
		var req BookingRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		useCreditCard := (req.CreditCardNumber == "" || req.CreditCardExpiry == "" || req.CreditCardCVV == "")
		usePayPal := (req.PayPalEmail == "" || req.PayPalPassword == "")

		if useCreditCard && usePayPal {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payment information is required"})
		}
		if !useCreditCard && !usePayPal {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Only one payment method should be provided"})
		}

		if c.Request().Context().Err() != nil {
			return c.JSON(http.StatusRequestTimeout, map[string]string{"error": "Request timed out"})
		}

		var ticket *booking.Ticket
		if useCreditCard {
			t, err := bookingService.BookWithCreditCard(c.Request().Context(), req.MovieId, req.ShowTimeId, req.Seats, req.CreditCardNumber, req.CreditCardExpiry, req.CreditCardCVV)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to book tickets"})
			}
			ticket = t
		} else if usePayPal {
			t, err := bookingService.BookWithPayPal(c.Request().Context(), req.MovieId, req.ShowTimeId, req.Seats, req.PayPalEmail, req.PayPalPassword)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to book tickets"})
			}
			ticket = t
		}
		return c.JSON(http.StatusOK, ticket)
	})

	e.Start(":8080")
}
