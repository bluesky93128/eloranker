<template>
  <div v-if="joined">
    <RoomStatus />
    <router-view />
  </div>
</template>

<script lang="ts">
import { Vue, Component, Watch } from 'vue-property-decorator';
import { State, Getter } from 'vuex-class';
import { Route } from 'vue-router/types/router';
import RoomStatus from '@/components/Room/RoomStatus.vue';

@Component({ components: { RoomStatus } })
export default class RoomList extends Vue {
  @State roomName!: string;
  @State roomSecret!: string;
  @State joined!: boolean;
  @Getter canVote!: boolean;

  mounted() {
    this.joinRoom(this.$route);
  }

  destroyed() {
    this.$store.dispatch('leaveRoom');
  }

  @Watch('$route')
  private async joinRoom(_to: Route, from?: Route) {
    if (from != null && from.name != null && from.name.startsWith('room-')) {
      return;
    }

    this.$store.dispatch('joinRoom');
  }
}
</script>
