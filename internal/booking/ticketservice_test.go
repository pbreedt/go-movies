package booking

import (
	"testing"
	"time"

	"github.com/pbreedt/go-movies/internal/catalog"
)

var mockMovie = catalog.Movie{
	ID:    0,
	Title: "Mock Movie",
	ShowTimes: []time.Time{
		time.Date(2024, time.June, 1, 19, 0, 0, 0, time.UTC),
	},
	Price: 10.0,
}

type MockPaymentProvider struct{}

func (m *MockPaymentProvider) Pay(amount float64) error {
	// Simulate successful payment processing
	return nil
}

type MockCatalogRepository struct{}

func (m *MockCatalogRepository) GetMovie(movieId int) (*catalog.Movie, error) {
	// Simulate fetching a movie from the catalog
	return &mockMovie, nil
}
func (m *MockCatalogRepository) GetAllMovies() ([]catalog.Movie, error) {
	// Simulate fetching all movies from the catalog
	return []catalog.Movie{
		mockMovie,
	}, nil
}
func (m *MockCatalogRepository) AddMovie(movie catalog.Movie) error {
	// Simulate adding a movie to the catalog
	return nil
}
func (m *MockCatalogRepository) FindMovieByTitle(title string) (*catalog.Movie, error) {
	// Simulate finding a movie by title
	return &mockMovie, nil
}
func (m *MockCatalogRepository) FindMoviesByGenre(genre catalog.Ganre) ([]catalog.Movie, error) {
	// Simulate finding movies by genre
	return []catalog.Movie{
		mockMovie,
	}, nil
}

type MockBookingRepository struct{}

func (m *MockBookingRepository) SaveBooking(ticket *Ticket) error {
	// Simulate saving the ticket to the repository
	return nil
}
func (m *MockBookingRepository) GetBookingById(ticketId int) (*Ticket, error) {
	// Simulate fetching a ticket by ID
	return &Ticket{
		TicketId: ticketId,
		MovieId:  1,
		ShowTime: time.Date(2024, time.June, 1, 19, 0, 0, 0, time.UTC),
		Seats:    2,
		Price:    20.0,
	}, nil
}
func (m *MockBookingRepository) CancelBooking(ticketId int) error {
	// Simulate canceling a booking
	return nil
}

func TestCreateTicket(t *testing.T) {
	// Setup
	movieRepo := &MockCatalogRepository{}
	bookingRepo := &MockBookingRepository{}
	service := NewService(bookingRepo, movieRepo) // Pass mock repositories
	movie, _ := movieRepo.GetMovie(1)             // Get a mock movie

	// Test data
	movieId := movie.ID
	showTimeId := 0 // Use a valid show time ID if needed
	showTime := movie.ShowTimes[showTimeId]
	seats := 2
	payment := &MockPaymentProvider{}
	price := movie.Price * float64(seats)

	// Execute
	ticket, err := service.create(t.Context(), movieId, showTimeId, seats, payment)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ticket == nil {
		t.Fatal("Expected a ticket, got nil")
	}
	if ticket.MovieId != movieId {
		t.Errorf("Expected movie ID %d, got %d", movieId, ticket.MovieId)
	}
	if ticket.ShowTime != showTime {
		t.Errorf("Expected show time %v, got %v", showTime, ticket.ShowTime)
	}
	if ticket.Price != price {
		t.Errorf("Expected price %f, got %f", price, ticket.Price)
	}
}
