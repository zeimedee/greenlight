package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zeimedee/greenlight/internal/data"
	"github.com/zeimedee/greenlight/internal/validator"
)

//Add a createMovieHandler for the "POST /v1/movies" endpoint. for now we simply returna plain-text placeholder response.

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

//Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve the interpolated "id" parameter from the
//current URL include it in a placeholder response

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	//when httprouter is parsing a request, any interpolated URL parameters will be stored in the request
	//context.
	id, err := app.readIDParameter(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Panicln(err)
		http.Error(w, "the server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
