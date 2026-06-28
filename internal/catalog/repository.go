package catalog

import (
	"errors"
	"time"
)

type Repository interface {
	AddMovie(movie Movie) error
	GetAllMovies() ([]Movie, error)
	GetMovie(id int) (*Movie, error)
	FindMovieByTitle(title string) (*Movie, error)
	FindMoviesByGenre(genre Ganre) ([]Movie, error)
}

type InMemoryRepository struct {
	movies []Movie
}

func NewInMemoryRepository() *InMemoryRepository {
	r := &InMemoryRepository{}
	r.movies = []Movie{
		{ID: 1, Title: "Inception", Genre: SciFi, Rating: 8.8, Price: 10.0, ShowTimes: []time.Time{time.Now().Add(1 * time.Hour), time.Now().Add(3 * time.Hour)}},
		{ID: 2, Title: "The Dark Knight", Genre: Action, Rating: 9.0, Price: 12.0, ShowTimes: []time.Time{time.Now().Add(1 * time.Hour), time.Now().Add(3 * time.Hour)}},
		{ID: 3, Title: "Interstellar", Genre: SciFi, Rating: 8.6, Price: 11.0, ShowTimes: []time.Time{time.Now().Add(1 * time.Hour), time.Now().Add(3 * time.Hour)}},
		{ID: 4, Title: "The Shawshank Redemption", Genre: Drama, Rating: 9.3, Price: 9.0, ShowTimes: []time.Time{time.Now().Add(1 * time.Hour), time.Now().Add(3 * time.Hour)}},
		{ID: 5, Title: "The Godfather", Genre: Drama, Rating: 9.2, Price: 10.0, ShowTimes: []time.Time{time.Now().Add(1 * time.Hour), time.Now().Add(3 * time.Hour)}},
	}
	return r
}

func (r *InMemoryRepository) AddMovie(movie Movie) error {
	r.movies = append(r.movies, movie)
	return nil
}

func (r *InMemoryRepository) GetAllMovies() ([]Movie, error) {
	return r.movies, nil
}

func (r *InMemoryRepository) GetMovie(id int) (*Movie, error) {
	for _, movie := range r.movies {
		if movie.ID == id {
			return &movie, nil
		}
	}
	return nil, errors.New("movie not found")
}

func (r *InMemoryRepository) FindMovieByTitle(title string) (*Movie, error) {
	for _, movie := range r.movies {
		if movie.Title == title {
			return &movie, nil
		}
	}
	return nil, errors.New("movie not found")
}

func (r *InMemoryRepository) FindMoviesByGenre(genre Ganre) ([]Movie, error) {
	var filteredMovies []Movie
	for _, movie := range r.movies {
		if movie.Genre == genre {
			filteredMovies = append(filteredMovies, movie)
		}
	}
	if len(filteredMovies) == 0 {
		return nil, errors.New("no movies found for the given genre")
	}
	return filteredMovies, nil
}
