<template>
  <div>
    <div :class="['columns', 'is-mobile', 'is-centered']">
      <VariantElement v-if="pair" voting :variant="pairVariants[0]" @click.native="vote(pair[0])" />
      <div :class="['column', 'is-one-third']">
        <button class="button" @click="nextPair">SKIP</button>
      </div>
      <VariantElement v-if="pair" voting :variant="pairVariants[1]" @click.native="vote(pair[1])" />
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';
import { Getter } from 'vuex-class';
import { Variant } from '@/room';
import connection from '@/connection';
import VariantElement from './VariantElement.vue';

@Component({ components: { VariantElement } })
export default class RoomVoting extends Vue {
  @Getter findVariant!: (id: string) => Variant | undefined;
  pair: [string, string] | null = null;

  get pairVariants() {
    if (this.pair == null) return null;
    return this.pair.map(id => this.findVariant(id));
  }

  mounted() {
    this.nextPair();
    connection.on('voting:get', ({ variants, error }) => {
      if (error) {
        if (error === 'not enough variants to vote') {
          this.$router.push({ name: 'room-list', params: this.$route.params });
          return;
        }
        throw new Error(error);
      }

      this.pair = variants;
    });
  }

  async nextPair() {
    await connection.waitOpen();
    connection.getVoting();
  }

  vote(id: string) {
    connection.submitVote(id);
    this.nextPair();
  }
}
</script>

<style lang="scss" module>

</style>
