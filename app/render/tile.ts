import {Tile, TileValue} from "../game/tile";
import {load_texture} from "./texture";
import * as THREE from "three";

const map = new Map<number, string>([
    // Manzu (Characters) - 1-9
    [TileValue.Manzu, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man1.png', import.meta.url).href],
    [TileValue.Manzu + 1, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man2.png', import.meta.url).href],
    [TileValue.Manzu + 2, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man3.png', import.meta.url).href],
    [TileValue.Manzu + 3, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man4.png', import.meta.url).href],
    [TileValue.Manzu + 4, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man5.png', import.meta.url).href],
    [TileValue.Manzu + 5, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man6.png', import.meta.url).href],
    [TileValue.Manzu + 6, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man7.png', import.meta.url).href],
    [TileValue.Manzu + 7, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man8.png', import.meta.url).href],
    [TileValue.Manzu + 8, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Man9.png', import.meta.url).href],

    // Pinzu (Circles) - 1-9
    [TileValue.Pinzu, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin1.png', import.meta.url).href],
    [TileValue.Pinzu + 1, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin2.png', import.meta.url).href],
    [TileValue.Pinzu + 2, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin3.png', import.meta.url).href],
    [TileValue.Pinzu + 3, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin4.png', import.meta.url).href],
    [TileValue.Pinzu + 4, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin5.png', import.meta.url).href],
    [TileValue.Pinzu + 5, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin6.png', import.meta.url).href],
    [TileValue.Pinzu + 6, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin7.png', import.meta.url).href],
    [TileValue.Pinzu + 7, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin8.png', import.meta.url).href],
    [TileValue.Pinzu + 8, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pin9.png', import.meta.url).href],

    // Souzu (Bamboo) - 1-9
    [TileValue.Souzu, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou1.png', import.meta.url).href],
    [TileValue.Souzu + 1, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou2.png', import.meta.url).href],
    [TileValue.Souzu + 2, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou3.png', import.meta.url).href],
    [TileValue.Souzu + 3, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou4.png', import.meta.url).href],
    [TileValue.Souzu + 4, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou5.png', import.meta.url).href],
    [TileValue.Souzu + 5, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou6.png', import.meta.url).href],
    [TileValue.Souzu + 6, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou7.png', import.meta.url).href],
    [TileValue.Souzu + 7, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou8.png', import.meta.url).href],
    [TileValue.Souzu + 8, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Sou9.png', import.meta.url).href],

    // Wind tiles
    [TileValue.EastTile, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Ton.png', import.meta.url).href],
    [TileValue.SouthTile, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Nan.png', import.meta.url).href],
    [TileValue.WestTile, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Shaa.png', import.meta.url).href],
    [TileValue.NorthTile, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Pei.png', import.meta.url).href],

    // Dragon tiles
    [TileValue.White, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Haku.png', import.meta.url).href],
    [TileValue.Red, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Chun.png', import.meta.url).href],
    [TileValue.Green, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Hatsu.png', import.meta.url).href],

    [TileValue.Invalid, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Blank.png', import.meta.url).href],
    [TileValue.Hidden, new URL('/app/assets/riichi-mahjong-tiles/Export/Regular/Back.png', import.meta.url).href],
])

export function get_image_texture(tile: Tile) {
    let result = map.get(tile.clear_red_or_dora().value)
    if (result === undefined) {
        throw new Error('Unknown tile type')
    }
    return load_texture(result)
}

export function load_all_textures() {
    let texture_map = new Map<number, THREE.Texture>()
    for (let [value, url] of map.entries()) {
        texture_map.set(value, load_texture(url))
    }
    return texture_map
}