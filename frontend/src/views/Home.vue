<template>
  <div>
    <button class="button is-link" @click="createRoom">CREATE</button>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import connection from '@/connection';

@Component
export default class Home extends Vue {
  async createRoom() {
    await connection.waitOpen();
    const { name, secret } = await connection.newRoom();
    this.$router.push({ name: 'room-list', params: { roomId: `${name}!${secret}` } });
  }
}
</script>

<style lang="scss" module>

</style>
