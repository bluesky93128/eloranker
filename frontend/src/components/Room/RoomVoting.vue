<template>
  <div>
    <div class="section">
      <div class="container">
        <div class="columns is-centered is-vcentered">
          <transition name="fade-variant-left" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[0]" :key="pairVariants[0].uuid">
              <div class="is-overlay is-4by3" @click="vote(pair[0])">
              </div>
            </VariantElement>
          </transition>
          <div class="column is-one-third">
            <div class="has-text-centered is-hidden-mobile">
              <p class="title is-1 is-unselectable" id="orLabel">VS</p>
            </div>
            <div class="buttons has-addons is-centered">
              <button class="button is-centered" @click="nextPair">Skip</button>
            </div>
          </div>
          <transition name="fade-variant" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[1]" :key="pairVariants[1].uuid">
              <div class="is-overlay is-4by3" @click="vote(pair[1])">
              </div>
            </VariantElement>
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
          <div class="columns">
            <div class="column is-half">
              <div class="field has-text-left">
                <h2 class="title is-2">
                  See the results...
                </h2>
              </div>
              <div class="field is-grouped">
                <div class="control">
                  <a class="button is-warning" @click="showList = !showList">
                    <span class="icon">
                      <i :class="showList ? 'icon-eye-off' : 'icon-eye'"></i>
                    </span>
                    <span>{{ showList ? 'Hide' : 'Show' }}</span>
                  </a>
                </div>
              </div>
            </div>
            <div class="column">
              <div class="field has-text-right">
                <h2 class="title is-2">
                  Share
                </h2>
              </div>
              <ShareBlock class="is-pulled-right"/>
            </div>
          </div>
          <VariantList v-if="showList" :order="SortingOrder.RATING" listing />
        </div>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Watch } from 'vue-property-decorator';
import { State } from 'vuex-class';
import { Variant, SortingOrder } from '@/room';
import connection from '@/connection';
import VariantElement from './VariantElement.vue';
import ShareBlock from './ShareBlock.vue';
import VariantList from './VariantList.vue';

@Component({ components: { VariantElement, VariantList, ShareBlock } })
export default class RoomVoting extends Vue {
  SortingOrder = SortingOrder;

  @State variants!: Variant[];
  pair: [string, string] | null = null;
  showList = false;

  get pairVariants() {
    if (this.pair == null) return [];
    return this.pair.map(id => this.$store.getters.findVariant(id));
  }

  mounted() {
    this.nextPair();
    connection.on('voting:get', ({ variants, error }) => {
      if (error) {
        if (error === 'not enough variants to vote') {
          this.$router.push({ name: 'room-edit', params: this.$route.params });
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
