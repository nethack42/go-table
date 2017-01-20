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
		Prefix:  "\t",
		Rows: []table.Row{
			table.Row{
				table.Column{
					Value: "Cyan",
					Color: color.FgCyan,
				},
				table.Column{
					Value: "This is a bit longer",
				},
				table.Column{
					Value: "I'm bold",
					Bold:  true,
				},
			},
			table.Row{
				table.Column{
					Value:  "col1",
					Italic: true,
				},
				table.Column{
					Value:      "right-aligned",
					RightAlign: true,
					Italic:     true,
				},
				table.Column{
					Value:     "Error",
					Color:     table.Error,
					Underline: true,
				},
			},
			table.Row{
				table.Column{
					Value: "This is a rather long line",
				},
				table.Column{
					Value:      "Warning",
					Color:      table.Warn,
					RightAlign: true,
				},
				table.Column{
					Value:      "This is special",
					Attributes: []color.Attribute{color.FgBlack, color.BgWhite, color.BlinkSlow},
				},
			},
		},
	}

	if err := table.Write(os.Stdout); err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("test")
}
