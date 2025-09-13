<script setup lang="ts">
import {BoxStyling, ButtonStyling, FlexBox, H1Styling, InputStyling, ULStyling} from "../styling";
import {Ref, ref} from "vue";
import ListItem from "../components/list_item.vue";
import {MessageType} from "../messaging/message";
import {use_room_state, use_websocket_state} from "../index";
import {ServerActionType} from "../messaging/server_action_generated";
import {ServerResponseType} from "../messaging/server_response_generated";
import {register_request} from "../messaging/event_handler";
import router from "../router";

const websocket_state = use_websocket_state()
const room_state = use_room_state()

const create_room_name = ref('')
const show_error = ref('')
const avail_rooms: Ref<string[]> = ref([])

async function check_avail_rooms() {
  let msg_idx = websocket_state.conn.send({
    message_type: MessageType.REQUEST,
    data: {
      serveraction_type: ServerActionType.ListArenasAction,
    }
  })

  let msg = await register_request(msg_idx)

  if (msg.serverresponse_type !== ServerResponseType.ListArenasResponse) {
    throw new Error("Wrong type")
  }

  if (!msg.success) {
    throw new Error("Failed to list rooms")
  }

  avail_rooms.value = msg.arena_list.sort()
}

async function find_room() {
    let msg_idx = websocket_state.conn.send(
        {
          message_type: MessageType.REQUEST,
          data: {
            serveraction_type: ServerActionType.JoinArenaAction,
            arena_name: room_state.room_name,
          }
        }
    )

    let ret = await register_request(msg_idx)
    if (ret.serverresponse_type !== ServerResponseType.GenericResponse) {
      show_error.value = "Connection error: Wrong type"
      return
    }

    if (!ret.success) {
      show_error.value = "Could not join room: " + ret.fail_reason
      return
    }

    console.log("Joined room")
    room_state.room_set = true
    await router.push({name: 'arena_page'})
}

async function create_room() {
  let msg_idx = websocket_state.conn.send({
    message_type: MessageType.REQUEST,
    data: {
      serveraction_type: ServerActionType.CreateArenaAction,
      arena_name: create_room_name.value,
    }
  })

  let msg = await register_request(msg_idx)

  if (msg === undefined) {
    // TODO: Give a reason for why
    show_error.value = "Connection error: Failed to create room"
  }

  if (msg.serverresponse_type !== ServerResponseType.GenericResponse) {
    show_error.value = "Connection error: Wrong type"
    return
  }

  if (!msg.success) {
    show_error.value = "Could not create room: " + msg.fail_reason
  }
}

</script>

<template>
  <div :class="BoxStyling">
  <h1 :class="H1Styling">Join Room</h1>
  <p>Room Name</p>
  <input
      :class="InputStyling"
      v-model="room_state.room_name">
  <button
      :class="ButtonStyling"
      @click="find_room">Find</button>
  </div>
  <div :class="BoxStyling">
    <h1 :class="H1Styling">Create room</h1>
    <input
        :class="InputStyling"
        v-model="create_room_name">
    <button
        :class="ButtonStyling"
        @click="create_room">Create</button>
  </div>

  <div :class="BoxStyling">
    <div :class="FlexBox">
      <h1 :class="H1Styling">Available rooms</h1>
      <button :class="ButtonStyling" @click="check_avail_rooms">Find</button>
      <br>
    </div>
    <ul :class="ULStyling" v-if="avail_rooms.length !== 0">
      <ListItem v-for="avail_room in avail_rooms">{{ avail_room }}</ListItem>
    </ul>
  </div>



</template>

<style scoped>

</style>