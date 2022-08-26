package main

import (
	"database/sql"
	"encoding/json"
	"runtime/debug"
	"time"

	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var db *sql.DB

// This function will make a connection to the database only once.
func init() {
	var err error

	connStr := "postgres://postgres:admin@localhost:2000/HELLO?sslmode=disable"
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("Connected to database")
	fmt.Println(db.Stats().OpenConnections)
	fmt.Println(db.Stats())
}

type sandbox struct {
	a int64
	b int64
	c int64
	d int64
	e int64
	f int64
	g int64
	h int64
}

// var (
// 	locations = []Geolocation{
// 		{-97, 37.819929, -122.478255},
// 		{1899, 39.096849, -120.032351},
// 		{2619, 37.865101, -119.538329},
// 		{42, 33.812092, -117.918974},
// 		{15, 37.77493, -122.419416},
// 	}
// )

func main() {

	rows, err := db.Query("SELECT * FROM newtable")
	if err != nil {
		echo.NewHTTPError(500, "could not fetch")
		return
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		for rows.Next() {
			snb := sandbox{}
			err := rows.Scan(&snb.a, &snb.b, &snb.c, &snb.d, &snb.e, &snb.f, &snb.g, &snb.h)
			if err != nil {
				log.Println(err)
				echo.NewHTTPError(500, "could not scan")
			}
			type box struct {
				A int64 `json:"1"`
				B int64 `json:"2"`
				C int64 `json:"3"`
				D int64 `json:"4"`
				E int64 `json:"5"`
				F int64 `json:"6"`
				G int64 `json:"7"`
				H int64 `json:"8"`
			}
			var (
				vals = []box{
					{snb.a, snb.b, snb.c, snb.d, snb.e, snb.f, snb.g, snb.h},
				}
			)
			enc := json.NewEncoder(c.Response())
			for _, l := range vals {
				if err := enc.Encode(l); err != nil {
					return err
				}
				c.Response().Flush()
				debug.FreeOSMemory()
				time.Sleep(1 * time.Microsecond)
			}

		}
		defer rows.Close()
		return nil
	})
	e.Logger.Fatal(e.Start(":1323"))
}
