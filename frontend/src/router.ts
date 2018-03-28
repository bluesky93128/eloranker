import Vue from 'vue';
import Router from 'vue-router';
import Home from '@/views/Home.vue';
import Room from '@/views/Room.vue';
import RoomList from '@/components/Room/RoomList.vue';
import RoomVoting from '@/components/Room/RoomVoting.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  linkExactActiveClass: 'is-active',
  routes: [
    { path: '/', component: Home, name: 'home' },
    {
      path: '/:roomId',
      component: Room,
      children: [
        { path: '/', component: RoomVoting, name: 'room-voting' },
        { path: 'edit', component: RoomList, name: 'room-edit' },
      ],
    },
  ],
});
