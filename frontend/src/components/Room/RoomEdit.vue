<template>
  <div class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-four-fifths">
          <VariantList :order="this.$store.state.sortingOrder" canCreate />
        </div>
        <RoomSettings/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import VariantElement from './VariantElement.vue';
import VariantList from './VariantList.vue';
import RoomSettings from './RoomSettings.vue';

@Component({ components: { VariantElement, VariantList, RoomSettings } })
export default class RoomEdit extends Vue {}
</script>
