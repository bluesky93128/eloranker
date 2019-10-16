<template>
  <div>
    <div class="section">
      <div class="container">
        <div class="columns is-centered is-vcentered">
          <transition name="fade-variant-left" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[0]" :key="pairVariants[0].uuid">
              <div class="is-overlay is-4by3" @click="vote(pair[0])" :style="{'background-color':activeColor1}">
              </div>
            </VariantElement>
          </transition>
          <div class="column is-one-third" >
            <div class="has-text-centered is-hidden-mobile">
              <p class="title is-1 is-unselectable" id="orLabel">VS</p>
            </div>
            <div class="buttons has-addons is-centered">
              <button class="button is-centered" @click="nextPair">Skip</button>
            </div>
			<div class="buttons has-addons is-centered">
              <input type="checkbox" class="is-centered" v-model="predict_vote" id="predict_check" @input="predict()">
			  <label for="predict_check"> Predictive Voting </label>
            </div>
			<div class="buttons has-addons is-centered">
              <input type="checkbox" class="is-centered" v-model="auto_vote" id="auto_vote">
			  <label for="auto_vote"> Automatic Voting </label>
            </div>
          </div>
          <transition name="fade-variant" mode="out-in">
            <VariantElement v-if="pair" voting :variant="pairVariants[1]" :key="pairVariants[1].uuid">
              <div class="is-overlay is-4by3" @click="vote(pair[1])" :style="{'background-color':activeColor2}">
              </div>
            </VariantElement>
          </transition>
        </div>
      </div>
    </div>
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
import * as events from '@/events';

@Component({ components: { VariantElement, VariantList, ShareBlock } })

export default class RoomVoting extends Vue {
  SortingOrder = SortingOrder;
  
  auto_vote = false;
  predict_vote=false;
  data_set:any[]=[];
  
  activeColor1 = "transparent";
  activeColor2 = 'transparent';

  @State variants!: Variant[];
  pair: [string, string] | null = null;
  showList = false;
  voted_pair=-1;
  handler: any;

  get pairVariants() {
    if (this.pair == null) return [];
    return this.pair.map(id => this.$store.getters.findVariant(id));
  }
  
  created() {
	console.log("created");
  }
  
  mounted() {
	console.log("mounted");
    this.nextPair();
	
	this.handler = (event: events.GetVotingEvent) => {
      if (event.error) {
        if (event.error === 'not enough variants to vote') {
          this.$router.push({ name: 'room-edit', params: this.$route.params });
          return;
        }
        throw new Error(event.error);
      }

      this.pair = event.variants;
	  console.log("asdjfkl;jas;dlkfj;aslkdf");
	  
	  
		
		console.log("voted");
		console.log(this.pairVariants);
		
		if(this.predict_vote) {
			setTimeout(this.predict, 1000);
		}
		this.autovote();
		
    };
	
    connection.on('voting:get', this.handler);
  }

  @Watch('variants')
  onVariantsUpdate() {
    if (this.pairVariants.some(v => v == null)) {
		this.nextPair();
	}  
	
  }
  
  @Watch('predict_vote')
  predict() {

	if(this.pair == null) {
		return;
	}
	var res = this.getPredict(this.pair[0], this.pair[1]);
	console.log(res);
	if(this.predict_vote) {
		if(res == 0) {
			this.activeColor1 = this.predict_vote ? 'RGBA(255, 0, 0, 20%)' : 'transparent';
			this.activeColor2 = 'transparent';
		} else if (res == 1) {
			this.activeColor1 = 'transparent';
			this.activeColor2 = this.predict_vote ? 'RGBA(255, 0, 0, 20%)' : 'transparent';
		} else {
			this.activeColor1 = 'transparent';
			this.activeColor2 = 'transparent';
		}
	} else {
			this.activeColor1 = 'transparent';
			this.activeColor2 = 'transparent';
		}
  }
  
  @Watch('auto_vote')
  async autovote() {
	if(this.auto_vote && this.predict_vote) {
		if(this.pair != null) {
			var res = this.getPredict(this.pair[0], this.pair[1]);
			if(res >= 0) {
				await connection.waitOpen();
				await this.delay(2000);
				this.vote(this.pair[res]);
			}
		}
	}
  }
  
  private delay(ms: number)
  {
	return new Promise(resolve => setTimeout(resolve, ms));
  }

  async nextPair() {
    await connection.waitOpen();
    await connection.getVoting();
  }
  
  genData_set(a:any, b:any){
	this.data_set[a+b] = [a, b];
	this.genPredict(a, b);
  }
  
  genPredict(a:any, b:any){
	
    for(var ind in this.data_set){
        if(this.data_set[ind][1] == a){
            this.data_set[this.data_set[ind][0]+b] = [this.data_set[ind][0], b];
        }
		if(this.data_set[ind][0] == b){
            this.data_set[a+this.data_set[ind][1]] = [a, this.data_set[ind][1]];
        }
    }
	
    console.log(this.data_set);
  }
  
  getPredict(a:string, b:string) {
	console.log("getpredict");
	console.log(a+b);
	for(var ind in this.data_set){
		if(ind==(a+b)) {
			return 0;
		} else if (ind==(b+a)) {
			return 1;
		}
	}
	return -1;
  }

  async vote(id: string) {
    connection.submitVote(id);
    await this.nextPair();
	
	this.activeColor1 = 'transparent';
	this.activeColor2 = 'transparent';
//	this.predict_vote = false;

	if(this.pair != null) {
		if(this.pair[0]==id){
			this.genData_set(this.pair[0],this.pair[1]);
		} else if(this.pair[1]==id) {
			this.genData_set(this.pair[1],this.pair[0]);
		}
	}
  }
  
  destroyed() {
	console.log("destroyed");
	connection.removeVotingListener(this.handler);
  }
  
}
</script>
