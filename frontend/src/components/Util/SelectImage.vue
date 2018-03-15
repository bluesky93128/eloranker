<template>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Select Image</p>
      </header>
      <section class="modal-card-body">
        <div class="tabs">
          <ul>
            <li :class="[{ 'is-active' : optionTab == 0 }]" @click="optionTab = 0"><a>From URL</a></li>
            <li :class="[{ 'is-active' : optionTab == 1 }]" @click="optionTab = 1"><a>Upload an image</a></li>
            <li :class="[{ 'is-active' : optionTab == 2 }]" @click="optionTab = 2"><a>Google Images</a></li>
          </ul>
        </div>
        <div v-if="optionTab == 0">
          <input class="input" type="text" placeholder="URL">
        </div>
        <div v-else-if="optionTab == 1">
          <div class="file">
            <label class="file-label">
              <input class="file-input" type="file" name="resume">
              <span class="file-cta">
                <span class="file-icon">
                  <svg class="svg-inline--fa fa-upload fa-w-16" aria-hidden="true" data-prefix="fas" data-icon="upload" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" data-fa-i2svg=""><path fill="currentColor" d="M296 384h-80c-13.3 0-24-10.7-24-24V192h-87.7c-17.8 0-26.7-21.5-14.1-34.1L242.3 5.7c7.5-7.5 19.8-7.5 27.3 0l152.2 152.2c12.6 12.6 3.7 34.1-14.1 34.1H320v168c0 13.3-10.7 24-24 24zm216-8v112c0 13.3-10.7 24-24 24H24c-13.3 0-24-10.7-24-24V376c0-13.3 10.7-24 24-24h136v8c0 30.9 25.1 56 56 56h80c30.9 0 56-25.1 56-56v-8h136c13.3 0 24 10.7 24 24zm-124 88c0-11-9-20-20-20s-20 9-20 20 9 20 20 20 20-9 20-20zm64 0c0-11-9-20-20-20s-20 9-20 20 9 20 20 20 20-9 20-20z"></path></svg>
                </span>
                <span class="file-label">
                  Choose a file…
                </span>
              </span>
            </label>
          </div>
        </div>
        <div v-else>
          First Google Images search result will be used.
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-success" @click="apply">Apply</button>
        <button class="button" @click="cancel">Cancel</button>
      </footer>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';

@Component
export default class SelectImage extends Vue {
  optionTab = 0;

  async apply() {
    this.$emit('close', this.optionTab);
  }

  async cancel() {
    this.$emit('close');
  }
}
</script>

<style lang="scss" module>
</style>