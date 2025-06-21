package core

type BadMessage struct{}

func (BadMessage) Error() string {
	return "Bad message"
}
