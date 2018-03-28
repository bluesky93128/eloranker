<template>
  <div>
    <section class="hero is-large">
      <div class="hero-body">
        <div class="container">
          <div class="columns is-vcentered">
            <div class="column is-5 is-narrow">
              <div class="content">
                <div class="field">
                  <h1 class="title">
                    Ranking things made easy.
                  </h1>
                  <h4 class="subtitle">
                    Rank a list of items by comparing them 1 vs 1
                  </h4>
                </div>
                <div class="field">
                  <div class="control">
                    <input
                      class="input is-rounded is-large"
                      type="text"
                      placeholder="Type your question here"
                      v-model="roomTitle"
                    >
                  </div>
                </div>
                <nav class="buttons">
                  <a class="button is-primary is-large is-red is-rounded" @click="createRoom">Create Poll</a>
                  <a class="button is-large is-rounded">Try with sample data</a>
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
  roomTitle: string = '';

  async createRoom() {
    const roomTitle = this.roomTitle;
    this.roomTitle = '';
    await connection.waitOpen();
    const { name, secret } = await connection.newRoom(roomTitle);
    this.$router.push({ name: 'room-edit', params: { roomId: `${name}!${secret}` } });
  }
}
</script>
