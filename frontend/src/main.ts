import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './bulma.scss';

router.afterEach(to =>
  store.commit('recalculateRoomId', to.params.roomId != null ? to.params.roomId : ''),
);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
