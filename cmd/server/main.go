package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pbreedt/go-movies/internal/booking"
	"github.com/pbreedt/go-movies/internal/catalog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// metrics:
	reg := prometheus.NewRegistry()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	movieRepo := catalog.NewInMemoryRepository()
	bookingRepo := booking.NewInMemoryRepository()
	bookingService := booking.NewService(bookingRepo, movieRepo)
	bookingService.InitMetrics(reg)

	http.ListenAndServe(":2112", nil)

	// browse movies
	movies, _ := movieRepo.GetAllMovies()
	for _, m := range movies {
		println(m.Title)
	}

	// book a ticket
	ticket, err := bookingService.BookWithCreditCard(context.Background(), 1, 0, 2, "1234-5678-9012-3456", "12/25", "123")
	if err != nil {
		println("Error booking ticket:", err.Error())
		return
	}
	fmt.Printf("Booked ticket %d for movie ID: %d\n", ticket.TicketId, ticket.MovieId)

}
