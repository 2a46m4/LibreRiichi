package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/msg"
)

func PlayerJoined(player Agent) ArenaMessage {
	return ArenaMessage{
		MessageType: PlayerJoinedEventType,
		Data:        player,
	}
}

func GameStarted(setup []Setup) ArenaMessage {
	return ArenaMessage{
		MessageType: GameStartedEventType,
		Data:        setup,
	}
}
