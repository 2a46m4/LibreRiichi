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

export function is_red(i: number): boolean {
    return (i & TileValue.RedTile) !== 0
}

export function is_dora(i: number): boolean {
    return (i & TileValue.DoraTile) !== 0
}

// Not good
export function decode_tile (i: number): Tile {
    try {
        return NumberTile.from_number(i)
    }
    catch (error) {}
    try {
        return WindTile.from_number(i)
    }
    catch (error) {}
    return DragonTile.from_number(i)
}

export abstract class Tile {
    abstract red: boolean
    abstract dora: boolean
    abstract get_text_representation(): string

    dora_mask(): number {
        if (this.dora) {
            return TileValue.DoraTile
        } else {
            return 0
        }
    }

    red_mask() {
        if (this.red) {
            return TileValue.RedTile
        } else {
            return 0
        }
    }
}

export class NumberTile extends Tile {
    base: TileValue.Manzu | TileValue.Pinzu | TileValue.Souzu
    number: number
    red: boolean
    dora: boolean

    constructor(base: TileValue.Manzu | TileValue.Pinzu | TileValue.Souzu,
                number: number,
                red: boolean,
                dora: boolean) {
        super();
        this.base = base
        this.number = number
        this.red = red
        this.dora = dora
    }

    static from_number(i: number): NumberTile {
        if (i >= TileValue.Manzu && i < TileValue.Manzu + 9) {
            return new NumberTile(TileValue.Manzu, i - TileValue.Manzu, is_red(i), is_dora(i))
        }
        if (i >= TileValue.Pinzu && i < TileValue.Pinzu + 9) {
            return new NumberTile(TileValue.Pinzu, i - TileValue.Pinzu, is_red(i), is_dora(i))
        }
        if (i >= TileValue.Souzu && i < TileValue.Souzu + 9) {
            return new NumberTile(TileValue.Souzu, i - TileValue.Souzu, is_red(i), is_dora(i))
        }
        throw new Error("Invalid tile")
    }

    get_text_representation(): string {
        switch (this.base) {
            case TileValue.Manzu: return this.number + "M"
            case TileValue.Pinzu: return this.number + "P"
            case TileValue.Souzu: return this.number + "S"
        }
    }

    get tile_value(): TileValue {
        return (this.base + this.number) | this.dora_mask() | this.red_mask()
    }

    set tile_value(value: number) {
       Object.assign(this, NumberTile.from_number(value as number))
    }
}

export class WindTile extends Tile {
    tile_value: TileValue.EastTile | TileValue.SouthTile | TileValue.WestTile | TileValue.NorthTile
    dora: boolean;
    red: boolean;

    constructor(tile_value: TileValue.EastTile | TileValue.SouthTile | TileValue.WestTile | TileValue.NorthTile,
                red: boolean,
                dora: boolean) {

        super();
        this.tile_value = tile_value
        this.dora = dora
        this.red = red
    }

    static from_number(i: number): WindTile {
        if (i >= TileValue.EastTile && i < TileValue.EastTile + 4) {
            return new WindTile(i, is_red(i), is_dora(i))
        }
        throw new Error("Invalid tile")
    }

    get_text_representation(): string {
        return (this.tile_value - TileValue.EastTile) + "Z";
    }
}

export class DragonTile extends Tile {
    tile_value: TileValue.Sangenpai
    dora: boolean;
    red: boolean;
    constructor(tile_value: TileValue.Sangenpai, red: boolean, dora: boolean) {
        super();
        this.tile_value = tile_value
        this.dora = dora
        this.red = red
    }

    static from_number(i: number): DragonTile {
        if (i >= TileValue.Sangenpai && i < TileValue.Sangenpai + 4) {
            return new DragonTile(i, is_red(i), is_dora(i))
        }
        throw new Error("Invalid tile")
    }

    get_text_representation(): string {
        return (this.tile_value - TileValue.Sangenpai) + "Z";
    }

}