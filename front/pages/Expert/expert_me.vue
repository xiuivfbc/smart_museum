<template>
  <view class="me-container">
    <view v-if="isLogin">
      <text class="welcome">欢迎，{{ userform.username }}</text>

      <!-- 个人信息卡片 -->
      <uni-card title="个人信息管理" is-full>
        <view class="grid-container">
          <view class="grid-item" @click="goProfileDetail">
            <image class="grid-icon" src="/static/icons/profile.png" mode="aspectFit"></image>
            <text class="grid-text">个人资料</text>
          </view>
          <view class="grid-item" @click="goEditProfile">
            <image class="grid-icon" src="/static/icons/edit.png" mode="aspectFit"></image>
            <text class="grid-text">更新信息</text>
          </view>
          <view class="grid-item" @click="goChangePassword">
            <image class="grid-icon" src="/static/icons/password.png" mode="aspectFit"></image>
            <text class="grid-text">修改密码</text>
          </view>
          <view class="grid-item" @click="logout">
            <image class="grid-icon" src="/static/icons/logout.png" mode="aspectFit"></image>
            <text class="grid-text">注销登录</text>
          </view>
        </view>
      </uni-card>

     

      <!-- 自定义 tabBar -->
      <expert_tabbar />
    </view>
  </view>
</template>

<script>
import expert_tabbar from '../Common/expert_tabbar.vue';

export default {
  components: {
    expert_tabbar
  },
  data() {
    return {
      isLogin: false,
      userform: {
        id: "",
        username: "",
        phone: "",
        email: "",
        password: "",
        role: "user"
      }
    }
  },
  onShow() {
    const token = uni.getStorageSync("expert_token")
    const userInfo = uni.getStorageSync("expert_user")
    this.userform = userInfo
    if (token) {
      this.isLogin = true
      this.userform = userInfo
    } else {
      uni.redirectTo({
        url: "/pages/Expert/expert_login"
      })
    }
  },
  methods: {
    logout() {
      uni.removeStorageSync("expert_token")
      uni.removeStorageSync("expert_user")
      this.isLogin = false

      uni.showToast({
        title: "已退出登录",
        icon: "none"
      })

      uni.redirectTo({
        url: "/pages/Expert/expert_login"
      })
    },
    goProfileDetail() {
      uni.navigateTo({
        url: "/pages/Expert/expert_profile"
      })
    },
    goEditProfile() {
      uni.navigateTo({
        url: "/pages/Expert/expert_edit_profile"
      })
    },
    goChangePassword() {
      uni.navigateTo({
        url: "/pages/Expert/expert_change_password"
      })
    }
    
  }
}
</script>

<style lang="scss">
.me-container {
  padding: 40rpx 30rpx;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.welcome {
  font-size: 36rpx;
  font-weight: 600;
  color: #2c3e50;
  text-align: center;
  display: block;
  margin-bottom: 40rpx;
}

/* 网格布局 */
.grid-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  padding: 20rpx 0;
}

.grid-item {
  width: 23%; /* 留出一些间距 */
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 30rpx;
  padding: 20rpx 0;
  border-radius: 12rpx;
  background-color: #fff;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.grid-item:active {
  background-color: #f5f5f5;
  transform: scale(0.98);
}

.grid-icon {
  width: 60rpx;
  height: 60rpx;
  margin-bottom: 15rpx;
}

.grid-text {
  font-size: 24rpx;
  color: #333;
  text-align: center;
}

.logout-text {
  color: #e74c3c;
}

/* 卡片间距 */
.uni-card {
  margin-bottom: 30rpx;
}
</style>