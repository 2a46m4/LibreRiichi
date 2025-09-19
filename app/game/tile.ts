export enum TileValue {
    Manzu = 0,
    Pinzu = 16,
    Souzu = 32,

    Kazehai    = 48,
    EastTile   = 48,
    SouthTile  = 49,
    WestTile   = 50,
    NorthTile  = 51,

    Sangenpai = 52,
    White     = 52,
    Red       = 53,
    Green     = 54,

    DoraTile = 64,
    RedTile  = 128,
    Hidden   = 254,
    Invalid  = 255,
}

const TileMask = 0b11 << 4;
const TileShift = 4;
const ManzuBit = 0;
const PinzuBit = 1;
const SouzuBit = 2;
const HonourBit = 3;
const NumberMask = 0b1111;
const SpecialMask = 0b11 << 6;

export function decode(array: string): Uint8Array {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
    const lookup = new Map<string, number>();
    for (let i = 0; i < chars.length; i++) {
        lookup.set(chars[i], i);
    }

    // Remove padding
    const cleanInput = array.replace(/=/g, '');
    const bytes: number[] = [];

    for (let i = 0; i < cleanInput.length; i += 4) {
        const a = lookup.get(cleanInput[i]) || 0;
        const b = lookup.get(cleanInput[i + 1]) || 0;
        const c = lookup.get(cleanInput[i + 2]) || 0;
        const d = lookup.get(cleanInput[i + 3]) || 0;

        const bitmap = (a << 18) | (b << 12) | (c << 6) | d;

        bytes.push((bitmap >> 16) & 255);
        if (i + 2 < cleanInput.length) bytes.push((bitmap >> 8) & 255);
        if (i + 3 < cleanInput.length) bytes.push(bitmap & 255);
    }

    return new Uint8Array(bytes);
}

export class Tile {
    public value: number

    constructor(value: number) {
        this.value = value
    }

	static from(base64_string: string): Tile[] {
		let result: Tile[] = []
		decode(base64_string).forEach((n)=>result.push(new Tile(n)))
		return result
	}

    static get_string_representation(tile_array: Tile[]) {
		let sorted = tile_array
			.map((a)=>a.clear_red_or_dora())
			.sort((a, b)=> (a.value > b.value)? 1 : -1)

		let str = ""

		for (let i = 0; i < sorted.length; ) {
			let current_meld = sorted[i].value & TileMask
			let meld_list = [sorted[i]]
			let j = i + 1

			while(j < sorted.length && (sorted[j].value & TileMask) == current_meld) {
				meld_list.push(sorted[j])
				j++
			}

			i = j

			switch(current_meld >> TileShift) {
				case ManzuBit:
					console.log(meld_list
							.map((tile) => (tile.value & NumberMask) + 1)
							.join(""))
					str += meld_list
							.map((tile) => (tile.value & NumberMask) + 1)
							.join("") + "M"
					break
				case SouzuBit:
					str += meld_list
							.map((tile) => (tile.value & NumberMask) + 1)
							.join("") + "S"
					break
				case PinzuBit:
					str += meld_list
							.map((tile) => (tile.value & NumberMask) + 1)
							.join("") + "P"
					break
				case HonourBit:
					str += meld_list
							.map((tile) => tile.value - 47)
							.join("") + "Z"
					break
				default:
					throw new Error("Wrong state")
			}
		}

		return str
    }

    static get_base64_representation(tile_array: Tile[]) {

    }

    clear_red_or_dora(): Tile {
        return new Tile(this.value & ~(TileValue.DoraTile | TileValue.RedTile));
    }

    is_invalid(): boolean {
        return this.value === TileValue.Invalid;
    }

    is_hidden(): boolean {
        return this.value === TileValue.Hidden;
    }

    is_honour(): boolean {
        return (this.value & TileMask) === HonourBit;
    }

    is_wind(): boolean {
        return this.value >= 48 && this.value <= 51;
    }

    is_dragon(): boolean {
        return this.value >= 52 && this.value <= 54;
    }

	is_manzu(): boolean {
        return (this.value & TileMask) === ManzuBit;
    }

    is_pinzu(): boolean {
        return (this.value & TileMask) === PinzuBit;
    }

    get_tile_number(): number {
        return this.value & NumberMask;
    }

    set_red_tile(): Tile {
        return new Tile(this.value | TileValue.RedTile);
    }

    set_dora_tile(): Tile {
        return new Tile(this.value | TileValue.DoraTile);
    }

    set_tile_number(num: number): Tile {
        return new Tile((this.value & (TileMask | SpecialMask)) | num);
    }
}
