<template>
  <div id="app">
    <Header />
    <div v-if="joined">
      <RoomStatus :class="$style.roomStatus" />
    </div>
    <section class="section">
      <div class="container">
        <router-view />
      </div>
    </section>
    <div v-if="error" :class="$style.errorOverlay">
      <div :class="$style.errorOverlayMessage">{{ error }}</div>
    </div>
    <Footer />
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { State } from 'vuex-class';
import connection from '@/connection';
import Header from './components/Header/Header.vue';
import Footer from './components/Footer/Footer.vue';
import RoomStatus from '@/components/Header/RoomStatus.vue';

@Component({ components: { Header, Footer, RoomStatus } })
export default class App extends Vue {
  @State joined!: boolean;

  socketClosed = false;

  mounted() {
    connection.on('state', state => {
      if (state >= WebSocket.CLOSING) {
        this.socketClosed = true;
      }
    });
    connection.on('error', () => (this.socketClosed = true));
  }

  get error() {
    if (WebSocket == null) {
      return 'WebSocket is required for application to work. Update your browser';
    }
    if (this.socketClosed) return 'Socket is closed';
  }
}
</script>

<style lang="scss" module>
@import url('https://fonts.googleapis.com/css?family=Fira+Sans|Fira+Sans+Condensed:400,500');

:global #app {
  background: #f5f5f5;
}

.errorOverlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(black, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
}

.errorOverlayMessage {
  margin: auto;
  color: white;
  font-size: 2rem;
}
</style>
