<template>
  <div id="app">
    <WatermarkLoader />
  </div>
</template>


<script>
  import WatermarkLoader from './components/WatermarkLoader.vue'

  export default {
    name: 'app',
    components: {
      WatermarkLoader
    }
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;

    color: #2c3e50;

    display: flex;
    justify-content: center;
    margin-top: 60px;
    text-align: center;
  }
</style>
