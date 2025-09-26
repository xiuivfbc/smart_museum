<template>
  <view class="chat-container">
    <scroll-view class="messages" scroll-y>
      <view class="msg-row"
            v-for="(msg,i) in messages"
            :key="i"
            :class="msg.type">
        <!-- 1. AI 头像在左边 -->
        <image v-if="msg.type==='ai'"
               class="avatar"
               src="/static/touxiang/ai.jpg" mode="aspectFill" />
        <!-- 2. 气泡 -->
        <view class="bubble">
          <view class="bubble-content">{{ msg.content }}</view>
          <view v-if="msg.data" class="bubble-data">
            <text>温度: {{ msg.data.温度 }}°C</text>
            <text>湿度: {{ msg.data.湿度 }}%</text>
          </view>
        </view>
        <!-- 3. 用户头像在右边 -->
        <image v-if="msg.type==='user'"
               class="avatar"
               src="/static/touxiang/user.jpg" mode="aspectFill" />
      </view>
    </scroll-view>

    <view class="input-area">
      <input v-model="inputQuestion" @confirm="sendMessage" placeholder="输入问题…" />
      <button @click="sendMessage">发送</button>
    </view>

    <expert_tabbar />
  </view>
</template>

<script>
import expert_tabbar from 'pages/Common/expert_tabbar.vue';
export default {
  components: { expert_tabbar },
  data() {
    return { inputQuestion: '', messages: [], token: '' };
  },
  onLoad() {
    this.token = uni.getStorageSync('expert_token');
    if (!this.token) uni.redirectTo({ url: '/pages/Expert/expert_login' });
  },
  methods: {
    async sendMessage() {
      const q = this.inputQuestion.trim();
      if (!q) return;
      this.messages.push({ type: 'user', content: q });
      this.inputQuestion  = '';

      try {
        const [, res] = await uni.request({
          url: 'http://127.0.0.1:8081/ai/control',
          method: 'POST',
          data: { question: q },
          header: { Authorization: `Bearer ${this.token}`, 'Content-Type': 'application/json' }
        });
        const { content, data } = res.data;
        this.messages.push({ type: 'ai', content, data });
        this.sendToDevice(data);
      } catch (e) {
        uni.showToast({ title: '获取回答失败', icon: 'none' });
      }
    },
    async sendToDevice(data) {
      await uni.request({
        url: 'http://127.0.0.1:8081/device/autoenv',
        method: 'POST',
        data: { device: 'device1', 温度: data.温度, 湿度: data.湿度 },
        header: { Authorization: `Bearer ${this.token}`, 'Content-Type': 'application/json' }
      });
      uni.showToast({ title: '环境设置完成' });
    }
  }
};
</script>

<style scoped>
/* 整体容器 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 10px 10px 100rpx 10px;
  /* 设置背景图 */
  background-image: url(/static/background/3.jpg);
  background-size: cover;
  background-position: center;
}

/* 聊天区 */
.messages {
  flex: 1;
  overflow-y: auto;
  padding: 10rpx 20rpx;
}

/* 单条消息 */
.msg-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 24rpx;
}
.msg-row.user { justify-content: flex-end; }
.msg-row.ai   { justify-content: flex-start; }

/* 头像圆形（未加载时灰色占位） */
.avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin: 0 20rpx;
}

/* 气泡 */
.bubble {
  max-width: 70%;
  padding: 20rpx 26rpx;
  border-radius: 20rpx;
  position: relative;
  /* 毛玻璃效果相关样式 */
  background-color: rgba(255, 255, 255, 0.1); /* 透明背景 */
  backdrop-filter: blur(8px); /* 模糊效果 */
  -webkit-backdrop-filter: blur(8px); /* 兼容 Safari */
  border: 1px solid rgba(255, 255, 255, 0.2); /* 浅色边框增强质感 */
}
.bubble-content {
  font-size: 30rpx;
  line-height: 1.4;
  color: #333;
}
.bubble-data {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #666;
}

/* 左右气泡颜色微调，更适配背景 */
.msg-row.user .bubble {
  background-color: rgba(0, 122, 255, 0.1); /* 透明背景 */
  color: #fff;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border: 1px solid rgba(0, 122, 255, 0.2);
}
.msg-row.ai .bubble {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

/* 右侧气泡右边距 */
.msg-row.user .bubble { margin-right: 20rpx; }
.msg-row.ai   .bubble { margin-left:  20rpx; }

/* 输入区固定 100rpx，背景半透明更和谐 */
.input-area {
  flex: none;
  height: 100rpx;
  display: flex;
  align-items: center;
  padding: 0 20rpx;
  background-color: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-top: 1rpx solid rgba(238, 238, 238, 0.5);
}
input {
  flex: 1;
  padding: 20rpx;
  border-radius: 40rpx;
  border: 1rpx solid rgba(221, 221, 221, 0.5);
  background-color: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  color: #333;
}
button {
  margin-left: 20rpx;
  padding: 20rpx 30rpx;
  background-color: rgba(0, 122, 255, 0.1);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  color: #fff;
  border-radius: 40rpx;
  border: 1px solid rgba(0, 122, 255, 0.2);
}
</style>