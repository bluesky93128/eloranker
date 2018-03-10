<template>
  <div :class="$style.header">
    <!-- <router-link to="/">HOME</router-link> -->
    <div :class="$style.feedback">
      <a href="https://discord.gg/pz76E6D">Feedback (Discord server)</a>
    </div>
    <RoomStatus :class="$style.roomStatus" v-if="joined" />
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { State } from 'vuex-class';
import RoomStatus from '@/components/Header/RoomStatus.vue';

@Component({ components: { RoomStatus } })
export default class Header extends Vue {
  @State joined!: boolean;
}
</script>

<style lang="scss" module>
.header {
  display: flex;
  flex-flow: row;
  background-color: azure;
}

.feedback {
  margin-left: 1.25rem;
}

.roomStatus {
  flex: 1;
}
</style>
