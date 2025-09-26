<template>
  <view class="container">
    <!-- 首页轮播 -->
    <swiper class="swiper" indicator-dots="true" autoplay="true" interval="3000" duration="500" circular="true">
      <swiper-item v-for="(img, index) in images" :key="index">
        <image :src="img" mode="aspectFill" class="swiper-image"></image>
        <view class="swiper-title">{{ getTitleByIndex(index) }}</view>
      </swiper-item>
    </swiper>
    <!-- 镇馆之宝 -->
    <view class="title-section">镇馆之宝</view>
    <!-- 卡片模块 -->
    <view class="card" v-for="(item, index) in cards" :key="index">
      <navigator :url="`/pages/cardDetail/cardDetail?index=${index}`" class="card-link">
        <view class="card-content">
          <view class="card-cover">
            <image :src="item.cover" mode="aspectFill" class="cover-image"></image>
          </view>
          <view class="card-title">{{ item.title }}</view>
        </view>
      </navigator>
    </view>
    <!-- 自定义 tabBar -->
    <public_tabbar />
  </view>
</template>

<script>
import public_tabbar from '../Common/public_tabbar.vue';

export default {
  components: {
    public_tabbar
  },
  data() {
    return {
      images: [
        '/static/rounds/1.jpg',
        '/static/rounds/2.jpg',
        '/static/rounds/3.jpg'
      ],
      cards: [
        {
          cover: '/static/cards/4.jpg',
          title: '金缕玉衣',
          description: '金缕玉衣是汉代皇帝和高级贵族死时穿用的殓服，外观和人体形状相同，由玉片拼成，并用金丝加以编缀。这种玉衣不仅只是简单磨成玉片而已，上面还雕有花纹，反映了玉师杰出的技艺和达官奢侈的生活。金缕玉衣的起源可以追溯到东周时的“缀玉面幕”“缀玉衣服”，到三国时曹丕下诏禁用玉衣，共流行了四百年。汉代人认为玉是“山岳精英”，将金玉置于人的九窍，人的精气不会外泄，就能使尸骨不腐，可求来世再生，所以用于丧葬的玉器在汉玉中占有重要的地位。金缕玉衣是汉代规格最高的丧葬殓服，大致出现在西汉文景时期。'
        },
        {
          cover: '/static/cards/6.jpg',
          title: '乳钉纹青铜爵',
          description: '乳钉纹青铜爵是一件来自二里头文化的文物，以其俊巧的艺术美感和丰富的历史内涵著称。它不仅是夏代青铜冶铸技术的实物见证，更是贵族们身份和地位的象征。乳钉纹青铜爵的铸造工艺堪称夏代青铜铸造技术的巅峰之作，其采用的是复合范铸法，这是一种较为复杂的铸造工艺，展现了惊人的工艺水平。乳钉纹青铜爵的艺术成就对后世青铜器的发展产生了深远的影响，成为中国古代青铜器艺术的经典之作。'
        },
        {
          cover: '/static/cards/7.jpg', // 替换为你的图片路径
          title: '青花寿山福海纹瓷香炉',
          description: '造型为青铜鼎样式，盘口，方唇，束颈，圆腹，小平底。纹样以苏麻离青料绘蓝地白花寿山福海纹，青花发色浓重。'
        }
      ]
    }
  },
  methods: {
    getTitleByIndex(index) {
      const titles = [
        '探秘千年：金缕玉衣修复工坊',
        '青铜技艺：乳钉纹爵铸造体验',
        '玉润华夏：汉玉文化深度研学'
      ];
      return titles[index];
    }
  }
}
</script>

<style scoped>
/* =========== 极简毛玻璃首页 =========== */
:root {
  --primary: #8c6f58;
  --glass-bg: rgba(255, 255, 255, 0.58);
  --glass-border: rgba(255, 255, 255, 0.75);
  --text-dark: #2e2e2e;
  --text-light: #666;
  --shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.06);
  --radius: 28rpx;
}

/* 页面容器：背景+毛玻璃 */
.container {
  min-height: 100vh;
  padding: 30rpx 30rpx 140rpx;
  background: linear-gradient(rgba(255, 255, 255, 0.3), rgba(255, 255, 255, 0.3)),
              url('/static/background/3.jpg') center/cover no-repeat;
  position: relative;
}

/* 毛玻璃通用类 */
.frost-glass {
  background: var(--glass-bg);
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1rpx solid var(--glass-border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  transition: transform .25s, box-shadow .25s;
}
.frost-glass:active {
  transform: translateY(2rpx);
  box-shadow: 0 6rpx 24rpx rgba(0, 0, 0, 0.04);
}

/* ===== 轮播 ===== */
.swiper {
  width: 100%;
  height: 360rpx;
  border-radius: var(--radius);
  overflow: hidden;
  margin-bottom: 40rpx;
}
.swiper-image {
  width: 100%;
  height: 100%;
}
.swiper-title {
  position: absolute;
  left: 30rpx;
  bottom: 30rpx;
  font-size: 28rpx;
  color: #fff;
  background: rgba(0, 0, 0, 0.45);
  padding: 8rpx 18rpx;
  border-radius: 12rpx;
}

/* ===== 镇馆之宝 ===== */
.title-section {
  width: 100%;
  height: 90rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36rpx;
  font-weight: 600;
  color: var(--text-dark);
  letter-spacing: 4rpx;
  margin-bottom: 40rpx;
}

/* ===== 卡片列表 ===== */
.card {
  width: 92%;
  margin: 0 auto 40rpx;
}
.card-link {
  display: block;
}
.card-content {
  background: var(--glass-bg);
  border-radius: var(--radius);
  overflow: hidden;
  padding: 20rpx;
}
.card-cover {
  width: 100%;
  height: 280rpx;
  border-radius: var(--radius);
  overflow: hidden;
}
.cover-image {
  width: 100%;
  height: 100%;
}
.card-title {
  margin-top: 20rpx;
  font-size: 34rpx;
  font-weight: 600;
  color: var(--text-dark);
  text-align: center;
}
</style>