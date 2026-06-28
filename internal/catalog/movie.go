package catalog

import "time"

type Ganre string

const (
	Action  Ganre = "Action"
	Comedy  Ganre = "Comedy"
	Drama   Ganre = "Drama"
	SciFi   Ganre = "Sci-Fi"
	Romance Ganre = "Romance"
)

type Movie struct {
	ID        int
	Title     string
	Genre     Ganre
	Rating    float64
	Price     float64
	ShowTimes []time.Time
}
