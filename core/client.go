package core

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type Client struct {
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

func MakeClient(connection ConnChan) (Client, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return Client{}, err
	}

	client := Client{
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

func (client Client) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection)
	for {
		select {
		case send := <-client.Recv:
			var msg Message

			switch send.(type) {
			case ServerAction:
				panic("Wrong type")
			case ServerResponse:
				msg.MessageType = RESPONSE
				msg.MessageIndex = client.ResponseIndex
				msg.Data = send
				client.ResponseIndex += 1
			case ServerEvent:
				msg.MessageType = EVENT
				msg.MessageIndex = client.EventIndex
				msg.Data = send
				client.EventIndex += 1
			}

			bytes, err := json.Marshal(send)
			if err != nil {
				panic(err)
			}
			fmt.Println("Sending", string(bytes))
			client.Connection.Send(bytes)

		case recv := <-client.Connection.RecvChan():
			if err, ok := recv.(error); ok {
				fmt.Println("Error: ", err)
				client.Connection.CloseConnChan()
				client.HandleClientDestruction()
				return
			}

			msg, err := ReceiveRequest(recv.([]byte),
				client.ResponseIndex,
			)

			if err != nil {
				fmt.Println(err)
				continue
			}
			dispatchResult, err := ServerActionDecode(&client, msg, nil)
			if err != nil {
				fmt.Println("Problem with message during dispatch:", err)
				continue
			}

			client.GetSendChannel() <- dispatchResult
		}
	}
}

func (client *Client) HandleInitialMessageAction(msg InitialMessageAction, other any) (any, error) {
	client.Name = msg.Name
	return GenericResponse{
		Success:    true,
		FailReason: "",
	}, nil
}

func (client *Client) HandleListArenasAction(data ListArenasAction, other any) (any, error) {
	list := ListArenas()
	return ListArenasResponse{
		Success:   true,
		ArenaList: list,
	}, nil
}

func (client *Client) HandleJoinArenaAction(data JoinArenaAction, other any) (any, error) {
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

	err = arena.JoinArena(client, true)
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

func (client *Client) HandleServerArenaAction(action ServerArenaAction, other any) (any, error) {
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

func (client *Client) HandleCreateArenaAction(data CreateArenaAction, other any) (any, error) {
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

func (client *Client) HandleArenaInfoAction(data ArenaInfoAction, other any) (any, error) {
	if client.Arena == nil {
		return GenericResponse{
			Success:    false,
			FailReason: "Not in an arena",
		}, nil
	}

	return client.Arena.GetArenaInfo(), nil
}

func (client *Client) HandleClientDestruction() {
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

func (client Client) GetSendChannel() chan<- any {
	return client.Recv
}
