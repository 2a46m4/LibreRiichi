export enum SetupType {
  INITIAL_TILES = 0,
  DORA = 1,
  STARTING_POINTS = 2,
  PLAYER_NUMBER = 3,
  PLAYER_ORDER = 4,
  ROUND_WIND = 5,
  ROUND_NUMBER = 6,
}

export type Setup =
  | InitialTiles
  | Dora
  | StartingPoints
  | PlayerNumber
  | PlayerOrder
  | RoundWind
  | RoundNumber

export type TileArray = string

export type InitialTiles = {
  setup_type: SetupType.INITIAL_TILES
  data: TileArray
}

export type Dora = {
  setup_type: SetupType.DORA
  data: number
}

export type StartingPoints = {
  setup_type: SetupType.STARTING_POINTS
  data: number[]
}

export type PlayerNumber = {
  setup_type: SetupType.PLAYER_NUMBER
  data: number
}

export type PlayerOrder = {
  setup_type: SetupType.PLAYER_ORDER
  data: TileArray
}

export type RoundWind = {
  setup_type: SetupType.ROUND_WIND
  data: number
}

export type RoundNumber = {
  setup_type: SetupType.ROUND_NUMBER
  data: number // int
}
