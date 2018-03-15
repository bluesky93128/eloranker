<template>
  <div :class="$style.variant">
    <img
      :class="$style.image"
      :src="variant.image || fallbackImage"
      @error="$event.target.src = fallbackImage"
    >
    <input
      ref="imageInput"
      v-if="!voting"
      :class="$style.imageInput"
      v-model="variant.image"
      @input="pushVariantUpdate"

      :readonly="!hasEditPermissions"
    >
    <input
      ref="textInput"
      class="input"
      :class="[$style.textInput, { 'is-static': voting }]"
      v-model="variant.text"
      @input="onTextInput"
      :list="autocompleteId"
      maxlength="100"

      :readonly="!hasEditPermissions"
    >
    <!-- TODO: Serverside maxlength validation -->

    <datalist :id="autocompleteId">
      <option v-for="text in autocomplete" :key="text" :value="text" />
    </datalist>

    <div v-if="!voting && !isNew">{{ variant.rating }}</div>
    <button
      class="button"
      @click="findImage"
      v-if="!voting"
      :disabled="!canFindImage"
    >Google Images</button>

    <button
      class="button"
      v-if="!isNew"
      :disabled="!isIgnored && !canIgnoreVariant"
      @click="setIgnored(!isIgnored)"
    >{{ isIgnored ? 'UNIGNORE' : 'IGNORE' }}</button>

    <button
      class="button"
      @click="remove"
      v-if="!isNew && !voting"
      :disabled="!hasEditPermissions"
    >RM</button>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';
import { Getter } from 'vuex-class';
import connection from '@/connection';
import { Variant, emptyVariant } from '@/room';

import GoogleImageSearch from 'free-google-image-search';
import googleAutocomplete from '@/google-autocomplete';

@Component
export default class VariantElement extends Vue {
  $refs!: { textInput: HTMLInputElement; imageInput: HTMLInputElement };

  @Getter canIgnoreVariant!: boolean;

  @Prop({ type: Object, default: emptyVariant })
  variant!: Variant;

  @Prop({ type: Boolean, default: false })
  voting!: boolean;

  autocomplete: string[] = [];
  fallbackImage: string = require('@/assets/no-image.png');
  waitingForAllocation = false;

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

  async findImage() {
    if (!this.canFindImage) return;
    const variant = this.variant!;

    const images = await GoogleImageSearch.searchImage(variant.text);
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
.variant {
  border: 1px solid red;
  margin: 2px;
  width: 350px;
  display: flex;
  flex-flow: column;
}

.image {
  height: 150px;
  object-fit: contain;
}

.imageInput {
  width: 100%;
}

.textInput {
  width: 100%;
  min-height: 3rem;
  font-size: 1.25rem;
}
</style>
