package command

import "flag"

type AddCommand struct {
	name    string
	handler func()
	flags   *flag.FlagSet
}

func (cmd *AddCommand) Name() string {
	return cmd.Name()
}
