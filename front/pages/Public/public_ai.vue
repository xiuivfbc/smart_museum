<template>
  <view class="container">
    <!-- web-view 占剩余高度 -->
    <web-view
      src="http://154.39.79.242:8080/chat/eae0c9f1e7c3b42c?mode=mobile"
      style="width: 100%; height: calc(100vh - 100rpx);"
      @message="onWebMessage"
    ></web-view>

    <!-- 底部 tabBar -->
    <public_tabbar />
  </view>
</template>

<script>
import public_tabbar from 'pages/Common/public_tabbar.vue';

export default {
  components: { public_tabbar },
  onLoad() {
    window.addEventListener('message', this.onWebMessage);
  },
  methods: {
    onWebMessage(e) {
      if (e.data === 'back') {
        uni.switchTab({ url: '/pages/Public/public_home' });
      }
    }
  }
};
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  height: 100vh;
}
/* web-view 占剩余高度，底部 100rpx 留给 tabBar */
</style>