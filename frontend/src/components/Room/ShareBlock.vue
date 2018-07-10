<template>
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
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter, State } from 'vuex-class';

@Component
export default class ShareBlock extends Vue {
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
}
</script>

<style lang="scss" module>
</style>
