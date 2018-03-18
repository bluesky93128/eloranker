<template>
  <div>
    <div class="columns is-mobile is-centered">
      <transition name="fade-variant-left" mode="out-in">
        <VariantElement v-if="pair" voting :variant="pairVariants[0]" :key="pairVariants[0].uuid" @click.native="vote(pair[0])" />
      </transition>
      <div class="column is-one-third">
        <div class="buttons has-addons is-centered">
          <button class="button is-centered" @click="nextPair">SKIP</button>
        </div>
      </div>
      <transition name="fade-variant" mode="out-in">
        <VariantElement v-if="pair" voting :variant="pairVariants[1]" :key="pairVariants[1].uuid" @click.native="vote(pair[1])" />
      </transition>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Watch } from 'vue-property-decorator';
import { State, Getter } from 'vuex-class';
import { Variant } from '@/room';
import connection from '@/connection';
import VariantElement from './VariantElement.vue';

@Component({ components: { VariantElement } })
export default class RoomVoting extends Vue {
  @State variants!: Variant[];
  @Getter findVariant!: (id: string) => Variant | undefined;
  pair: [string, string] | null = null;

  get pairVariants() {
    if (this.pair == null) return [];
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

  @Watch('variants')
  onVariantsUpdate() {
    if (this.pairVariants.some(v => v == null)) {
      this.nextPair();
    }
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
