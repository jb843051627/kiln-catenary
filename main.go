package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jb843051627/kiln-catenary/internal/handler"
	"github.com/jb843051627/kiln-catenary/internal/model"
	"github.com/jb843051627/kiln-catenary/internal/service"
	"github.com/jb843051627/kiln-catenary/internal/store"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func main() {
	smoke := flag.Bool("smoke-test", false, "run startup probe")
	flag.Parse()
	path := os.Getenv("KILN_CATENARY_DB")
	if path == "" {
		path = "data/kiln-catenary.db"
	}
	db, err := store.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	app := service.NewApp(db)
	defer app.Close()
	defer db.Close()
	if *smoke {
		if err := smokeTest(app); err != nil {
			log.Fatal(err)
		}
		return
	}
	addr := os.Getenv("KILN_CATENARY_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("kiln-catenary listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.New(app)))
}

func smokeTest(app *service.App) error {
	server := httptest.NewServer(handler.New(app))
	defer server.Close()
	response, err := http.Get(server.URL + "/healthz")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("health returned %d", response.StatusCode)
	}
	kiln, err := app.CreateKiln(context.Background(), model.Kiln{Code: "SMOKE", Cell: "north", MaxTemperature: 1400, MaxPressure: 2, Atmosphere: "neutral", Active: true})
	if err != nil {
		return err
	}
	if _, err := app.GetKiln(context.Background(), kiln.ID); err != nil {
		return err
	}
	fmt.Println("smoke ok")
	return nil
}
