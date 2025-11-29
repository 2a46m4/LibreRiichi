<script setup lang="ts">
import { ref } from 'vue'
import { BoxStyling, ButtonStyling, H1Styling, InputStyling } from '../styling'
import ErrorDisplay from '../components/error_display.vue'
import { router, use_player_state, use_websocket_state } from '../index'
import { MessageType } from '../messaging/message'
import { ServerResponseType } from '../messaging/server_response_generated'
import { ServerActionType } from '../messaging/server_action_generated'
import { register_request } from '../messaging/event_handler'

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

  let message_return = await register_request(return_index)
  if (
    message_return.serverresponse_type !== ServerResponseType.GenericResponse
  ) {
    status.value = 'Unexpected message type'
    return
  }

  if (!message_return.success) {
    status.value = message_return.fail_reason
    return
  }

  await router.push({ name: 'connected_page' })
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
