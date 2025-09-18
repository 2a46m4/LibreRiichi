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
const ManzuBit = 0;
const PinzuBit = 1;
const SouzuBit = 2;
const HonourBit = 3;
const NumberMask = 0b1111;
const SpecialMask = 0b11 << 6;

function decode(array: string): Uint8Array {
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
    private value: number

    constructor(value: number) {
        this.value = value
    }

    get getValue(): number {
        return this.value;
    }

    clearRedOrDora(): Tile {
        return new Tile(this.value & ~(TileValue.DoraTile | TileValue.RedTile));
    }

    isInvalid(): boolean {
        return this.value === TileValue.Invalid;
    }

    isHidden(): boolean {
        return this.value === TileValue.Hidden;
    }

    isHonour(): boolean {
        return (this.value & TileMask) === HonourBit;
    }

    isWind(): boolean {
        return this.value >= 48 && this.value <= 51;
    }

    isDragon(): boolean {
        return this.value >= 52 && this.value <= 54;
    }

    isManzu(): boolean {
        return (this.value & TileMask) === ManzuBit;
    }

    isPinzu(): boolean {
        return (this.value & TileMask) === PinzuBit;
    }

    getTileNumber(): number {
        return this.value & NumberMask;
    }

    setRedTile(): Tile {
        return new Tile(this.value | TileValue.RedTile);
    }

    setDoraTile(): Tile {
        return new Tile(this.value | TileValue.DoraTile);
    }

    setTileNumber(num: number): Tile {
        return new Tile((this.value & (TileMask | SpecialMask)) | num);
    }
}