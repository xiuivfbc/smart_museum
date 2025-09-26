<template>
  <view class="login-container">
    <text class="welcome">登录</text>

    <view class="tab-group">
      <text :class="['tab-item', activeTab === 'phone' ? 'tab-active' : '']" @click="activeTab = 'phone'">
        手机号登录
      </text>
      <text :class="['tab-item', activeTab === 'email' ? 'tab-active' : '']" @click="activeTab = 'email'">
        邮箱登录
      </text>
    </view>

    <!-- 手机号 -->
    <view v-if="activeTab === 'phone'" class="input-group">
      <text class="input-label">手机号</text>
      <input v-model="userform.phone" class="underline-input" placeholder="请输入手机号" type="number" />
    </view>

    <!-- 邮箱 -->
    <view v-if="activeTab === 'email'" class="input-group">
      <text class="input-label">邮箱</text>
      <input v-model="userform.email" class="underline-input" placeholder="请输入邮箱" type="email" />
    </view>

    <!-- 密码 -->
    <view class="input-group">
      <text class="input-label">密码</text>
      <input v-model="userform.password" class="underline-input" placeholder="请输入密码" password />
    </view>

    <button @click="login" class="login-btn">登录</button>
    <button @click="goRegister" class="register-btn">去注册</button>
  </view>
</template>

<script>
export default {
  data() {
    return {
      activeTab: 'phone',
      userform: { phone: '', email: '', password: '', role: 'admin' },
      loading: false
    };
  },
  methods: {
    async login() {
      if (this.loading) return;
      // 简易前端校验
      const identifier = this.activeTab === 'phone' ? this.userform.phone : this.userform.email;
      if (!identifier || !this.userform.password) {
        uni.showToast({ title: '请输入完整信息', icon: 'none' });
        return;
      }

      this.loading = true;
      const [err, res] = await uni.request({
        url: 'http://127.0.0.1:8081/auth/login',
        method: 'POST',
        data: { identifier, password: this.userform.password, role: this.userform.role },
        header: { 'Content-Type': 'application/json' }
      });

      this.loading = false;
      if (err) {
        uni.showToast({ title: '网络异常', icon: 'none' });
        return;
      }
      if (res.statusCode === 200 && res.data?.token) {
        // 关键：缓存 token 及用户信息
        uni.setStorageSync('expert_token', res.data.token);
        uni.setStorageSync('expert_user', res.data.userform);
console.log('Saved token:', uni.getStorageSync('expert_token')); 
        uni.showToast({ title: '登录成功', icon: 'none' });
        setTimeout(() => uni.redirectTo({ url: '/pages/Expert/expert_home' }), 500);
      } else {
        uni.showToast({ title: res.data?.error || '登录失败', icon: 'none' });
      }
    },
    goRegister() {
      uni.navigateTo({ url: '/pages/Expert/expert_register' });
    }
  }
};
</script>

<style>
	.login-container {
		padding: 60rpx 40rpx;
		background-color: #f8f9fa;
		min-height: 100vh;
		background:
			linear-gradient(rgba(255, 255, 255, 0.5), rgba(255, 255, 255, 0.5)),
			url('/static/background/3.jpg');

	}

	.welcome {
		font-size: 36rpx;
		font-weight: 600;
		color: #2c3e50;
		text-align: center;
		display: block;
		margin-bottom: 40rpx;
	}

	/* 切换标签样式 */
	.tab-group {
		display: flex;
		justify-content: center;
		margin-bottom: 40rpx;
		border-bottom: 1rpx solid #dcdfe6;
	}

	.tab-item {
		padding: 20rpx 40rpx;
		font-size: 28rpx;
		color: #666;
		cursor: pointer;
	}

	.tab-active {
		color: #3498db;
		border-bottom: 2rpx solid #3498db;
	}

	/* 复用注册页面的输入框样式 */
	.input-group {
		margin: 40rpx 0;
		position: relative;
	}

	.input-label {
		font-size: 28rpx;
		color: #34495e;
		margin-bottom: 16rpx;
		display: block;
		font-weight: 500;
	}

	.underline-input {
		border: none;
		border-bottom: 1px solid #dcdfe6;
		padding: 20rpx 0;
		font-size: 30rpx;
		color: #2c3e50;
		background: transparent;
		width: 100%;
		transition: border-color 0.3s ease;
	}

	.underline-input:focus {
		border-bottom-color: #3498db;
		outline: none;
	}

	.underline-placeholder {
		color: #a0a0a0;
		font-size: 28rpx;
	}

	.login-btn,
	.register-btn {
		width: 80%;
		/* 缩小宽度 */
		margin-left: auto;
		margin-right: auto;
		padding: 10rpx 0;
		/* 减小高度 */
		font-size: 28rpx;
		/* 缩小字体 */
		border-radius: 100rpx;
		/* 圆角稍小 */
		margin-top: 80rpx;
		/* 增加与输入框的距离 */
	}

	.login-btn {
		background-color: #4a89b3;
		color: #fff;
		border: none;
	}

	.register-btn {
		margin-top: 30rpx;
		/* 注册按钮与登录按钮间距 */
		background-color: transparent;
		color: #4a89b3;
		border: 1rpx solid #4a89b3;
	}

	/* 按钮悬停效果 */
	.login-btn:active {
		background-color: rgba(52, 152, 219, 0.1);
	}

	.register-btn:active {
		background-color: rgba(52, 152, 219, 0.1);
	}
</style>