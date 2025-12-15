<script setup lang="ts">
import { ref } from 'vue'
import { BoxStyling, ButtonStyling, H1Styling, InputStyling } from '../styling'
import ErrorDisplay from '../components/error_display.vue'
import { router, use_player_state, use_websocket_state } from '../index'
import { MessageType } from '../messaging/message'
import { ServerActionType } from '../messaging/server_action_generated'
import { get_response } from '../messaging/event_handler'

const player_state = use_player_state()
const websocket_state = use_websocket_state()
const status = ref('')

async function connect() {
  await websocket_state.conn.wait_until_ready()
  let return_index = websocket_state.conn.send({
    message_type: MessageType.REQUEST,
    data: {
      serveraction_type: ServerActionType.InitialMessageAction,
      name: player_state.username,
    },
  })

  try {
    await get_response(return_index)
    await router.push({ name: 'connected_page' })
  } catch (e) {
    console.error(e)
    if (e instanceof Error) {
      status.value = e.message
    }
  }
}

// TMP
connect()
</script>

<template>
  <div :class="BoxStyling">
    <h1 :class="H1Styling">LibreRiichi</h1>
    <p>Username</p>
    <input :class="InputStyling" v-model="player_state.username" />
    <button :class="ButtonStyling" @click="connect" @keyup.enter="connect">
      Connect
    </button>
    <ErrorDisplay :error="status" v-if="status.length !== 0"></ErrorDisplay>
  </div>
</template>

<style>
button {
  margin-left: 10px;
}
</style>
