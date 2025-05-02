package core

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/msg"
)

func PlayerJoined(player Player) Message {
	return Message{
		MessageType: PlayerJoinedEvent,
		Data:        player,
	}
}

func GameStarted(setup []Setup) Message {
	return Message{
		MessageType: GameStartedEvent,
		Data:        setup,
	}
}
