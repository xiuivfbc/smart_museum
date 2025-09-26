<template>
	<view class="tab-bar">
		
		<!-- item : 当前遍历到的 tab 对象; index : 当前索引（0,1,2） -->
		<view class="tab-item" v-for="(item, index) in tabs" :key="index" @click="goPage(item.path)">
			<image class="tab-icon" :src="isActive(item) ? item.activeIcon : item.icon" mode="aspectFit" />
			<text class="tab-text">{{ item.name }}</text>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				tabs: [{
						name: '首页',
						path: '/pages/Expert/expert_home',
						alias: ['/pages/Expert/expert_home'],
						icon: '/static/icons/home.png',
						activeIcon: '/static/icons/home-fill.png'
					},
					{
						name: 'AI',
						path: '/pages/Expert/expert_ai',
						alias: ['/pages/Expert/expert_ai'],
						icon: '/static/icons/ai.png',
						activeIcon: '/static/icons/ai-fill.png'
					},
					{
						name: '我的',
						path: '/pages/Expert/expert_me',
						alias: [
							'/pages/Expert/expert_me',
							'/pages/Expert/expert_login',
							'/pages/Expert/expert_register'
						],
						icon: '/static/icons/me.png',
						activeIcon: '/static/icons/me-fill.png'
					}
				],
				currentPath: ''
			}
		},
		
		// 当组件实例挂载到 DOM 后 会触发 mounted(),比如 public_tabbar 挂载后
		mounted() {
			const pages = getCurrentPages()
			const currentPage = pages[pages.length - 1]
			this.currentPath = '/' + currentPage.route
		},
		methods: {
			isActive(item) {
				return item.alias.includes(this.currentPath)
			},
			goPage(url) {
				if (this.currentPath === url) return;

				// 一级 tab 页面（首页 / AI / 我的）用 redirectTo / navigateTo 都可以，因为 tabbar 独立存在
				this.currentPath = url
				uni.redirectTo({
					url
				})
			}
		}
	}
</script>

<style scoped>
	.tab-bar {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		height: 60px;
		display: flex;
		background-color: #fff;
	}

	.tab-item {
		flex: 1;
		text-align: center;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	.tab-icon {
		width: 24px;
		height: 24px;
	}

	.tab-text {
		font-size: 12px;
		color: #666;
	}
</style>
