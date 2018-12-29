<template>
  <transition enter-active-class="animated zoomIn" leave-active-class="animated zoomOut">
    <div v-show="show" class="wrap">
      <div class="bg">
        <div class="title">
          <div class="title-inner">
            <img class="title-icon" :src="icon" />
            <h2 class="title-text">{{title}}</h2>
          </div>
        </div>
        <div class="confirm-button" @click="click_button">
          <span class="confirm-button-text">{{button_text}}</span>
        </div>
      </div>
    </div>
  </transition>
</template>
<script>
export default {
  name: 'dialog',
  props: {
    'title': {
      default: ''
    },
    'button_text': {
      default: 'OK'
    },
    'show': {
      default: false
    },
    'icon': {
      default: require('../assets/trophy.png')
    }
  },
  methods: {
    click_button: function () {
      this.show = false
      this.$emit('close')
    }
  },
  data () {
    return {}
  }
}
</script>
<style scoped>
.wrap {
  width: 100%;
  height: 100%;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  position: fixed;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
}

.bg {
  width: 400px;
  height: 208px;
  background-color: #4abdac;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-radius: 10px;
  border: 1px solid #4c4c4c;
  box-shadow: 0px 5px 10px #0e0c0d;
}

.title {
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
}

.title-inner {
  flex: 0 0 75%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.title-icon {
  flex: 0 0 15%;
  width: 15%;
  height: auto;
  margin-right: 10px;
}

.title-text {
  display: inline;
  color: white;
}

.confirm-button {
  height: 54px;
  width: 50%;
  background-color: #f78733;
  border-radius: 8px;
  margin-top: 6.5%;
  margin-left: 0.5%;
  margin-right: 0.5%;
  margin-bottom: 5px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.confirm-button-text {
  color: white;
  font-size: 18px;
  font-weight: bold;
}

</style>
