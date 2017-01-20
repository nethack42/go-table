package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/nethack42/go-table"
)

func main() {
	table := table.Table{
		Spacing: 5,
		Rows: []table.Row{
			table.Row{
				table.Column{
					Value: "col1",
					Color: color.FgWhite,
				},
				table.Column{
					Value: "col2adfasdfasdf",
				},
				table.Column{
					Value: "col3",
					Bold:  true,
				},
			},
			table.Row{
				table.Column{
					Value:  "col1",
					Italic: true,
				},
				table.Column{
					Value:      "col2",
					Color:      table.Warn,
					RightAlign: true,
				},
				table.Column{
					Value:     "col3",
					Color:     table.Error,
					Underline: true,
				},
			},
		},
	}

	if err := table.Write(os.Stdout); err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("test")
}
