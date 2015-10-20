package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type TplData struct {
	Table   string
	Fields  []Field
	DbrUsed string
}

func CreateTableModel(path string, table string, db *sql.DB, verbose bool) {
	var (
		name  string
		typ   string
		null  string
		key   string
		def   string
		extra string
	)

	const data = `package dbmodels

import {{.DbrUsed}}"github.com/gocraft/dbr"

type {{.Table}} struct {
	{{range .Fields}}{{.Name}}{{/*tab*/}} {{.Type}}{{/*tab*/}} {{.Tag}}
	{{end}}
}

`

	d := TplData{}
	d.Table = strings.Title(table)
	d.DbrUsed = "_"

	var dbrUsed int

	q := fmt.Sprintf("SHOW COLUMNS FROM %s", table)
	if rows, err := db.Query(q); err == nil {
		if verbose {
			fmt.Println("\tfields:")
		}
		for rows.Next() {
			err := rows.Scan(&name, &typ, &null, &key, &def, &extra)
			if err != nil {
				fmt.Errorf("%s", err.Error())
			}
			if verbose {
				fmt.Printf("\t\t`%s` %s\n", name, typ)
			}
			if typ == "tinyint(1)" {
				if null == "YES" {
					dbrUsed++
					typ = "dbr.NullBool"
				} else {
					typ = "bool"
				}
			} else if strings.Contains(typ, "int") {
				if null == "YES" {
					dbrUsed++
					typ = "dbr.NullInt64"
				} else {
					typ = "int64"
				}
			} else if strings.Contains(typ, "float") {
				if null == "YES" {
					dbrUsed++
					typ = "dbr.NullFloat64"
				} else {
					typ = "float64"
				}
			} else {
				if null == "YES" {
					dbrUsed++
					typ = "dbr.NullString"
				} else {
					typ = "string"
				}
			}

			tag := fmt.Sprintf("`db:\"%s\"`", name)
			if verbose {
				fmt.Printf("\t\t\t => %s %s %s\n", name, typ, tag)
			}
			f := Field{strings.Title(name), typ, tag}
			d.Fields = append(d.Fields, f)
		}
	}
	if dbrUsed > 0 {
		d.DbrUsed = ""
	}
	t := template.Must(template.New("struct").Parse(data))
	file, err := os.Create(path + "/" + table + ".go")
	defer file.Close()
	if err == nil {
		if err := t.Execute(file, d); err != nil {
			fmt.Errorf("%s", err.Error())
		}
		cmd := exec.Command("go", "fmt", path+"/"+table+".go")
		err = cmd.Start()
		if err == nil {
			err = cmd.Wait()
		}
	}
}
