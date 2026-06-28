package booking

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/pbreedt/go-movies/internal/catalog"
	"github.com/pbreedt/go-movies/internal/payment"
)

type Service struct {
	BookingRepo Repository
	MovieRepo   catalog.Repository
	metrics     *metrics
}

type Ticket struct {
	TicketId int
	MovieId  int
	ShowTime time.Time
	Seats    int
	Price    float64
}

func NewService(bookingRepo Repository, movieRepo catalog.Repository) *Service {
	return &Service{
		BookingRepo: bookingRepo,
		MovieRepo:   movieRepo,
	}
}

func (s *Service) BookWithCreditCard(ctx context.Context, movieId int, showTimeId int, seats int, cardNumber string, expiryDate string, cvv string) (*Ticket, error) {
	return s.create(ctx, movieId, showTimeId, seats, &payment.CreditCardProvider{
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
	})
}

func (s *Service) BookWithPayPal(ctx context.Context, movieId int, showTimeId int, seats int, email string, password string) (*Ticket, error) {
	return s.create(ctx, movieId, showTimeId, seats, &payment.PayPalProvider{
		Email:    email,
		Password: password,
	})
}

func (s *Service) create(ctx context.Context, movieId int, showTimeId int, seats int, payment payment.Provider) (*Ticket, error) {
	var err error
	defer func() {
		slog.Info("booking attempt",
			"movie_id", movieId,
			"show_time_id", showTimeId,
			"seats", seats,
			"error", err,
		)
		if err != nil {
			if s.metrics != nil {
				s.metrics.bookingErrors.Inc()
			}
		}
	}()

	if s.metrics != nil {
		s.metrics.bookingRequests.Inc()
	}

	if ctx.Err() != nil {
		err = ctx.Err()
		return nil, err
	}

	m, err := s.MovieRepo.GetMovie(movieId)
	if err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}

	price := m.Price * float64(seats)
	err = payment.Pay(price)
	if err != nil {
		return nil, fmt.Errorf("create ticket: %w", err)
	}

	ticket := &Ticket{
		// no ID, ID assigned by BookingRepo
		MovieId:  movieId,
		ShowTime: m.ShowTimes[showTimeId],
		Seats:    seats,
		Price:    price,
	}

	err = s.BookingRepo.SaveBooking(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *Service) Cancel(ticketId int) error {
	// TODO: implement refund logic here

	return s.BookingRepo.CancelBooking(ticketId)
}

func (s *Service) GetById(ticketId int) (*Ticket, error) {
	return s.BookingRepo.GetBookingById(ticketId)
}
