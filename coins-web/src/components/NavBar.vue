<template>
  <nav class="navbar is-dark">
    <div class="navbar-start">
      <router-link :class="'navbar-item teh-brand'" :to="{name: 'Home'}">COINS</router-link>
    </div>
    <div class="navbar-end">
      <router-link class="navbar-item" :to="{name: 'Dashboard'}">Dashboard</router-link>
      <router-link class="navbar-item" :to="{name: 'Transactions'}">Transactions</router-link>
      <router-link class="navbar-item" :to="{name: 'Settings'}">Settings</router-link>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'NavBar'
}
</script>

<style lang="stylus">
  .teh-brand
    font-weight: 700
  .navbar-end
    font-size: 1.3rem
  .navbar-item.router-link-active
    background-color: rgba(0, 0, 0, 0.7)
</style>
