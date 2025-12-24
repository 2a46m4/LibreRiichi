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

### Message hierarchy



### Generated messages
Generated messages have a `_generated` suffix in their file name. They shouldn't be edited as they will be overwritten when generated again.

The purpose of generating messages is to allow Go to look a little more like a language with sum types. To take an example, `core/game_data/action.go` contains the list of actions available to the user, which each contain some data. In Go we can't do the usual functional pattern of sum types + pattern matching. So 

