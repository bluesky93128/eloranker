<template>
  <div>
    <section class="hero is-primary is-red">
      <div class="hero-body">
        <div class="container">
          <div class="columns">
            <div class="column is-one-third">
              <div class="content">
                <h1 class="title">
                  Creating new poll
                </h1>

                <h5 class="subtitle">
                  Enter poll name and add minimum of 3 options
                </h5>

                <div class="field">
                  <div class="control">
                    <input :class="['input', 'is-rounded', 'is-large']" :value="roomTitle" @input="onTitleChange" :disabled="!isAdmin" placeholder="Poll Name">
                  </div>
                </div>

                <div class="field is-grouped">
                  <router-link
                    class="control button is-rounded is-large is-warning"
                    tag="button"
                    :to="{ name: 'room-list', params: $route.params }"
                  >
                    <span class="icon is-small">
                      <i class="icon-pencil-squared"></i>
                    </span>
                    <span>
                      Edit
                    </span>
                  </router-link>
                  <router-link
                    class="control button is-rounded is-large"
                    tag="button"
                    :to="{ name: 'room-voting', params: $route.params }"
                    :disabled="!canVote"
                  >   Start voting   </router-link>
                </div>
              </div>
            </div>
            <div class="column is-one-third">
            </div>
            <div class="column is-one-third">
              <div class="content">
                <h1 class="title">
                  Settings
                </h1>

                <input type="checkbox" v-model="exposeSecret" v-if="isAdmin">
                <div class="field has-addons">
                  <p class="control">
                    <input class="input" type="text" readonly :value="shareableLink" ref="shareableLink">
                  </p>
                  <p class="control">
                    <a class="button is-warning">
                      <span class="icon is-small is-right" @click="copyLink">
                        <i class="icon-clipboard"></i>
                      </span>
                    </a>
                  </p>
                </div>

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

                <input type="checkbox" :checked="roomQuotaEnabled" @change="onQuotaChange" :disabled="!isAdmin">
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
              </div>

              <h5 class="subtitle">Connected Clients: {{ clientNumber }}</h5>
            </div>
          </div>
        </div>
      </div>
    </section>
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
