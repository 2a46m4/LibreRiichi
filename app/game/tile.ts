import { Tile as TileValue } from "../messaging/tile"

const TileMask = 0b11 << 4
const TileShift = 4
const ManzuBit = 0
const PinzuBit = 1
const SouzuBit = 2
const HonourBit = 3
const NumberMask = 0b1111
const SpecialMask = 0b11 << 6

export function decode(array: string): Uint8Array {
  const chars =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
  const lookup = new Map<string, number>()
  for (let i = 0; i < chars.length; i++) {
    lookup.set(chars[i], i)
  }

  // Remove padding
  const cleanInput = array.replace(/=/g, '')
  const bytes: number[] = []

  for (let i = 0; i < cleanInput.length; i += 4) {
    const a = lookup.get(cleanInput[i]) || 0
    const b = lookup.get(cleanInput[i + 1]) || 0
    const c = lookup.get(cleanInput[i + 2]) || 0
    const d = lookup.get(cleanInput[i + 3]) || 0

    const bitmap = (a << 18) | (b << 12) | (c << 6) | d

    bytes.push((bitmap >> 16) & 255)
    if (i + 2 < cleanInput.length) bytes.push((bitmap >> 8) & 255)
    if (i + 3 < cleanInput.length) bytes.push(bitmap & 255)
  }

  return new Uint8Array(bytes)
}

// We should use singleton variables that store a single reference to texture/mesh?
// Or multiple singletons store a single reference to texture/mesh
export class Tile {
  public value: number

  constructor(value: number) {
    this.value = value
  }

  static sort(a: Tile, b: Tile) {
    return a.value - b.value
  }

  static from(base64_string: string): Tile[] {
    let result: Tile[] = []
    decode(base64_string).forEach((n) => result.push(new Tile(n)))
    return result
  }

  static get_string_representation(tile_array: Tile[]) {
    let sorted = tile_array
      .map((a) => a.clear_red_or_dora())
      .sort((a, b) => (a.value > b.value ? 1 : -1))

    let str = ''

    for (let i = 0; i < sorted.length;) {
      let current_meld = sorted[i].value & TileMask
      let meld_list = [sorted[i]]
      let j = i + 1

      while (
        j < sorted.length &&
        (sorted[j].value & TileMask) == current_meld
      ) {
        meld_list.push(sorted[j])
        j++
      }

      i = j

      switch (current_meld >> TileShift) {
        case ManzuBit:
          console.log(
            meld_list.map((tile) => (tile.value & NumberMask) + 1).join(''),
          )
          str +=
            meld_list.map((tile) => (tile.value & NumberMask) + 1).join('') +
            'M'
          break
        case SouzuBit:
          str +=
            meld_list.map((tile) => (tile.value & NumberMask) + 1).join('') +
            'S'
          break
        case PinzuBit:
          str +=
            meld_list.map((tile) => (tile.value & NumberMask) + 1).join('') +
            'P'
          break
        case HonourBit:
          str += meld_list.map((tile) => tile.value - 47).join('') + 'Z'
          break
        default:
          throw new Error('Wrong state')
      }
    }

    return str
  }

  static get_base64_representation(tile_array: Tile[]) { }

  clear_red_or_dora(): Tile {
    if (this.value === TileValue.Invalid || this.value === TileValue.Hidden)
      return this
    return new Tile(this.value & ~(TileValue.DoraTile | TileValue.RedTile))
  }

  is_invalid(): boolean {
    return this.value === TileValue.Invalid
  }

  is_hidden(): boolean {
    return this.value === TileValue.Hidden
  }

  is_honour(): boolean {
    return (this.value & TileMask) >> TileShift === HonourBit
  }

  is_wind(): boolean {
    return this.value >= 48 && this.value <= 51
  }

  is_dragon(): boolean {
    return this.value >= 52 && this.value <= 54
  }

  get_tile_type(): number {
    return (this.value & TileMask) >> TileShift
  }

  // returns the number of the tile, 0-indexed
  get_tile_number(): number {
    return this.value & NumberMask
  }

  is_souzu(): boolean {
    return this.get_tile_type() === SouzuBit
  }

  is_manzu(): boolean {
    return this.get_tile_type() === ManzuBit
  }

  is_pinzu(): boolean {
    return this.get_tile_type() === PinzuBit
  }

  set_red_tile(): Tile {
    return new Tile(this.value | TileValue.RedTile)
  }

  set_dora_tile(): Tile {
    return new Tile(this.value | TileValue.DoraTile)
  }

  set_tile_number(num: number): Tile {
    return new Tile((this.value & (TileMask | SpecialMask)) | num)
  }

  get_tile_image_path(): URL {
    const base = '/app/assets/riichi-mahjong-tiles/Export/Regular/'

    const honour_map = new Map<number, string>([
      [TileValue.EastTile, 'Ton.png'],
      [TileValue.SouthTile, 'Nan.png'],
      [TileValue.WestTile, 'Shaa.png'],
      [TileValue.NorthTile, 'Pei.png'],
      [TileValue.White, 'Haku.png'],
      [TileValue.Red, 'Chun.png'],
      [TileValue.Green, 'Hatsu.png'],
    ])

    let url_construct = (name: string) =>
      new URL(
        base + name + (this.get_tile_number() + 1) + '.png',
        import.meta.url,
      )

    if (this.is_honour()) {
      return new URL(base + honour_map.get(this.value)!, import.meta.url)
    } else if (this.is_manzu()) {
      return url_construct('Man')
    } else if (this.is_pinzu()) {
      return url_construct('Pin')
    } else if (this.is_souzu()) {
      return url_construct('Sou')
    } else {
      throw new Error('Unknown tile type')
    }
  }
}

export const HiddenTile = new Tile(TileValue.Hidden)