<script setup lang="ts">
import { ref, computed } from 'vue'
import { Tile, TileValue } from "../game/tile"

const props = defineProps<{
  tiles: Tile[]
}>()

const gameContainer = ref<HTMLDivElement>()
const isFullscreen = ref(false)

// Mock data for demonstration - in real game this would come from props/store
const playerTiles = computed(() => props.tiles || [
  new Tile(TileValue.Manzu + 0),
  new Tile(TileValue.Manzu + 1),
  new Tile(TileValue.Manzu + 2),
  new Tile(TileValue.Pinzu + 0),
  new Tile(TileValue.Pinzu + 1),
  new Tile(TileValue.Souzu + 0),
  new Tile(TileValue.Souzu + 1),
  new Tile(TileValue.EastTile),
  new Tile(TileValue.SouthTile),
  new Tile(TileValue.WestTile),
  new Tile(TileValue.White),
  new Tile(TileValue.Red),
  new Tile(TileValue.Green)
])

// Other players' tiles (face down)
const otherPlayerTiles = computed(() =>
  Array(13).fill(null).map(() => new Tile(TileValue.Hidden))
)

// Center discarded tiles
const discardedTiles = computed(() => [
  new Tile(TileValue.Manzu + 5),
  new Tile(TileValue.Pinzu + 3),
  new Tile(TileValue.Souzu + 7),
  new Tile(TileValue.EastTile),
  new Tile(TileValue.White),
  new Tile(TileValue.Manzu + 1),
  new Tile(TileValue.Pinzu + 8),
  new Tile(TileValue.Souzu + 2)
])

function getTileImageUrl(tile: Tile): string {
  return tile.get_image_path_static()
}
</script>

<template>
  <div ref="gameContainer" class="game-view-container" :class="{ fullscreen: isFullscreen }">
    <div class="mahjong-table">
      <!-- North player (top) -->
      <div class="player-area north">
        <div class="tile-row">
          <img
            v-for="(tile, index) in otherPlayerTiles"
            :key="`north-${index}`"
            :src="getTileImageUrl(tile)"
            class="tile tile-small"
            alt="Hidden tile"
          />
        </div>
        <div class="player-info">Player 3</div>
      </div>

      <!-- West player (left) -->
      <div class="player-area west">
        <div class="tile-column">
          <img
            v-for="(tile, index) in otherPlayerTiles"
            :key="`west-${index}`"
            :src="getTileImageUrl(tile)"
            class="tile tile-small rotated-90"
            alt="Hidden tile"
          />
        </div>
        <div class="player-info">Player 2</div>
      </div>

      <!-- Center area with discarded tiles -->
      <div class="center-area">
        <div class="table-surface">
          <div class="discard-pile">
            <img
              v-for="(tile, index) in discardedTiles"
              :key="`discard-${index}`"
              :src="getTileImageUrl(tile)"
              class="tile tile-discarded"
              :style="{
                transform: `rotate(${Math.random() * 20 - 10}deg)`,
                zIndex: index
              }"
              alt="Discarded tile"
            />
          </div>
          <div class="dealer-marker"></div>
        </div>
      </div>

      <!-- East player (right) -->
      <div class="player-area east">
        <div class="tile-column">
          <img
            v-for="(tile, index) in otherPlayerTiles"
            :key="`east-${index}`"
            :src="getTileImageUrl(tile)"
            class="tile tile-small rotated-270"
            alt="Hidden tile"
          />
        </div>
        <div class="player-info">Player 4</div>
      </div>

      <!-- South player (bottom - current player) -->
      <div class="player-area south">
        <div class="player-info">You</div>
        <div class="tile-row">
          <img
            v-for="(tile, index) in playerTiles"
            :key="`south-${index}`"
            :src="getTileImageUrl(tile)"
            class="tile tile-player"
            alt="Player tile"
            @click="() => console.log('Clicked tile:', tile)"
          />
        </div>
      </div>
    </div>

    <div class="game-ui">
      <div class="game-info">
        <h3>Riichi Mahjong</h3>
        <p>East Round - Hand 1</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-view-container {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  background: linear-gradient(135deg, #0d4f3c 0%, #1a7c5a 100%);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.mahjong-table {
  position: relative;
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-areas:
    ". north ."
    "west center east"
    ". south .";
  grid-template-columns: 1fr 2fr 1fr;
  grid-template-rows: 1fr 2fr 1fr;
  padding: 20px;
  gap: 10px;
}

.player-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.player-area.north {
  grid-area: north;
}

.player-area.west {
  grid-area: west;
  flex-direction: row;
}

.player-area.east {
  grid-area: east;
  flex-direction: row-reverse;
}

.player-area.south {
  grid-area: south;
  flex-direction: column-reverse;
}

.center-area {
  grid-area: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

.table-surface {
  position: relative;
  width: 300px;
  height: 300px;
  background: radial-gradient(ellipse at center, #0a5c42 0%, #083d2e 70%);
  border-radius: 20px;
  border: 4px solid #654321;
  box-shadow: inset 0 2px 10px rgba(0,0,0,0.3), 0 4px 20px rgba(0,0,0,0.2);
}

.discard-pile {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
}

.dealer-marker {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 24px;
  height: 24px;
  background: #ff6b6b;
  border-radius: 50%;
  border: 2px solid white;
  box-shadow: 0 2px 6px rgba(0,0,0,0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.tile-row {
  display: flex;
  gap: 4px;
  align-items: center;
}

.tile-column {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
}

.tile {
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
  transition: all 0.2s ease;
  cursor: pointer;
  border: 2px solid #fff;
}

.tile-player {
  width: 48px;
  height: 66px;
}

.tile-player:hover {
  transform: translateY(-6px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.4);
  border-color: #4a90e2;
}

.tile-small {
  width: 32px;
  height: 44px;
  opacity: 0.9;
}

.tile-discarded {
  position: absolute;
  width: 28px;
  height: 38px;
  opacity: 0.8;
  left: 50%;
  top: 50%;
  transform-origin: center;
  margin-left: -14px;
  margin-top: -19px;
}

.rotated-90 {
  transform: rotate(90deg);
}

.rotated-270 {
  transform: rotate(-90deg);
}

.player-info {
  margin: 8px 0;
  padding: 4px 12px;
  background: rgba(0,0,0,0.7);
  color: white;
  border-radius: 16px;
  font-size: 12px;
  font-weight: 500;
  text-align: center;
  min-width: 60px;
}

.player-area.south .player-info {
  background: rgba(255, 107, 107, 0.8);
}

.game-ui {
  position: absolute;
  top: 16px;
  left: 16px;
  z-index: 10;
}

.game-info {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
  padding: 12px 16px;
  border-radius: 8px;
  color: white;
  font-size: 14px;
}

.game-info h3 {
  margin: 0 0 4px 0;
  font-size: 16px;
  font-weight: 600;
}

.game-info p {
  margin: 0;
  opacity: 0.8;
  font-size: 12px;
}

.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  border-radius: 0;
}

/* Responsive design */
@media (max-width: 768px) {
  .tile-player {
    width: 36px;
    height: 50px;
  }

  .tile-small {
    width: 24px;
    height: 33px;
  }

  .table-surface {
    width: 250px;
    height: 250px;
  }

  .mahjong-table {
    padding: 10px;
  }
}
</style>
