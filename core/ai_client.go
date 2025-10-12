package core

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/game_data"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type ComputerClient struct {
	Name string
	ID   uuid.UUID
	Recv chan any

	Arena *Arena

	// AI state
	PlayerIndex    uint8
	InitialTiles   []Tile
	Dora           Tile
	PlayerOrder    []uint8
	RoundWind      Wind
	StartingPoints [4]uint32
	CurrentHand    []Tile
	KnownDiscards  []Tile
	KnownMelds     map[uint8][]Tile // Track other players' exposed melds
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
		Name:       "Computer Client " + uuid.String(),
		ID:         uuid,
		Recv:       make(chan any, 32),
		Arena:      nil,
		KnownMelds: make(map[uint8][]Tile),
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
	isOurMove := client.ourMove(event.FromPlayer)

	switch action := event.Action.(type) {
	case Toss:
		if isOurMove {
			client.KnownDiscards = append(client.KnownDiscards, action.TileToToss)
			client.removeFromHand(action.TileToToss)
		}

	case Draw:
		if isOurMove {
			client.CurrentHand = append(client.CurrentHand, action.DrawnTile)
		}

	case Chii:
		if !isOurMove {
			client.KnownMelds[event.FromPlayer] = append(client.KnownMelds[event.FromPlayer],
				action.TileToChii, action.TilesInHand[0], action.TilesInHand[1])
		} else {
			client.removeFromHand(action.TilesInHand[0])
			client.removeFromHand(action.TilesInHand[1])
		}
	case Pon:
		if !isOurMove {
			client.KnownMelds[event.FromPlayer] = append(client.KnownMelds[event.FromPlayer],
				action.TileToPon, action.TileToPon, action.TileToPon)
		} else {
			client.removeFromHand(action.TileToPon)
			client.removeFromHand(action.TileToPon)
		}

	case Kan:
		if !isOurMove {
			client.KnownMelds[event.FromPlayer] = append(client.KnownMelds[event.FromPlayer],
				action.TileToKan, action.TileToKan, action.TileToKan, action.TileToKan)
		} else {
			for i := 0; i < 4; i++ {
				client.removeFromHand(action.TileToKan)
			}
		}

	case Riichi:
		if isOurMove {
			client.removeFromHand(action.TileToRiichi)
		}

	case Ron:
	case Tsumo:
	}

	return Unit, nil
}

func (client *ComputerClient) HandlePotentialActionEvent(event PotentialActionEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Potential actions available: %+v", client.Name, event.Actions)

	// Immediately do a potential action
	for _, action := range event.Actions {
		idx, err := client.Arena.getPlayerIdx(client)
		if err != nil {
			panic(err)
		}
		switch action.(type) {
		case Toss:
			_, err = ArenaActionDecode(client.Arena, PlayerActionData{
				Action: Toss{
					TileToToss: Last(client.InitialTiles),
				},
			}, idx)
			if err != nil {
				panic(err)
			}
		default:
			_, err = ArenaActionDecode(client.Arena, PlayerActionData{
				Action: action,
			}, idx)
			if err != nil {
				panic(err)
			}
		}

	}

	return Unit, nil
}

func (client *ComputerClient) HandleGameSetupEvent(event GameSetupEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Game setup: %+v", client.Name, event.Setup)

	// Initialize AI with game setup information
	for _, setup := range event.Setup {
		switch setup.Type {
		case INITIAL_TILES:
			if tiles, ok := setup.Data.([]Tile); ok {
				client.InitialTiles = tiles
				client.CurrentHand = make([]Tile, len(tiles))
				copy(client.CurrentHand, tiles)
				log.Printf("ComputerClient %s: Received %d initial tiles", client.Name, len(tiles))
			}
		case DORA:
			if dora, ok := setup.Data.(Tile); ok {
				client.Dora = dora
				log.Printf("ComputerClient %s: Dora tile: %v", client.Name, dora)
			}
		case PLAYER_NUMBER:
			if playerNum, ok := setup.Data.(uint8); ok {
				client.PlayerIndex = playerNum
				log.Printf("ComputerClient %s: Player index: %d", client.Name, playerNum)
			}
		case PLAYER_ORDER:
			if order, ok := setup.Data.([]uint8); ok {
				client.PlayerOrder = order
				log.Printf("ComputerClient %s: Player order: %v", client.Name, order)
			}
		case ROUND_WIND:
			if wind, ok := setup.Data.(Wind); ok {
				client.RoundWind = wind
				log.Printf("ComputerClient %s: Round wind: %v", client.Name, wind)
			}
		case STARTING_POINTS:
			if points, ok := setup.Data.([4]uint32); ok {
				client.StartingPoints = points
				log.Printf("ComputerClient %s: Starting points: %v", client.Name, points)
			}
		case ROUND_NUMBER:
			// Handle round number if needed
			log.Printf("ComputerClient %s: Round number setup received", client.Name)
		}
	}

	return Unit, nil
}

func (client *ComputerClient) HandleGameEndEvent(event GameEndEvent, extraData UnitType) (UnitType, error) {
	log.Printf("ComputerClient %s: Game ended: %+v", client.Name, event.GameResult)
	return Unit, nil
}

func (client *ComputerClient) ourMove(fromPlayer uint8) bool {
	return client.PlayerIndex == fromPlayer
}

// Helper method to remove a tile from the AI's tracked hand
func (client *ComputerClient) removeFromHand(tile Tile) {
	for i, handTile := range client.CurrentHand {
		if handTile == tile {
			Remove(&client.CurrentHand, i)
			return
		}
	}
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
