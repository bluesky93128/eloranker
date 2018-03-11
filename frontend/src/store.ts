// TODO: https://github.com/vuejs/vuex/pull/1121

import Vue from 'vue';
import router from '@/router';
import Vuex, { GetterTree, MutationTree, ActionTree } from 'vuex';
import { Variant, SortingOrder, emptyVariant, EditMode } from './room';
import connection from './connection';
import sorters from './sorters';
import { JoinRoomEvent } from '@/events';

Vue.use(Vuex);

interface State {
  joined: boolean;
  isAdmin: boolean;
  roomName: string;
  roomSecret: string;

  roomTitle: string;
  roomQuotaEnabled: boolean;
  roomEditMode: EditMode;

  identifier: string;
  ignoredVariants: { [id: string]: boolean };
  clientNumber: number;
  variants: Variant[];
  sortingOrder: SortingOrder;
}

const baseState: State = {
  joined: false,
  isAdmin: false,
  roomName: '',
  roomSecret: '',

  roomTitle: '',
  roomQuotaEnabled: false,
  roomEditMode: EditMode.Normal,

  identifier: '',
  ignoredVariants: {},
  clientNumber: 0,
  variants: [],
  sortingOrder: SortingOrder.DATE,
};

const getters: GetterTree<State, any> = {
  findVariant: state => (id: string) => state.variants.find(v => v.uuid === id),
  isIgnoredVariant: state => (id: string) => state.ignoredVariants[id],
  canIgnoreVariant: state => {
    if (state.variants.length < 6) return false;

    const { ignoredVariants } = state;
    const ignoredVariantsLen = Object.keys(ignoredVariants).filter(k => ignoredVariants[k]).length;
    return ignoredVariantsLen * 2 < state.variants.length;
  },

  canVote: state => state.variants.length >= 2,
  sortedVariants: state => state.variants.slice().sort(sorters.get(state.sortingOrder)),
  hasWriteAccess: (state, getters) => (id: string) => {
    switch (state.roomEditMode) {
      case EditMode.Trust:
        return true;

      case EditMode.Restricted:
        return state.isAdmin;

      case EditMode.Normal:
        if (state.isAdmin) return true;
        if (id === '') return true;

        const variant: Variant | null = getters.findVariant(id);
        if (variant == null) return false;

        return variant.author === state.identifier;

      default:
        return false;
    }
  },
};

const mutations: MutationTree<State> = {
  setJoined: (state, joined: boolean) => (state.joined = joined),
  setClientNumber: (state, clients: number) => (state.clientNumber = clients),

  createVariant(state, init: Partial<Variant>) {
    state.variants.push(emptyVariant(init));
  },

  updateVariant(state, variant: Partial<Variant> & { uuid: string }) {
    const localVariant = state.variants.find(v => v.uuid === variant.uuid);
    if (localVariant == null) {
      window.location.reload();
      return;
    }

    Object.assign(localVariant, variant);
  },

  recalculateRoomId(state, roomId: string) {
    const [name, secret] = roomId.split(/!/);
    state.roomName = name;
    state.roomSecret = secret;
  },

  loadRoom(state, event: JoinRoomEvent) {
    state.isAdmin = event.isAdmin;
    state.variants = event.variants;
    state.identifier = event.identifier;
    state.ignoredVariants = event.ignoredVariants;

    store.commit('setTitle', event.title);
    store.commit('setQuotaEnabled', event.quotaEnabled);
    store.commit('setEditMode', event.editMode);
  },

  setVariantIgnored(state, { id, ignored }: { id: string; ignored: boolean }) {
    Vue.set(state.ignoredVariants, id, ignored);
  },

  setSortingOrder(state, value: SortingOrder) {
    state.sortingOrder = value;
  },

  setTitle(state, value: string) {
    state.roomTitle = value;
  },
  setQuotaEnabled(state, value: boolean) {
    state.roomQuotaEnabled = value;
  },
  setEditMode(state, value: EditMode) {
    state.roomEditMode = value;
  },
};

const actions: ActionTree<State, any> = {
  async joinRoom(store) {
    await connection.waitOpen();
    try {
      const result = await connection.joinRoom(store.state.roomName, store.state.roomSecret);
      store.commit('loadRoom', result);
      store.commit('setJoined', true);
    } catch (err) {
      store.commit('setJoined', false);
      if (err.message === 'room not exists') {
        router.push({ name: 'home' });
      } else {
        throw err;
      }
    }
  },
};

const store = new Vuex.Store<State>({
  actions,
  mutations,
  getters,
  state: baseState,
});

connection.on('room:clients', ({ clients }) => store.commit('setClientNumber', clients));
connection.on('variant:allocate', init => store.commit('createVariant', init));
connection.on('variant:update', event => {
  if (event.error) throw new Error(event.error);
  store.commit('updateVariant', event);
});

connection.on('settings:title', ({ error, value }) => {
  if (error) throw new Error(error);
  store.commit('setTitle', value);
});
connection.on('settings:quotaEnabled', ({ error, value }) => {
  if (error) throw new Error(error);
  store.commit('setQuotaEnabled', value);
});
connection.on('settings:editMode', ({ error, value }) => {
  if (error) throw new Error(error);
  store.commit('setEditMode', value);
});

export default store;
