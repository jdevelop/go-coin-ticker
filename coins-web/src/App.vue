<template>
  <div id="app">
    <nav-bar/>
    <div class="container fluid root-container">
      <router-view/>
    </div>
  </div>
</template>

<script>
import NavBar from './components/NavBar'
export default {
  name: 'app',
  components: {
    NavBar
  }
}
</script>

<style lang="sass">
  // 1. Import the initial variables
  @import "~bulma/sass/utilities/initial-variables"
  @import "~bulma/sass/utilities/functions"

  // 2. Set your own initial variables
  @import "./styles/overrides.sass"

  // 5. Add new color variables to the color map.
  @import "~bulma/sass/utilities/derived-variables.sass"
  // $addColors: ("twitter":($twitter, $twitter-invert), "linkedin": ($linkedin, $linkedin-invert), "github": ($github, $github-invert))
  // $colors: map-merge($colors, $addColors)

  // 6. Import the rest of Bulma
  @import "~bulma/bulma"

  // SHARED STYLES
  @import "./styles/shared.sass"

  html
    // font-size: 62.5%
    // font-size: 125%
  html, body
    min-height: 100%
  #app
    font-size: 1.6em
    min-height: 100vh
    font-family: 'Open Sans', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased
    -moz-osx-font-smoothing: grayscale

</style>
