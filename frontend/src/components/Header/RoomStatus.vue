<template>
  <div :class="$style.roomStatus">
    <div>
      <input type="checkbox" v-model="exposeSecret" v-if="isAdmin">
      <input :class="$style.shareableLink" type="text" readonly :value="shareableLink" ref="shareableLink">
      <button @click="copyLink">COPY</button>
    </div>

    <div>Connected Clients: {{ clientNumber }}</div>


    <select v-model="sortingOrder">
      <option :value="SortingOrder.DATE">Sort by Date</option>
      <option :value="SortingOrder.RATING">Sort by Rating</option>
    </select>


    <div>
      <router-link
        class="button"
        tag="button"
        :to="{ name: 'room-list', params: $route.params }"
      >EDIT</router-link>
      <router-link
        class="button"
        tag="button"
        :to="{ name: 'room-voting', params: $route.params }"
        :disabled="!canVote"
      >VOTE</router-link>
    </div>


    <div>
      <input type="checkbox" :checked="roomQuotaEnabled" @change="onQuotaChange" :disabled="!isAdmin">
      <select :value="roomEditMode" @input="onEditModeChange" :disabled="!isAdmin">
        <option :value="EditMode.Trust">Trust Mode</option>
        <option :value="EditMode.Normal">Normal Mode</option>
        <option :value="EditMode.Restricted">Restricted Mode</option>
      </select>
      <input :class="$style.roomTitle" :value="roomTitle" @input="onTitleChange" :disabled="!isAdmin">
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter, State } from 'vuex-class';
import connection from '@/connection';
import { EditMode, SortingOrder } from '@/room';

@Component
export default class RoomStatus extends Vue {
  $refs!: { shareableLink: HTMLInputElement };
  @State roomName!: string;
  @State roomSecret!: string;
  @State clientNumber!: number;
  @Getter canVote!: boolean;
  @State isAdmin!: boolean;
  exposeSecret = false;

  SortingOrder = SortingOrder;
  get sortingOrder() {
    return this.$store.state.sortingOrder;
  }
  set sortingOrder(value: SortingOrder) {
    this.$store.commit('setSortingOrder', value);
  }

  get shareableLink() {
    const shouldExposeSecret = this.exposeSecret && this.roomSecret;
    const roomId = `${this.roomName}${shouldExposeSecret ? `!${this.roomSecret}` : ''}`;

    return `${window.location.origin}/${roomId}`;
  }

  copyLink() {
    const el = this.$refs.shareableLink;

    el.select();
    document.execCommand('copy');
    window.getSelection().empty();
  }

  @State roomQuotaEnabled!: boolean;
  onQuotaChange(event: Event & { target: HTMLInputElement }) {
    const quotaEnabled = event.target.checked;
    connection.setQuotaEnabled(quotaEnabled);
    this.$store.commit('setQuotaEnabled', quotaEnabled);
  }

  EditMode = EditMode;
  @State roomEditMode!: EditMode;
  onEditModeChange(event: Event & { target: HTMLInputElement }) {
    const editMode = Number(event.target.value);
    connection.setEditMode(editMode);
    this.$store.commit('setEditMode', editMode);
  }

  @State roomTitle!: string;
  onTitleChange(event: Event & { target: HTMLInputElement }) {
    const title = event.target.value;
    connection.setTitle(title);
    this.$store.commit('setTitle', title);
  }
}
</script>

<style lang="scss" module>
.shareableLink {
  width: 400px;
}

.roomTitle {
  width: 400px;
}
</style>
