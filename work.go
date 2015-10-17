package main

import (
	"database/sql"
	"fmt"
)

func DoWork(name string, verbose bool) {
	
	p := GetProject(name)

	db, err := sql.Open("mysql", 
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", p.DBUser, p.DBPass, p.DBHost, p.DBPort, p.DBName))
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer db.Close()

	var tabl string
	var ty string
	if rows, err := db.Query("show full tables where Table_Type != 'VIEW'"); err == nil {
		for rows.Next() {
			rows.Scan(&tabl, &ty)
			if verbose {
				fmt.Printf("\ngenerate `%s` table\n", tabl)
			}
			CreateTableModel(p.Path, tabl, db, verbose)
		}
	}
}
