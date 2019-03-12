<template>
  <div id="wm-loader">
    <h1>Watermark creator</h1>
    <div class="form" v-on:change="handleFileChange()" >
      <div class="item">
        <label for="imageFile">Image File</label>
        <input id="imageFile" type="file" accept="image/jpeg" />
      </div>
      <div class="item">
        <label for="watermarkFile">Watermark File</label>
        <input id="watermarkFile" type="file" accept="image/x-png" />
      </div>

      <div class="control">
        <button v-on:click="upload()">Send</button>
      </div>
    </div>
  </div>
</template>


<script>
  let axios = require("axios");
  let fileDownload = require('js-file-download');

  export default {
    name: 'WatermarkLoader',
    data() {
      return {
        images: {}
      }
    },
    methods: {
      handleFileChange() {
        console.log(this);
        this.$el.querySelectorAll("input[type='file']").forEach(f => {
          this.images[f.id] = f.files[0];
        });
      },
      upload() {
        let formData = new FormData();

        Object.keys(this.images).forEach(k => {
          formData.append(k, this.images[k]);
        });

        axios.post(
          "",
          formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            },
          responseType: 'blob'
        }).then(response => {
          console.log('Sended');
          fileDownload(response.data, 'result.jpg');
         }).catch(response => {
          console.log('Failed');
        });

      }
    }
  }
</script>


<style scoped>
  #wm-loader {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .form {
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

