package core

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type HumanClient struct {
	Name       string
	ID         uuid.UUID
	Connection ConnChan
	Recv       chan any
	Arena      *Arena
	// We will respond with this index
	ResponseIndex uint
	// We will respond with this index
	EventIndex uint
}

func (client HumanClient) GetName() string {
	return client.Name
}

func (client HumanClient) GetID() uuid.UUID {
	return client.ID
}

func (client *HumanClient) SetArena(arena *Arena) {
	client.Arena = arena
}

func (client *HumanClient) GetRecv() chan<- any {
	return client.Recv
}

func (HumanClient) IsAI() bool {
	return false
}

func MakeClient(connection ConnChan) (HumanClient, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return HumanClient{}, err
	}

	client := HumanClient{
		Name:          "Unnamed User",
		ID:            uuid,
		Connection:    connection,
		Recv:          make(chan any, 32),
		Arena:         nil,
		ResponseIndex: 0,
		EventIndex:    0,
	}
	fmt.Println("Making new client", client)
	return client, nil
}

func (client HumanClient) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection)
	for {
		select {
		case send := <-client.Recv:
			var msg Message

			switch send.(type) {
			case ServerAction:
				panic("Wrong type")
			case ServerResponse:
				panic("Wrong type")
			case ServerEvent:
				msg.MessageType = EVENT
				msg.MessageIndex = client.EventIndex
				msg.Data = send
				client.EventIndex += 1
			}

			bytes, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}
			client.Connection.Send(bytes)

		case recv := <-client.Connection.RecvChan():
			if err, ok := recv.(error); ok {
				log.Println("Error: ", err)
				client.Connection.CloseConnChan()
				client.HandleClientDestruction()
				return
			}

			msg, err := ReceiveRequest(recv.([]byte),
				client.ResponseIndex,
			)

			log.Printf("Message received: %+v\n", msg)

			if err != nil {
				log.Println(err)
				continue
			}
			dispatchResult, err := ServerActionDecode(&client, msg, nil)
			if err != nil {
				log.Println("Problem with message during dispatch:", err)
				continue
			}

			var ret_msg Message
			ret_msg.MessageType = RESPONSE
			ret_msg.MessageIndex = client.ResponseIndex
			ret_msg.Data = dispatchResult

			client.ResponseIndex += 1
			bytes, err := json.Marshal(ret_msg)
			if err != nil {
				panic(err)
			}

			client.Connection.Send(bytes)
		}
	}
}

func (client *HumanClient) HandleInitialMessageAction(msg InitialMessageAction, other any) (any, error) {
	log.Println("Handling initial message action")
	client.Name = msg.Name
	return GenericResponse{
		Success:    true,
		FailReason: "",
	}, nil
}

func (client *HumanClient) HandleListArenasAction(data ListArenasAction, other any) (any, error) {
	log.Println("Handling list arenas action")

	list := ListArenas()
	return ListArenasResponse{
		Success:   true,
		ArenaList: list,
	}, nil
}

func (client *HumanClient) HandleJoinArenaAction(data JoinArenaAction, other any) (any, error) {
	if client.Arena != nil {
		return GenericResponse{
			Success:    false,
			FailReason: "Already in an arena",
		}, nil
	}

	arena, err := GetArenaFromName(data.ArenaName)
	if err != nil {
		return GenericResponse{
			Success:    false,
			FailReason: err.Error(),
		}, nil
	}

	err = arena.JoinArena(client)
	if err != nil {
		return GenericResponse{
			Success:    false,
			FailReason: err.Error(),
		}, nil
	}

	client.Arena = arena
	return GenericResponse{
		Success:    true,
		FailReason: "",
	}, nil
}

func (client *HumanClient) HandleServerArenaAction(action ServerArenaAction, other any) (any, error) {
	if client.Arena == nil {
		return GenericResponse{
			Success:    false,
			FailReason: "No arena found",
		}, nil
	}

	idx, err := client.Arena.getPlayerIdx(client)
	if err != nil {
		return GenericResponse{
			Success:    false,
			FailReason: err.Error(),
		}, nil
	}

	_, err = ArenaActionDecode(client.Arena, action.ArenaAction, idx)
	if err != nil {
		return GenericResponse{
			Success:    false,
			FailReason: err.Error(),
		}, nil
	}

	return GenericResponse{
		Success:    true,
		FailReason: "",
	}, nil
}

func (client *HumanClient) HandleCreateArenaAction(data CreateArenaAction, other any) (any, error) {
	err := CreateAndAddArena(data.ArenaName)
	if err != nil {
		return GenericResponse{
			Success:    false,
			FailReason: err.Error(),
		}, nil
	}
	return GenericResponse{
		Success:    true,
		FailReason: "",
	}, nil
}

func (client *HumanClient) HandleArenaInfoAction(data ArenaInfoAction, other any) (any, error) {
	if client.Arena == nil {
		return GenericResponse{
			Success:    false,
			FailReason: "Not in an arena",
		}, nil
	}

	return client.Arena.GetArenaInfo(), nil
}

func (client *HumanClient) HandleClientDestruction() {
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

func (client HumanClient) GetSendChannel() chan<- any {
	return client.Recv
}
