<template>
  <div :class="['column', voting ? 'is-one-third' : 'is-one-quarter']">
    <div class="card">
      <header class="card-header">
        <p class="card-header-title">
          {{ (number > 0 ? (number) + ". " : "") + variant.text }}
        </p>
        <div class="card-header-icon dropdown is-hoverable">
          <div class="dropdown-trigger">
            <span class="icon">
              <i class="icon-menu"></i>
            </span>
          </div>
          <div class="dropdown-menu" id="dropdown-menu4" role="menu">
            <div class="dropdown-content">
              <a class="dropdown-item">
                {{ 'ELO: ' + variant.rating }}
              </a>
              <hr class="dropdown-divider">
              <a class="dropdown-item" @click="findImage">
                <span class="icon">
                  <i class="icon-google"></i>
                </span>
                From Google Images
              </a>
              <a class="dropdown-item" @click="openImageSelector">
                <span class="icon">
                  <i class="icon-picture"></i>
                </span>
                From file/URL
              </a>
              <hr class="dropdown-divider">
              <a v-if="voting" class="dropdown-item" @click="setIgnored(!isIgnored)" :disabled="!isIgnored && !canIgnoreVariant">
                <span class="icon">

                </span>
                {{ isIgnored ? 'Unignore' : 'Ignore' }}
              </a>
              <a class="dropdown-item" @click="remove">
                <span class="icon">
                  <i class="icon-trash"></i>
                </span>
                Delete
              </a>
            </div>
          </div>
        </div>
      </header>
      <div class="card-image">
        <figure class="image is-4by3">
          <img
            :src="variant.image || fallbackImage"
            @error="$event.target.src = fallbackImage"
          >
        </figure>
      </div>
      <div v-if="!voting" class="card-content">
        <div class="field">
          <div :class="['control', { 'is-loading': waitingForImage }]">
            <input
              ref="textInput"
              class="input"
              :class="[{ 'is-static': voting }]"
              v-model="variant.text"
              @input="onTextInput"
              :list="autocompleteId"
              maxlength="100"
              placeholder="Option Name"
              :readonly="!hasEditPermissions"
            >
          </div>
        </div>
        <div class="field">
          <div :class="['control']">
            <input
              type="hidden"
              ref="imageInput"
              v-if="!voting"
              class="input"
              v-model="variant.image"
              @input="pushVariantUpdate"

              :readonly="!hasEditPermissions"
            >
          </div>
        </div>
        <!-- TODO: Serverside maxlength validation -->

        <datalist :id="autocompleteId">
          <option v-for="text in autocomplete" :key="text" :value="text" />
        </datalist>

        <SelectImage v-show="selectingImage" @close="closeImageSelector"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';
import { Getter } from 'vuex-class';
import connection from '@/connection';
import { Variant, emptyVariant } from '@/room';

import GoogleImageSearch from 'free-google-image-search';
import googleAutocomplete from '@/google-autocomplete';

import SelectImage from '@/components/Util/SelectImage.vue';

@Component({ components: { SelectImage } })
export default class VariantElement extends Vue {
  $refs!: { textInput: HTMLInputElement; imageInput: HTMLInputElement };

  @Getter canIgnoreVariant!: boolean;

  @Prop({ type: Object, default: emptyVariant })
  variant!: Variant;

  @Prop({ type: Boolean, default: false })
  voting!: boolean;

  @Prop({ type: Number, default: 0 })
  number!: number;

  autocomplete: string[] = [];
  fallbackImage: string = require('@/assets/no-image.png');
  waitingForAllocation = false;
  waitingForImage = false;
  selectingImage = false;

  get isIgnored() {
    return this.$store.getters.isIgnoredVariant(this.variant.uuid);
  }

  setIgnored(ignored: boolean) {
    const id = this.variant.uuid;
    this.$store.commit('setVariantIgnored', { id, ignored });
    connection.setVariantIgnored(id, ignored);
  }

  get isNew() {
    return this.variant.uuid === '';
  }

  get canFindImage() {
    return this.hasEditPermissions && this.variant.text !== '' && !this.isNew;
  }

  get autocompleteId() {
    return `variant-${this.variant.uuid}`;
  }

  async openImageSelector() {
    this.selectingImage = true;
  }

  async closeImageSelector(option: Number) {
    console.log(option);
    this.selectingImage = false;
  }

  async findImage() {
    if (!this.canFindImage) return;
    const variant = this.variant!;

    this.waitingForImage = true;
    const images = await GoogleImageSearch.searchImage(variant.text);
    this.waitingForImage = false;
    if (images.length === 0) return;

    variant.image = images[0];

    this.pushVariantUpdate();
  }

  getSelectedElement() {
    switch (document.activeElement) {
      case this.$refs.textInput:
        return 'textInput';
      case this.$refs.imageInput:
        return 'imageInput';
    }
  }

  async pushVariantUpdate() {
    if (!this.isNew) {
      connection.updateVariant(this.variant);
      return;
    }

    if (this.waitingForAllocation) return;
    this.waitingForAllocation = true;

    const selected = this.getSelectedElement();
    const id = await connection.allocateNewVariant();
    if (selected != null) this.$emit('updateSelection', id, selected);

    const v: Variant | undefined = this.$store.getters.findVariant(id);
    if (v == null) return;

    v.text = this.variant.text;
    v.image = this.variant.image;
    connection.updateVariant(v);

    Object.assign(this.variant, emptyVariant());

    this.waitingForAllocation = false;
  }

  async onTextInput() {
    this.pushVariantUpdate();
    if (this.isNew) return;

    const { text } = this.variant;
    if (text.length <= 1) {
      this.autocomplete = [];
      return;
    }

    this.autocomplete = this.autocomplete.filter(v =>
      v.toLowerCase().startsWith(text.toLowerCase()),
    );
    this.autocomplete = await googleAutocomplete(text);
  }

  get hasEditPermissions() {
    return this.$store.getters.hasWriteAccess(this.variant.uuid);
  }

  async remove() {
    const id = this.variant.uuid;
    this.$store.commit('removeVariant', id);
    connection.removeVariant(id);
  }
}
</script>

<style lang="scss" module>
</style>
