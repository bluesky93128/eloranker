<template>
  <div class="column is-one-quarter">
    <div class="columns is-multiline">
      <div class="column">
        <div class="card">
          <header class="card-header">
            <p class="card-header-title">
              Share
            </p>
          </header>
          <div class="card-content">
            <input type="checkbox" v-model="exposeSecret" v-if="isAdmin" hidden>
            <ShareBlock/>
          </div>
        </div>
      </div>
      <div class="column">
        <div class="card">
          <header class="card-header">
            <p class="card-header-title">
              Settings
            </p>
            <div class="card-header-icon">
              <span class="tag is-info">Connected Clients: {{ clientNumber }}</span>
            </div>
          </header>
          <div class="card-content">
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

            <ImportBlock />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter, State } from 'vuex-class';
import connection from '@/connection';
import { EditMode, SortingOrder } from '@/room';
import ShareBlock from './ShareBlock.vue';
import ImportBlock from './ImportBlock.vue';

@Component({ components: { ShareBlock, ImportBlock } })
export default class RoomSettings extends Vue {
  $refs!: { shareableLink: HTMLInputElement };
  @State roomName!: string;
  @State roomSecret!: string;
  @State clientNumber!: number;
  @Getter canVote!: boolean;
  @State isAdmin!: boolean;
  exposeSecret = false;

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
