<template>
  <div>
    <div class="section">
      <div class="container">
        <div class="columns is-mobile is-centered is-vcentered">
          <transition name="fade-variant-left" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[0]" :key="pairVariants[0].uuid" @click.native="vote(pair[0])" />
          </transition>
          <div class="column is-one-third">
            <div class="has-text-centered">
              <p class="title is-1 is-unselectable" id="orLabel">VS</p>
            </div>
            <div class="buttons has-addons is-centered">
              <button class="button is-centered" @click="nextPair">Skip</button>
            </div>
          </div>
          <transition name="fade-variant" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[1]" :key="pairVariants[1].uuid" @click.native="vote(pair[1])" />
          </transition>
        </div>
      </div>
    </div>
    <!-- <article class="message">
      <div class="message-header">
        <p>Hello World</p>
        <button class="delete" aria-label="delete"></button>
      </div>
      <div class="message-body">
        Lorem ipsum dolor sit amet, consectetur adipiscing elit. <strong>Pellentesque risus mi</strong>, tempus quis placerat ut, porta nec nulla. Vestibulum rhoncus ac ex sit amet fringilla. Nullam gravida purus diam, et dictum <a>felis venenatis</a> efficitur. Aenean ac <em>eleifend lacus</em>, in mollis lectus. Donec sodales, arcu et sollicitudin porttitor, tortor urna tempor ligula, id porttitor mi magna a neque. Donec dui urna, vehicula et sem eget, facilisis sodales sem.
      </div>
    </article> -->
    <section class="hero is-primary is-red">
      <div class="hero-body">
        <div class="container">
          <h2 class="title is-2">See the rankings...</h2>
          <div class="columns is-multiline is-mobile">
            <VariantElement
              ref="elements"
              v-for="(variant, index) in sortedVariants"
              :number="index + 1"
              :key="variant.uuid"
              :variant="variant"
            />
          </div>
        </div>
      </div>
    </section>
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
  @Getter sortedVariants!: Variant[];
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
