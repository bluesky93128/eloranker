<template>
  <nav class="navbar">
    <div class="container">
      <div class="navbar-brand">
        <router-link class="navbar-item" to="/">
          <img src="~@/assets/logo.png" alt="Logo">
        </router-link>
        <span class="navbar-burger burger" data-target="navbarMenu">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </div>

      <div class="navbar-end is-hidden-mobile">
        <div class="navbar-item">
          <div class="field is-grouped">
            <p class="control">
              <a class="button is-primary is-red" href="">
                <span class="icon">
                  <i class="icon-pencil-squared icon-5x"></i>
                </span>
                <span>New Poll</span>
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator';

@Component
export default class Header extends Vue {}
</script>

<style lang="scss" module>
.feedback {
  margin-left: 1.25rem;
}

.roomStatus {
  flex: 1;
}
</style>
