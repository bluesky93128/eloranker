<template>
  <section class="hero is-primary is-red">
    <div class="hero-body">
      <div class="container">
        <div class="columns is-centered">
          <div class="column is-half has-text-centered">
            <div v-if="isCreating && isAdmin">
              <div class="field">
                <div class="control has-icons-right">
                  <input class="input is-rounded is-large" :value="roomTitle" @input="onTitleChange" :disabled="!isAdmin" placeholder="Type poll name here...">
                  <span class="icon is-large is-right">
                    <i class="icon-pencil-squared fa-3x"></i>
                  </span>
                </div>
              </div>
              <h5 v-if="this.$store.state.variants.length < 3" class="subtitle is-5">
                * You need to add <b>{{ 3 - this.$store.state.variants.length }}</b> more items to start voting
              </h5>
              <div v-else>
                <router-link
                  class="button is-rounded is-large is-warning"
                  tag="button"
                  :to="{ name: 'room-voting', params: $route.params }"
                >Start voting</router-link>
              </div>
            </div>
            <div v-else>
              <p :class="['title', 'is-1']">{{ roomTitle }}</p>
              <p class="subtitle is-5">
                Rank a list of items by comparing them 1 vs 1
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { State } from 'vuex-class';
import { EditMode } from '@/room';
import connection from '@/connection';

@Component
export default class RoomStatus extends Vue {
  $refs!: { shareableLink: HTMLInputElement };
  @State roomName!: string;
  @State roomSecret!: string;
  @State clientNumber!: number;
  @State isAdmin!: boolean;
  exposeSecret = false;
  editingName = false;

  get isCreating() {
    return this.$route.name == 'room-list';
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
.question {
  font-size: 6rem;
}

.shareableLink {
  width: 400px;
}

.roomTitle {
  width: 400px;
}
</style>
