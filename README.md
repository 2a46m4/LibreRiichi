LibreRiichi is an open-source Riichi Mahjong server.

# Building the game
Install Yarn and Go toolchain.

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

## Frontend
The frontend element is written in [VueJS](https://vuejs.org/) and [Typescript](https://www.typescriptlang.org/).

### Functionality
Players join the server with a username and can join rooms in the server. Rooms host standard 4-player Riichi Mahjong games based on Tenhou/Mahjong Soul rules. Once the room is filled the games can start.

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
  - 


## Message structure
Messages are sent via Websockets. All messages contain a field to distinguish the type of message and a field for the actual data itself. 

TODO: We'd want to eventually research and use a library for messaging. It's pretty janky writing custom serialization/deserialization code. It's also why I keep thinking about moving to a different language (Rust is more natural for writing code like this, ). Alternatively can think of a different structure. 

### Message hierarchy
Messages have three types:
- ServerRequest: A request sent from the client to the server. Example: client requesting an action
- ServerResponse: A reply from the server caused by a ServerRequest. Example: Server responding to an action that is not allowed
- ServerEvent: An event that occurs from the server. Example: the server processing a successful action, and broadcasting the action to everyone in the room.


### Generated message files
There are some generated message files for data type definitions and boilerplate code for JSON marshalling and unmarshalling. The message types that are used in the backend codebase don't include the type of the message itself to make it less annoying to write. It's only encoded when marshalling and decoded when unmarshalling the message from the websocket.

Additionally the messages look a bit like sum types with the generation. The generator creates the `*Handler` interface which is a function that handles all case (like pattern matching)

Generated messages have a `_generated` suffix in their file name. They shouldn't be edited as they will be overwritten when generated again.


