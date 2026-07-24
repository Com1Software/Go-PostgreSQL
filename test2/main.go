package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// PostgreSQL connection
	connStr := "host=localhost port=5432 user=postgres password=password dbname=yourdb sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	// Create table if needed
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// --- Fyne GUI ---
	a := app.New()
	w := a.NewWindow("PostgreSQL GUI Example")
	w.Resize(fyne.NewSize(400, 400))

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter a name")

	statusLabel := widget.NewLabel("Connected to PostgreSQL")

	// List widget
	listWidget := widget.NewList(
		func() int {
			count := 0
			rows, err := db.Query(`SELECT id FROM users`)
			if err == nil {
				for rows.Next() {
					count++
				}
				rows.Close()
			}
			return count
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(i int, o fyne.CanvasObject) {
			rows, err := db.Query(`SELECT id, name FROM users ORDER BY id`)
			if err != nil {
				o.(*widget.Label).SetText("Error loading")
				return
			}
			defer rows.Close()

			index := 0
			for rows.Next() {
				var id int
				var name string
				rows.Scan(&id, &name)
				if index == i {
					o.(*widget.Label).SetText(fmt.Sprintf("%d: %s", id, name))
					return
				}
				index++
			}
		},
	)

	// Insert button
	insertButton := widget.NewButton("Insert Name", func() {
		name := nameEntry.Text
		if name == "" {
			statusLabel.SetText("Name cannot be empty")
			return
		}

		_, err := db.Exec(`INSERT INTO users (name) VALUES ($1)`, name)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Insert error: %v", err))
			return
		}

		statusLabel.SetText("Inserted: " + name)
		nameEntry.SetText("")
		listWidget.Refresh()
	})

	w.SetContent(
		container.NewVBox(
			statusLabel,
			nameEntry,
			insertButton,
			widget.NewSeparator(),
			listWidget,
		),
	)

	w.ShowAndRun()
}
