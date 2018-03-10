<template>
  <div :class="$style.variants">
    <VariantElement
      ref="elements"
      v-for="variant in sortedVariants"
      :key="variant.uuid"
      :variant="variant"
    />

    <VariantElement @updateSelection="updateSelection" />
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
.variants {
  display: flex;
  flex-flow: row wrap;

  > :last-child {
    opacity: 0.5;
  }
}
</style>
