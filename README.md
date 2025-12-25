LibreRiichi is an open-source Riichi Mahjong server.

# Building the game
Install Yarn and the Go toolchain.

```sh
yarn
yarn run build-backend
yarn run build-frontend
```

The compiled binary will be in `dist-server/server`.

The front-end will be in `dist` and can be served with any web server.

## Development
Run `yarn run watch` which sets up a server that updates itself when code changes.

[air](https://github.com/air-verse/air) is used to build the server when code changes.

# Project structure

## Functionality
Players join the server with a username and can join rooms (called an arena in the codebase) in the server. Rooms host standard 4-player Riichi Mahjong games based on Tenhou/Mahjong Soul rules. Once the room is filled the games can start.

## Frontend
The frontend element is written in [VueJS](https://vuejs.org/) and [Typescript](https://www.typescriptlang.org/). I think the frontend is OK, maybe could clean up a bit but not a big issue.

### Code Structure
The frontend code is in the `app` directory. 
- `components` store reusable Vue components. These are things like pop-up bars, scoreboards, elements that are really their own thing.
  - TODO: This needs a cleanup. There's a lot of random + unused code here.
  - The most important file is `threejs_game_view.vue`. This handles game messages from the server and coordinates the rendering.
- `game` contains data structures in the game.
  - TODO: It's actually all pretty unfinished at the moment. 
  - `setup.ts` contains the data types that will be sent to the client when the game starts.
  - `tile.ts` contains the tile data type, and extra functions to deserialize off the "wire".
- `messaging` contains the data structures which the server and client use to communicate with each other. 
  - Some files are auto-generated from `generate_messages.go`. These shouldn't be edited. 
- `render` contains everything needed for the rendering.
  - TODO: A lot of functionality is still missing
- `views` are pages, essentially. The entrypoint of the frontend is in login.vue.

### Game loop
The game loop is done somewhat but I haven't tried making it work again yet. Originally I had an explicit state machine that receives events from the server and user and made explicit transitions to the state machine data structure (e.g. transitioning from waiting for a player to discard to a player discarding and waiting for any naki calls). But this was WAY too much plumbing and kind of buggy. Also the code is split up into a lot of small functions this way which makes it hard to reason about (like a handler function for each state)

Now it's essentially a big async/await function that awaits for messages from either user or server and then does stuff based on the message received. The code path is a lot more linear that way (you can see what happens to the message right after it gets received, and all the logic is contained in one section). 

## Backend
The backend is written in Go. All code is stored in `core`. The folder structure is pretty crappy, I should really merge `game` and `game_data`. 

When a client connects, they are represented as a `HumanClient`. They use Go channels wrapped over a websocket to communicate to the actual client. The code in `client.go` handles messages generally and also handles joining/leaving/searching for arenas. It forwards any arena-specific code to the arena struct.

The arena struct exists during the lifetime of a single game. I'm still kind of thinking this section through, it's really messy. The arena's purpose is to coordinate between different clients and the game object itself. The arena gets called when clients send a message to the server. The handle methods are functions that get called when the arena receives that particular message. 

calls into the game object Originally I had a goroutine (basically lightweight threads). Then I switched to an explicit state machine, and this was super messy and buggy, so I switched to a goroutine model again. 

### game/game_data
stores the actual game data and game loop. 
- meld_finder has functions that help find melds (a triplet sequence of tiles).
- action_testing tests whether a particular action from a client is legal. 


## Message structure {#message-structure}
Messages are sent via Websockets. All messages contain a field to distinguish the type of message and a field for the actual data itself. 

TODO: We'd want to eventually research and use a library for messaging. It's pretty janky writing custom serialization/deserialization code. Something like Protobuf would probably work pretty well. Alternatively can think of a different structure, like a flat message with just a type and a payload. Or just use Rust with serde or C++26 (lol)

### Message hierarchy
Messages have three types:
- ServerRequest: A request sent from the client to the server. Example: client requesting an action
- ServerResponse: A reply from the server caused by a ServerRequest. Example: Server responding to an action that is not allowed
- ServerEvent: An event that occurs from the server. Example: the server processing a successful action, and broadcasting the action to everyone in the room.

### Generated message files
There are some generated message files for data type definitions and boilerplate code for JSON marshalling and unmarshalling. The message type is only encoded when marshalling and decoded when unmarshalling the message from the websocket. 

Additionally the generator generates some code that helps the messages look a bit like sum types. The generator creates the `*Handler` interface which is a function that handles all cases (like pattern matching) and a `*Decoder` class which dispatches to the correct case. So anything that implements the Handler interface will be able to use the Decoder function to dispatch to the right method. You can see this in arena.go.

Generated messages have a `_generated` suffix in their file name. They shouldn't be edited as they will be overwritten when generated again.

# TODOS
- Separate the different types of kans in messages
