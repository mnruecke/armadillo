package main

import "strings"

const (
	usageTemplate = `usage: armadillo command [arguments]

				The commands are:
				{{range .}}
					{{.Name | printf "%-11s"}} {{.Short}}{{end}}

				Use "armadillo help [command]" for more information.`

	helpTemplate = `usage: revel {{.UsageLine}}
				{{.Long}}`
)

type Command struct {
	Run                    func(args []string)
	UsageLine, Short, Long string
}

var commands = []*Command{
	newCmd,
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}
