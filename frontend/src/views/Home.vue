<template>
  <div>
    <section class="hero is-medium intro">
      <div class="hero-body">
        <div class="container">
          <div class="intro-columns">
            <div class="intro-column is-content">
              <div class="intro-content">
                <h1 class="title">
                  Ranking things made easy.
                </h1>
                <nav class="buttons">
                  <a class="button is-primary is-large is-red" @click="createRoom">Create Poll</a>
                  <a class="button is-large">How it works</a>
                </nav>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
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
