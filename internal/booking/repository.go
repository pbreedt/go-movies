package booking

import (
	"errors"
	"slices"
)

type Repository interface {
	SaveBooking(ticket *Ticket) error
	CancelBooking(ticketId int) error
	GetBookingById(ticketId int) (*Ticket, error)
}

type InMemoryRepository struct {
	bookings []Ticket
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		bookings: []Ticket{},
	}
}

func (r *InMemoryRepository) SaveBooking(ticket *Ticket) error {
	ticket.TicketId = len(r.bookings) // len is already 1 more than index, so ID will be equal to index
	r.bookings = append(r.bookings, *ticket)
	return nil
}

func (r *InMemoryRepository) CancelBooking(ticketId int) error {
	for i, booking := range r.bookings {
		if booking.TicketId == ticketId {
			r.bookings = slices.Delete(r.bookings, i, i+1)
			return nil
		}
	}
	return errors.New("booking not found")
}

func (r *InMemoryRepository) GetBookingById(ticketId int) (*Ticket, error) {
	for _, booking := range r.bookings {
		if booking.TicketId == ticketId {
			return &booking, nil
		}
	}
	return nil, errors.New("booking not found")
}
