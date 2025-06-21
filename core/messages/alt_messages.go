package core

type AltMessageHandler struct{}

// HandleGameEndEventType implements BoardEventHandler.
func (a AltMessageHandler) HandleGameEndEventType(GameEndEventData) error {
	panic("unimplemented")
}

// HandleGameSetupEventType implements BoardEventHandler.
func (a AltMessageHandler) HandleGameSetupEventType(GameSetupEventData) error {
	panic("unimplemented")
}

// HandlePlayerActionEventType implements BoardEventHandler.
func (a AltMessageHandler) HandlePlayerActionEventType(PlayerActionEventData) error {
	panic("unimplemented")
}

// HandlePotentialActionEventType implements BoardEventHandler.
func (a AltMessageHandler) HandlePotentialActionEventType(PotentialActionEventData) error {
	panic("unimplemented")
}
