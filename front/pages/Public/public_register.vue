<template>
  <view class="register-container">
    <text class="welcome">注册</text>

    <!-- 邮箱 -->
    <view class="input-group">
      <text class="input-label">邮箱</text>
      <input v-model="userform.email" class="underline-input" placeholder="请输入邮箱"
        placeholder-class="underline-placeholder" type="email" />
    </view>

    <!-- 手机号（选填） -->
    <view class="input-group">
      <text class="input-label">手机号（选填）</text>
      <input v-model="userform.phone" class="underline-input" placeholder="请输入手机号"
        placeholder-class="underline-placeholder" type="number" />
    </view>

    <!-- 验证码 + 发送按钮 -->
    <view class="input-group">
      <text class="input-label">验证码</text>
      <view class="input-row">
        <input v-model="verificationCode" class="underline-input flex-1" placeholder="请输入验证码"
          placeholder-class="underline-placeholder" type="text" />
        <button @click="sendVerificationCode" class="send-btn-mini">发送</button>
      </view>
    </view>

    <!-- 密码 -->
    <view class="input-group">
      <text class="input-label">密码</text>
      <input v-model="userform.password" class="underline-input" placeholder="密码至少六位"
        placeholder-class="underline-placeholder" password />
    </view>

    <!-- 确认密码 -->
    <view class="input-group">
      <text class="input-label">确认密码</text>
      <input v-model="confirmPassword" class="underline-input" placeholder="请再次输入密码"
        placeholder-class="underline-placeholder" password />
    </view>

    <button @click="register" class="register-btn" :class="{'btn-loading': loading}" :disabled="loading">
      <text v-if="!loading">注册</text>
      <text v-else>注册中...</text>
    </button>

    <public_tabbar />
  </view>
</template>

<script>
import public_tabbar from '../Common/public_tabbar.vue';

export default {
  components: { public_tabbar },
  data() {
    return {
      userform: {
        email: '',
        phone: '',
        password: '',
        role: 'user'
      },
      confirmPassword: '',
      verificationCode: '',       // 用户手动输入
      serverCode: '',             // 前端从服务端拿到的验证码
      loading: false
    };
  },
  methods: {
    /* 发送验证码 */
    async sendVerificationCode() {
      if (!this.userform.email) {
        uni.showToast({ title: '邮箱不能为空', icon: 'none' });
        return;
      }

      try {
        const [err, res] = await uni.request({
          url: 'http://127.0.0.1:8081/auth/verificationCode',
          method: 'POST',
          data: { email: this.userform.email, phone: this.userform.phone },
          header: { 'Content-Type': 'application/json' },
          timeout: 3000
        });

        if (err) {
          uni.showToast({ title: '请求超时', icon: 'none' });
          return;
        }

        if (res.statusCode === 200) {
          this.serverCode = res.data.code || ''; // 保存验证码
          console.log('验证码已收到：', this.serverCode);
          uni.showToast({ title: '验证码已发送', icon: 'none' });
        } else {
          uni.showToast({ title: res.data.message || '发送失败', icon: 'none' });
        }
      } catch (e) {
        uni.showToast({ title: '网络异常', icon: 'none' });
      }
    },

    /* 注册：先比对验证码再提交 */
    async register() {
      if (!this.userform.email) {
        uni.showToast({ title: '邮箱不能为空', icon: 'none' });
        return;
      }
      if (!this.userform.password || this.userform.password.length < 6) {
        uni.showToast({ title: '密码至少六位', icon: 'none' });
        return;
      }
      if (this.userform.password !== this.confirmPassword) {
        uni.showToast({ title: '两次密码不一致', icon: 'none' });
        return;
      }
      if (!this.verificationCode) {
        uni.showToast({ title: '请填写验证码', icon: 'none' });
        return;
      }
      if (this.verificationCode !== this.serverCode) {
        uni.showToast({ title: '验证码错误', icon: 'none' });
        return;
      }

      this.loading = true;
      try {
        const [err, res] = await uni.request({
          url: 'http://127.0.0.1:8081/auth/register',
          method: 'POST',
          data: {
            email: this.userform.email,
            password: this.userform.password,
            phone: this.userform.phone,
            verificationCode: this.verificationCode
          },
          header: { 'Content-Type': 'application/json' },
          timeout: 3000
        });

        if (err) {
          uni.showToast({ title: '请求超时', icon: 'none' });
          return;
        }

        if (res.statusCode === 200) {
          uni.setStorageSync("public_user", { email: this.userform.email });
          uni.showToast({ title: '注册成功', icon: 'none' });
          setTimeout(() => uni.redirectTo({ url: "/pages/Public/public_login" }), 500);
        } else {
          uni.showToast({ title: res.data.message || '注册失败', icon: 'none' });
        }
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.register-container {
  padding: 60rpx 40rpx;
  background-color: #f8f9fa;
  min-height: 100vh;
}
.welcome {
  font-size: 36rpx;
  font-weight: 600;
  color: #2c3e50;
  text-align: center;
  margin-bottom: 60rpx;
}
.input-group {
  margin: 40rpx 0;
}
.input-label {
  font-size: 28rpx;
  color: #34495e;
  margin-bottom: 16rpx;
}
.underline-input {
  border: none;
  border-bottom: 1px solid #dcdfe6;
  padding: 20rpx 0;
  width: 100%;
}
.input-row {
  display: flex;
  align-items: center;
}
.flex-1 {
  flex: 1;
}
.send-btn-mini {
  background-color: #4a89b3;
  color: #fff;
  border: none;
  border-radius: 4rpx;
  padding: 10rpx 20rpx;
  font-size: 28rpx;
  margin-left: 20rpx;
}
.register-btn {
  width: 80%;
  margin: 80rpx auto 0;
  padding: 10rpx 0;
  font-size: 28rpx;
  border-radius: 100rpx;
  background-color: #4a89b3;
  color: #fff;
  border: none;
}
</style>