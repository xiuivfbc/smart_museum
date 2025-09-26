<template>
  <view class="card-detail">
    <image :src="cardDetail.cover" mode="aspectFit" class="card-detail-cover"></image>
    <view class="card-detail-info">
      <text class="card-detail-title">{{ cardDetail.title }}</text>
      <text class="card-detail-subtitle">藏品信息</text>
      <text class="card-detail-description">所属年代: {{ cardDetail.age }}</text>
      <text class="card-detail-subtitle">藏品介绍</text>
      <text class="card-detail-description">{{ cardDetail.description }}</text>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      cardDetail: {}
    };
  },
  onLoad(options) {
    const index = parseInt(options.index);
    const cards = [
      {
        cover: '/static/cards/4.jpg',
        title: '金缕玉衣',
        age: '汉',
        description: '金缕玉衣是汉代皇帝和高级贵族死时穿用的殓服，外观和人体形状相同，由玉片拼成，并用金丝加以编缀。这种玉衣不仅只是简单磨成玉片而已，上面还雕有花纹，反映了玉师杰出的技艺和达官奢侈的生活。'
      },
      {
        cover: '/static/cards/6.jpg',
        title: '乳钉纹青铜爵',
        age: '夏',
        description: '乳钉纹青铜爵是一件来自二里头文化的文物，以其俊巧的艺术美感和丰富的历史内涵著称。它不仅是夏代青铜冶铸技术的实物见证，更是贵族们身份和地位的象征,其铸造工艺堪称夏代青铜铸造技术的巅峰之作.'
      },
      {
        cover: '/static/cards/7.jpg', // 替换为你的图片路径
        title: '青花寿山福海纹瓷香炉',
        age: '明',
        description: '造型为青铜鼎样式，盘口，方唇，束颈，圆腹，小平底。纹样以苏麻离青料绘蓝地白花寿山福海纹，青花发色浓重。'
      }
    ];

    if (index !== undefined && index >= 0 && index < cards.length) {
      this.cardDetail = cards[index];
    } else {
      console.error('Invalid index:', index);
      // 处理无效的index，例如显示错误信息或默认内容
    }
  }
}
</script>

<style scoped>
/* ========== 毛玻璃详情页 ========== */
:root {
  --primary: #8c6f58;
  --glass-bg: rgba(245, 245, 245, 0.8); /* 更换背景颜色 */
  --glass-border: rgba(0, 0, 0, 0.1);
  --text-dark: #2e2e2e;
  --text-light: #555;
  --radius: 28rpx;
  --shadow: 0 12rpx 48rpx rgba(0, 0, 0, 0.08);
}

.card-detail {
  min-height: 100vh;
  padding: 40rpx 30rpx 120rpx;
  background: url('/static/background/detail-bg.jpg') center/cover no-repeat;
  display: flex;
  flex-direction: column;
  align-items: center; /* 垂直居中对齐 */
}

/* 图片 */
.card-detail-cover {
  width: 90%; /* 宽度为90%，不撑满页面 */
  height: 40vh; /* 占据上半部分五分之二 */
  border-radius: var(--radius);
  overflow: hidden;
  transition: transform .35s;
}
.card-detail-cover:active {
  transform: scale(1.02);
}

/* 信息区域 */
.card-detail-info {
  flex: 1;
  padding: 40rpx;
  background: var(--glass-bg);
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1rpx solid var(--glass-border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 90%; /* 宽度为90%，不撑满页面 */
}

/* 标题 */
.card-detail-title {
  font-size: 44rpx;
  font-weight: 700;
  color: var(--text-dark);
  margin: 30rpx 0 20rpx;
  text-align: center;
  font-family: KaiTi, serif; /* 使用楷体 */
}

/* 子标题 */
.card-detail-subtitle {
  font-size: 36rpx;
  font-weight: 600;
  color: var(--text-dark);
  margin: 20rpx 0 10rpx;
  font-family: KaiTi, serif; /* 使用楷体 */
}

/* 正文 */
.card-detail-description {
  font-size: 30rpx;
  line-height: 1.7;
  color: var(--text-light);
  text-align: justify;
  font-family: KaiTi, serif; /* 使用楷体 */
  text-indent: 2em; /* 首行缩进 */
}
</style>