<template>
  <div class="column is-one-quarter">
    <div class="card">
      <header class="card-header">
        <p class="card-header-title">
          Settings
        </p>
      </header>
      <div v-if="!voting" class="card-content">
        <div class="field">
          <div class="control">
            <div class="select">
              <select v-model="sortingOrder">
                <option :value="SortingOrder.DATE">Sort by Date</option>
                <option :value="SortingOrder.RATING">Sort by Rating</option>
              </select>
            </div>
          </div>
        </div>

        <input type="checkbox" :checked="roomQuotaEnabled" @change="onQuotaChange" :disabled="!isAdmin" hidden>
        <div class="field">
          <div class="control">
            <div class="select">
              <select :value="roomEditMode" @input="onEditModeChange" :disabled="!isAdmin">
                <option :value="EditMode.Trust">Trust Mode</option>
                <option :value="EditMode.Normal">Normal Mode</option>
                <option :value="EditMode.Restricted">Restricted Mode</option>
              </select>
            </div>
          </div>
        </div>

        <h5 class="subtitle">Connected Clients: {{ clientNumber }}</h5>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter, State } from 'vuex-class';
import connection from '@/connection';
import { EditMode, SortingOrder } from '@/room';

@Component
export default class RoomSettings extends Vue {
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
