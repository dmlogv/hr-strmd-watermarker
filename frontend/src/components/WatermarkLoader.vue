<template>
  <div id="wm-loader">
    <h1>Watermark creator</h1>
    <form @submit="checkForm" action="/" method="post">
      <div class="item">
        <label for="imageFile">Image File</label>
        <input id="imageFile" type="file" accept="image/jpeg" />
      </div>
      <div class="item">
        <label for="watermarkFile">Watermark File</label>
        <input id="watermarkFile" type="file" accept="image/x-png" />
      </div>
    </form>
  </div>
</template>


<script>
  export default {
    name: 'WatermarkLoader'
  }
</script>


<style scoped>
  #wm-loader {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  form {
    display: table;
    border-spacing: 1em .5em
  }

  .item {
    display: table-row;
  }

  .item label,
  .item input {
    display: table-cell;
  }

  .item label {
    text-align: right;
  }

  .item input {
    border: 1px solid #ddd;
  }
</style>

