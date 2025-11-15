package game

// All outgoing indices should be the arena index
// Game indices should be kept internal

import (
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
)

type MessageSendInfo struct {
	Events []BoardEvent
	SendTo uint8
}
