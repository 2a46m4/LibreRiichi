package game

// All outgoing indices should be the arena index
// Game indices should be kept internal

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MessageSendInfo struct {
	Events []BoardEvent
	SendTo uint8
}

// Changes the message from game to arena index
func ChangeToArenaIdx(infos []MessageSendInfo, ordering Ordering) {
	for infoI, info := range infos {
		infos[infoI].SendTo = ordering.ArenaIdx(infos[infoI].SendTo)
		for eventI, event := range info.Events {
			switch event := event.(type) {
			case GameEndEvent: // TODO
			case GameSetupEvent:
				handleSetup(event.Setup)
			case PlayerActionEvent:
				event.FromPlayer = ordering.ArenaIdx(event.FromPlayer)
				info.Events[eventI] = event
			}
		}
	}
}

func handleSetup([]Setup) {
	// Noop so far
}
