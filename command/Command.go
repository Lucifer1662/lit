package command

type Command interface {
	Name() string
	Invoke(args []string)
}
