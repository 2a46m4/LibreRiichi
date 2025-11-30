package game

// All outgoing indices should be the arena index
// Game indices should be kept internal

import (
	msg "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	gamedata "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
)

type MessageSendInfo struct {
	Events []msg.BoardEvent
	SendTo uint8
}

type AwaitAction struct {
    PotentialActions []gamedata.Action
    SentTo uint8
}

// Changes the message from game to arena index
//
// Only changes the index of the message sender and not the indices described in the event
func ChangeToArenaIdx(infos []MessageSendInfo, ordering Ordering) {
	for infoI := range infos {
		infos[infoI].SendTo = ordering.ArenaIdx(infos[infoI].SendTo)
	}
}
