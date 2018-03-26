<template>
  <section class="hero is-primary is-red">
    <div class="hero-body">
      <div class="container">
        <div class="columns">
          <div class="column is-two-thirds">
            <!-- <div class="content"> -->
            <div v-if="isAdmin">
              <!-- <h2 class="title is-2">
                Creating new poll
              </h2> -->
              <!-- <h2 class="title is-2">{{ roomTitle }}</h2> -->
              <div class="field has-addons">
                <div class="control has-icons-right">
                  <input class="input is-rounded is-large" :value="roomTitle" @input="onTitleChange" :disabled="!isAdmin" placeholder="Poll Name">
                  <span class="icon is-large is-right">
                    <i class="icon-pencil-squared fa-3x"></i>
                  </span>
                </div>
                <!-- <div class="control">
                  <a class="button is-large is-warning">
                    <span class="icon is-medium is-right" @click="editName">
                      <i class="icon-pencil-squared"></i>
                    </span>
                  </a>
                </div> -->
              </div>
              <h5 class="subtitle is-5">
                Enter poll name and add minimum of 3 options
              </h5>
            </div>
            <div v-else>
              <p class="title is-1">{{ roomTitle }}</p>
              <p class="subtitle is-5">
                Rank a list of items by comparing them 1 vs 1
              </p>
            </div>
            <!-- <div class="field is-grouped">
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
            </div> -->
            <!-- </div> -->
          </div>
          <div class="column is-third">
            <p class="title is-1">Share</p>
            <input type="checkbox" v-model="exposeSecret" v-if="isAdmin" hidden>
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
              <p class="control">
                <a class="button is-info">
                  <span class="icon is-small is-right" @click="copyLink">
                    <i class="icon-twitter"></i>
                  </span>
                </a>
              </p>
              <p class="control">
                <a class="button is-link">
                  <span class="icon is-small is-right" @click="copyLink">
                    <i class="icon-facebook-squared"></i>
                  </span>
                </a>
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
import { Getter, State } from 'vuex-class';
import connection from '@/connection';
import { EditMode } from '@/room';

@Component
export default class RoomStatus extends Vue {
  $refs!: { shareableLink: HTMLInputElement };
  @State roomName!: string;
  @State roomSecret!: string;
  @State clientNumber!: number;
  @Getter canVote!: boolean;
  @State isAdmin!: boolean;
  exposeSecret = false;
  editingName = false;

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
