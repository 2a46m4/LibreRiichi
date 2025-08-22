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
	// We will respond with this index
	ResponseIndex uint
	// We expect their next message to have this index
	RequestIndex uint
	// We will respond with this index
	EventIndex uint
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
				client.RequestIndex,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dispatchResult, err := ServerActionDecode(&client, msg, nil)
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

func (client *Client) HandleInitialMessageAction(InitialMessageAction, any) (Server, error) {
	return Message{}
}
func (client *Client) HandleListArenas(data ListArenasActionData) (DispatchResult, error) {
	list := ListArenas()
	return DispatchResult{
		Message: Message{
			MessageType: ListArenasResponseType,
			Data: ListArenasResponseData{
				Success:   true,
				ArenaList: list,
			},
		},
		DoSend: true,
	}, nil
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
			return FailureMsg(err.Error()), err
		}

		err = ArenaActionDispatch(client.Arena, action.ArenaMessage, idx)
		if err != nil {
			return FailureMsg(err.Error()), err
		} else {
			return SuccessMsg(), nil
		}
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

func (client *Client) HandleGetArenaInfo(data ArenaInfoActionData) (DispatchResult, error) {
	if client.Arena == nil {
		return FailureMsg("Not in arena"), nil
	}

	return DispatchResult{
		Message: Message{
			MessageType: ArenaInfoResponseType,
			Data:        client.Arena.GetArenaInfo(),
		},
		DoSend: true,
	}, nil
}

func (client *Client) HandleClientDestruction() {
	if client.Arena != nil {
		idx, err := client.Arena.getPlayerIdx(client)
		if err != nil {
			return
		}

		err = client.Arena.HandlePlayerQuitAction(PlayerQuitActionData{}, idx)
		if err != nil {
			return
		}
	}
}

func (client Client) GetSendChannel() chan<- Message {
	return client.Recv
}
