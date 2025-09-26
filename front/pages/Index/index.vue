<template>
  <view class="container">
    <text class="title">请选择您的身份</text>
    <text class="subtitle">选择后进入对应的功能界面</text>

    <view class="card-group">
      <!-- 公众身份卡片 -->
      <view class="role-card" @click="goPublic">
        <image src="/static/icons/public-icon1.jpg" class="card-icon" />
        <text class="card-title">公众用户</text>
        <text class="card-desc">浏览展览、预约参观</text>
      </view>

      <!-- 专家身份卡片 -->
      <view class="role-card expert" @click="goExpert">
        <image src="/static/icons/expert-icon1.jpg" class="card-icon" />
        <text class="card-title">专家用户</text>
        <text class="card-desc">文物管理、数据分析</text>
      </view>
    </view>

  </view>
</template>

<script>
export default {
  methods: {
    goPublic() {
      uni.navigateTo({
        url: '/pages/Public/public_home'
      });
    },
    goExpert() {
      uni.navigateTo({
        url: '/pages/Expert/expert_login'
      });
    }
  }
}
</script>

<style>
.container {
  padding: 40rpx;
  background:
  	linear-gradient(rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.1)),
  	url('/static/background/3.jpg');
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.title {
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 16rpx;
  text-shadow: 0 2rpx 4rpx rgba(0,0,0,0.1);
}

.subtitle {
  font-size: 28rpx;
  color: rgba(255,255,255,0.8);
  margin-bottom: 80rpx;
}

.card-group {
  display: flex;
  flex-direction: column;
  gap: 40rpx;
  width: 100%;
  max-width: 500rpx;
}

.role-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 10rpx 30rpx rgba(0,0,0,0.1);
  transition: all 0.3s ease;
}

.role-card:active {
  transform: scale(0.95);
  box-shadow: 0 5rpx 15rpx rgba(0,0,0,0.1);
}

.role-card.expert {
  background: #f8f9fa;
}

.card-icon {
  width: 120rpx;
  height: 120rpx;
  margin-bottom: 24rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 12rpx;
}

.card-desc {
  font-size: 24rpx;
  color: #666;
  text-align: center;
}

.tip {
  margin-top: 60rpx;
  font-size: 24rpx;
  color: rgba(255,255,255,0.7);
}
</style>