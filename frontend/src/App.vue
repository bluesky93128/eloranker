<template>
  <div id="app">
    <Header />
    <router-view />
    <div v-if="error" :class="$style.errorOverlay">
      <div :class="$style.errorOverlayMessage">{{ error }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import connection from '@/connection';
import Header from './components/Header/Header.vue';

@Component({ components: { Header } })
export default class App extends Vue {
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
:global #app {
  text-align: center;
  color: #2c3e50;
  min-width: 200px;
  width: 100%;
  height: 100%;
  display: flex;
  flex-flow: column;
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
