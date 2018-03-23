<template>
  <div class="columns is-mobile">
    <div :class="['column', 'is-four-fifths']">
      <div class="columns is-multiline is-mobile">
        <VariantElement @updateSelection="updateSelection" :number="0" />
        <VariantElement
          ref="elements"
          v-for="(variant, index) in sortedVariants"
          :number="index + 1"
          :key="variant.uuid"
          :variant="variant"
        />
      </div>
    </div>
    <RoomSettings />
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter } from 'vuex-class';
import { Variant } from '@/room';
import VariantElement from './VariantElement.vue';
import RoomSettings from '@/components/Room/RoomSettings.vue';

@Component({ components: { VariantElement, RoomSettings } })
export default class RoomList extends Vue {
  $refs!: { elements: VariantElement[] };
  @Getter sortedVariants!: Variant[];

  public updateSelection(id: string, selected: 'textInput' | 'imageInput') {
    const element = this.$refs.elements.find(e => e.variant.uuid === id);
    if (!element) return;

    element.$refs[selected].select();
  }
}
</script>

<style lang="scss" module>
</style>
