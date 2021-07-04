package main

import (
	"fmt"

	"github.com/vektah/gqlparser/v2"
	gqlast "github.com/vektah/gqlparser/v2/ast"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var schema = gqlparser.MustLoadSchema(
	&gqlast.Source{
		Name: "test.graphql",
		Input: `
		type User {
			id: ID!
			name: String!
			photo: Picture
		}
		
		type Picture {
			id: ID!
			url: String!
		}
		`,
	},
)

func main() {

	for _, t := range schema.Types {
		fmt.Println(t.Name)
	}

	// cli := cli.NewCli()
	// qm := &ast.File{}

	// for _, file := range cli.Files {
	// 	res, err := ioutil.ReadFile(file)
	// 	check(err)

	// 	err = parser.ParseString("", string(res), qm)
	// 	check(err)

	// 	toPrint, err := json.MarshalIndent(qm, "", "    ")
	// 	check(err)
	// 	fmt.Println(string(toPrint))
	// 	// repr.Println(qm, repr.Indent("  "), repr.OmitEmpty(true))
	// }

	// 	err = parser.ParseString("", `
	// package main

	// component Main {
	// 	prop name string = "test"
	// }
	// `, qm)

	// 	if err != nil {
	// 		panic(err)
	// 	}

}
