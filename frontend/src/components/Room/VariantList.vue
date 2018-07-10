<template>
  <div class="columns is-multiline">
    <VariantElement v-if="canCreate" @updateSelection="updateSelection" :number="0" />
    <VariantElement
      ref="elements"
      v-for="(variant, index) in sortedVariants"
      :key="variant.uuid"
      :number="index + 1"
      :variant="variant"
      :listing="listing"
    />
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';
import { SortingOrder } from '@/room';
import sorters from '@/sorters';
import VariantElement from './VariantElement.vue';

@Component({ components: { VariantElement } })
export default class VariantList extends Vue {
  $refs!: { elements: VariantElement[] };

  @Prop(Boolean) canCreate!: boolean;

  @Prop({ required: true })
  order!: SortingOrder;

  @Prop({ type: Boolean, default: false })
  listing!: boolean;

  get sortedVariants() {
    return this.$store.state.variants.slice().sort(sorters.get(this.order));
  }

  public updateSelection(id: string, selected: 'textInput' | 'imageInput') {
    const element = this.$refs.elements.find(e => e.variant.uuid === id);
    if (!element) return;

    element.$refs[selected].select();
  }
}
</script>
