package game

type State int

const (
	StateWaiting State = iota
	StatePlaying
	StateEnding
)

type ParticipantState int

const (
	ParticipantStatePlaying ParticipantState = iota
	ParticipantStateSpectating
)
