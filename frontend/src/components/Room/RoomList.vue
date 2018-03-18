<template>
  <div class="columns is-multiline is-mobile">
    <!-- <div :class="['tile', 'is-5', 'is-parent']"> -->
    <VariantElement @updateSelection="updateSelection" :number="0" />
    <!-- <VariantElement v-if="sortedVariants.length < 1" @updateSelection="updateSelection" :number="sortedVariants.length+1" /> -->
    <VariantElement
      ref="elements"
      v-for="(variant, index) in sortedVariants"
      :number="index + 1"
      :key="variant.uuid"
      :variant="variant"
    />
    <!-- </div> -->
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter } from 'vuex-class';
import { Variant } from '@/room';
import VariantElement from './VariantElement.vue';

@Component({ components: { VariantElement } })
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
