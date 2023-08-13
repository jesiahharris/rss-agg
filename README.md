# rss-agg
An rss aggregator built with Go

## Migrate DB ##
This project is configred to use Postgres and the [Goose](https://github.com/pressly/goose) database migration tool. To migrate the database, navigate to the sql/schema folder. To create the database, run
`goose postgres postgres://{username}:{password}@localhost:{PORT}/rssagg up` 
 To drop the database, run 
`goose postgres postgres://{username}:{password}@localhost:{PORT}/rssagg down`

## Make a Request ##
Requests can be made via API clients (such as Postman or Thunder Client) to POST, GET, and DEL users, feeds, feed_follows, and so on. 

Example request: 
  ```
	POST http://localhost:8080/v1/feeds
	Headers: Authorization: ApiKey {ApiKey}
	Body: {
		  "Name": "Blog Name",
		  "URL": "https://blog.com" 
		}
  ```
 
Refer to the router and handlers for more info
