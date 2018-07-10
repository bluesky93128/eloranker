<template>
  <div class="import-block">
    <h3 class="h3">Import from GitHub</h3>
    <input class="input" placeholder="Owner" v-model="owner">
    <input class="input" placeholder="Repository" v-model="repo">
    <button class="button" @click="doImport">Import</button>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import connection from '@/connection';

@Component
export default class ImportBlock extends Vue {
  owner = '';
  repo = '';
  doImport() {
    connection.importVariants(this.owner, this.repo);
  }
}
</script>

<style lang="scss" module>

</style>
