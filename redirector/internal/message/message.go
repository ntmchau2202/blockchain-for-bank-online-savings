package message

type Request struct {
	Cmd     Command
	Details map[string]string // everything going to be sent to blockchain will be string to avoid number overflow & error
}

type Response struct {
	Stat    Status
	Details map[string]interface{}
}

func (r Request) ValidateCommand(cmd Command) bool {
	return string(r.Cmd) == cmd.ToString()
}
