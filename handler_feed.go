package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jesiahharris/rss-agg/internal/database"
)

const createFeedsHTML = `
	<h1> New feed </h1>
	<form action="/feeds" method="POST"> 
		<table> 
			<tr>
				<td> Authorization </td>
				<td><input type="text" name="authorization" /></td>
			</tr>
			<tr>
				<td> Name </td>
				<td><input type="text" name="name" /></td>
			</tr>
			<tr>
				<td> URL </td>
				<td><input type="text" name="url" /></td>
			</tr>
		</table>
		<button type="submit"> Create feed</button>
	</form>

`

func (apiCfg *apiConfig) handlerNewFeed(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(createFeedsHTML))

	tmpl.Execute(w, nil)
}

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedtoFeed(feed))
}

func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIDStr := chi.URLParam(r, "feedID")
	feedId, err := uuid.Parse(feedIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeed(r.Context(), feedId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}
	respondWithJSON(w, 204, struct{}{})
}

const getFeedsHTML = `
<h1> Feeds </h1> 
<dl>
{{range .Feeds}}
<dt><strong>{{.Name}}</strong></dt>
<dd> URL:{{.Url}}</dd>
<dd> ID: {{.ID}}</dd>
<dd> Updated At: {{.UpdatedAt}} </dd>
{{end}}
`

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Feeds []Feed
	}

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	tmpl := template.Must(template.New("").Parse(getFeedsHTML))

	tmpl.Execute(w, data{Feeds: databaseFeedstoFeeds(feeds)})

	//respondWithJSON(w, 200, databaseFeedstoFeeds(feeds))
}
