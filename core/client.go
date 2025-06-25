package core

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	. "codeberg.org/ijnakashiar/LibreRiichi/core/messages"
	. "codeberg.org/ijnakashiar/LibreRiichi/core/util"
)

type Client struct {
	Name       string
	ID         uuid.UUID
	Connection ConnChan
	Recv       chan Message
	Arena      *Arena
}

type DispatchResult struct {
	Message Message
	DoSend  bool
}

func NoSend() DispatchResult {
	return DispatchResult{
		Message: Message{},
		DoSend:  false,
	}
}

func FormatMessage(msgType MessageType, data any) DispatchResult {
	return DispatchResult{
		Message: Message{
			MessageType: msgType,
			Data:        data,
		},
		DoSend: true,
	}
}

func SuccessMsg() DispatchResult {
	return DispatchResult{
		Message: Message{
			MessageType: GenericResponseType,
			Data: GenericResponseData{
				Success:    true,
				FailReason: "",
			},
		},
		DoSend: true,
	}
}

func FailureMsg(err string) DispatchResult {
	return DispatchResult{
		Message: Message{
			MessageType: GenericResponseType,
			Data: GenericResponseData{
				Success:    false,
				FailReason: err,
			},
		},
		DoSend: true,
	}
}

func MakeClient(connection ConnChan) (Client, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return Client{}, err
	}

	client := Client{
		Name:       "Unnamed User",
		ID:         uuid,
		Connection: connection,
		Recv:       make(chan Message, 32),
		Arena:      nil,
	}
	fmt.Println("Making new client", client)
	return client, nil
}

func (client Client) Loop() {
	fmt.Println(client.Name, client.ID, client.Connection)
	for {
		select {
		case send := <-client.Recv:
			bytes, err := json.Marshal(send)
			if err != nil {
				panic(err)
			}
			client.Connection.Send(bytes)
		case recv := <-client.Connection.RecvChan():
			if err, ok := recv.(error); ok {
				fmt.Println(err)
				client.Connection.CloseConnChan()
				return
			}

			msg := Message{}
			err := json.Unmarshal(recv.([]byte), &msg)
			if err != nil {
				fmt.Println("Error unmarshalling:", err)
				continue
			}
			dispatchResult, err := ServerActionDispatch(&client, msg)
			if err != nil {
				fmt.Println("Problem with message during dispatch:", err)
			}

			if dispatchResult.DoSend {
				dispatchResult.Message.MessageIndex = msg.MessageIndex
				client.GetSendChannel() <- dispatchResult.Message
			}
		}
	}
}

func (client *Client) HandleListArenas(data ListArenasActionData) (DispatchResult, error) {
	// TODO
	return DispatchResult{}, nil
}

// HandleJoinArenaAction implements ServerHandler.
func (client *Client) HandleJoinArena(data JoinArenaActionData) (DispatchResult, error) {
	if client.Arena != nil {
		err := errors.New("Already in an arena")
		return FailureMsg(err.Error()), err
	}

	arena, err := GetArenaFromName(data.ArenaName)
	if err != nil {
		return FailureMsg(err.Error()), err
	}

	err = arena.JoinArena(client, true)
	if err != nil {
		return FailureMsg(err.Error()), err
	}

	client.Arena = arena
	return SuccessMsg(), nil
}

func (client *Client) HandleInitialMessage(data InitialMessageActionData) (DispatchResult, error) {
	if len(data.Name) != 0 {
		fmt.Println("Renamed user to", data.Name)
		client.Name = data.Name
	}
	return SuccessMsg(), nil	
}

func (client *Client) HandleServerArena(action ServerArenaActionData) (DispatchResult, error) {
	if client.Arena != nil {
		idx, err := client.Arena.getPlayerIdx(client)
		if err != nil {
			return DispatchResult{}, err
		}

		// TODO: Do something with this result
		ArenaActionDispatch(client.Arena, action.ArenaMessage, idx)
		return DispatchResult{}, nil
	}
	err := errors.New("No arena found")
	return FailureMsg(err.Error()), err
}

func (client *Client) HandleCreateArena(data CreateArenaActionData) (DispatchResult, error) {
	err := CreateAndAddArena(data.ArenaName)
	if err != nil {
		return FailureMsg(err.Error()), err
	}
	return SuccessMsg(), nil
}

func (client Client) GetSendChannel() chan<- Message {
	return client.Recv
}
