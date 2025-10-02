package core

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type ComputerClient struct {
	Name  string
	ID    uuid.UUID
	Recv  chan any
	Arena *Arena
}

type Die struct{}

func (client ComputerClient) GetName() string {
	return client.Name
}

func (client ComputerClient) GetID() uuid.UUID {
	return client.ID
}

func (client *ComputerClient) SetArena(arena *Arena) {
	client.Arena = arena
}

func (client *ComputerClient) GetRecv() chan<- any {
	return client.Recv
}

func (ComputerClient) IsAI() bool {
	return true
}

func MakeComputerClient() (ComputerClient, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ComputerClient{}, err
	}

	client := ComputerClient{
		Name:  "Computer Client " + uuid.String(),
		ID:    uuid,
		Recv:  make(chan any, 32),
		Arena: nil,
	}
	fmt.Println("Making new client", client)

	return client, nil
}

func (client *ComputerClient) Loop() {
	log.Println(client.Name, client.ID)
	for send := range client.Recv {
		log.Println("Computer Loop")

		switch send := send.(type) {
		case ServerEvent:
			ServerEventDecode(client, send, Unit)
		case Die:
			break
		}
	}
}

func (client *ComputerClient) HandleServerArenaEvent(event ServerArenaEvent, extraData UnitType) (UnitType, error) {
	// Handle arena events from the server
	log.Printf("ComputerClient %s received arena event: %+v", client.Name, event)
	return ArenaEventDecode(client, event.ArenaMessage, extraData)
}

func (client *ComputerClient) HandlePlayerJoinedEvent(event PlayerJoinedEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Player %s joined the arena", client.Name, event.AgentInfo.Name)
	return Unit, nil
}

func (client *ComputerClient) HandlePlayerQuitEvent(event PlayerQuitEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Player %s quit the arena", client.Name, event.Name)
	return Unit, nil
}

func (client *ComputerClient) HandleGameStartedEvent(event GameStartedEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Game started", client.Name)
	// TODO: Initialize AI game state
	return Unit, nil
}

func (client *ComputerClient) HandleArenaBoardEvent(event ArenaBoardEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Received board event: %+v", client.Name, event.BoardEvent)
	return BoardEventDecode(client, event.BoardEvent, extraData)
}

func (client *ComputerClient) HandlePlayerActionEvent(event PlayerActionEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Player %d performed action: %+v", client.Name, event.FromPlayer, event.Action)
	// TODO: Update AI game state based on player action
	return Unit, nil
}

func (client *ComputerClient) HandlePotentialActionEvent(event PotentialActionEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Potential actions available: %+v", client.Name, event.Actions)
	log.Println("Test!!! ", *client)

	// Immediately do a potential action
	for _, action := range event.Actions {
		idx, err := client.Arena.getPlayerIdx(client)
		if err != nil {
			panic(err)
		}
		ArenaActionDecode(client.Arena, PlayerActionData{
			Action: action,
		}, idx)
	}

	// TODO: Implement AI decision making for potential actions
	return Unit, nil
}

func (client *ComputerClient) HandleGameSetupEvent(event GameSetupEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Game setup: %+v", client.Name, event.Setup)
	// TODO: Initialize AI with game setup information (tiles, dora, etc.)
	return Unit, nil
}

func (client *ComputerClient) HandleGameEndEvent(event GameEndEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Game ended: %+v", client.Name, event.GameResult)
	// TODO: Process game end results for AI learning
	return Unit, nil
}

func (client *ComputerClient) HandleClientDestruction() {
	if client.Arena != nil {
		idx, err := client.Arena.getPlayerIdx(client)
		if err != nil {
			return
		}

		_, err = client.Arena.HandlePlayerQuitActionData(PlayerQuitActionData{}, idx)
		if err != nil {
			return
		}
	}
}
