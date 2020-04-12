package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"text/template"

	"github.com/lib/pq"
)

func main2() {
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}

	var strs []string
	strs = append(strs, "test1")
	strs = append(strs, "test2")

	tmpl, err := template.New("test").
		Funcs(funcMap).
		Parse(`
{{range $index, $element := .}}
  Number: {{inc $index}}, Text:{{$element}}
{{end}}`)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, strs)
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func main() {
	qry := `
	SELECT id, name, age
	FROM(
	SELECT 1 AS id, 'Adam Gill' AS name, 15 AS age
	UNION ALL SELECT 2 AS id, 'Ray Cannon' AS name, 17 AS age
	UNION ALL SELECT 3 AS id, 'Birdie Carr' AS name, 16 AS age
	UNION ALL SELECT 4 AS id, 'Betty Bryan' AS name, 16 AS age
	UNION ALL SELECT 5 AS id, 'Lizzie Hamilton' AS name, 18 AS age
	UNION ALL SELECT 6 AS id, 'Zachary Zimmerman' AS name, 19 AS age
	UNION ALL SELECT 7 AS id, 'Shawn Myers' AS name, 18 AS age
	UNION ALL SELECT 8 AS id, 'Amanda Bowman' AS name, 18 AS age
	UNION ALL SELECT 9 AS id, 'Bowman Gutierrez' AS name, 18 AS age
	) students
	{{if .Aux.ByID}}
	WHERE students.id = {{param "id"}}
	{{else}}
	WHERE students.age = ANY({{param "ages"}})
	AND students.name LIKE {{param "name"}}
	{{end}}
	`
	parsed, params := NewQuery(qry, QryArgs{
		Aux:    map[string]interface{}{"ByID": true},
		Params: map[string]interface{}{"id": 5},
	})
	execute(parsed, params)
	parsed, params = NewQuery(qry, QryArgs{
		Aux:    map[string]interface{}{"ByID": false},
		Params: map[string]interface{}{"ages": pq.Array([]int{16, 18}), "name": "B%"},
	})
	execute(parsed, params)
}

func execute(qry string, params []interface{}) {
	fmt.Printf("Qry:%s\nParams:%v\n", qry, params)
	if true {
		return
	}
	rows, err := db.Query(qry, params...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var id, age int
	var name string
	for rows.Next() {
		if err = rows.Scan(&id, &name, &age); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Student (%d) %s is %d years old\n", id, name, age)
	}
}

type QryArgs struct {
	Aux    map[string]interface{}
	Params map[string]interface{}
}

func NewQuery(qry string, args QryArgs) (string, []interface{}) {
	params := make([]interface{}, 0, len(args.Params))
	i := 0
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"param": func(param string) string {
			params = append(params, args.Params[param])
			i++
			return "$" + strconv.Itoa(i)
		},
	}

	tmpl, err := template.New(fmt.Sprint(rand.Float64())).
		Funcs(funcMap).
		Parse(qry)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, args)
	if err != nil {
		panic(err)
	}
	return tpl.String(), params
}
