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
						<text class="grid-text logout-text">注销登录</text>
					</view>
				</view>
			</uni-card>

			<!-- 票务管理卡片 -->
			<uni-card title="票务管理" is-full>
				<view class="grid-container">
					<view class="grid-item" @click="goBuyTickets">
						<image class="grid-icon" src="/static/icons/buy.png" mode="aspectFit"></image>
						<text class="grid-text">购票</text>
					</view>
					<view class="grid-item" @click="goTicketList">
						<image class="grid-icon" src="/static/icons/dd.png" mode="aspectFit"></image>
						<text class="grid-text">订单</text>
					</view>
					<view class="grid-item" @click="goChange">
						<image class="grid-icon" src="/static/icons/gq.png" mode="aspectFit"></image>
						<text class="grid-text">改签</text>
					</view>
					<view class="grid-item" @click="goRefund">
						<image class="grid-icon" src="/static/icons/tp.png" mode="aspectFit"></image>
						<text class="grid-text">退票</text>
					</view>
				</view>
			</uni-card>
		

		<public_tabbar />
		</view>
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
			const token = uni.getStorageSync("public_token")
			const userInfo = uni.getStorageSync("public_user")
			this.userform = userInfo
			if (token) {
				this.isLogin = true
				this.userform = userInfo
			} else {
				uni.redirectTo({
					url: "/pages/Public/public_login"
				})
			}
		},
		methods: {
			logout() {
				uni.removeStorageSync("public_token")
				uni.removeStorageSync("public_user")
				this.isLogin = false

				uni.showToast({
					title: "已退出登录",
					icon: "none"
				})

				uni.redirectTo({
					url: "/pages/Public/public_login"
				})
			},
			goProfileDetail() {
				uni.navigateTo({
					url: "/pages/Public/public_profile"
				})
			},
			goEditProfile() {
				uni.navigateTo({
					url: "/pages/Public/public_edit_profile"
				})
			},
			goChangePassword() {
				uni.navigateTo({
					url: "/pages/Public/public_change_password"
				})
			},
			goBuyTickets() {
				uni.navigateTo({
					url: "/pages/Ticket/ticket_buy"
				})
			},
			goTicketList() {
				uni.navigateTo({
					url: "/pages/Ticket/ticket_check"
				})
			},
			goChange() {
				uni.navigateTo({
					url: "/pages/Ticket/ticket_change"
				})
			},
			goRefund() {
				uni.navigateTo({
					url: "/pages/Ticket/ticket_refund"
				})
			}
		}
	}
</script>

<style lang="scss">
	// 全局变量定义
	:root {
		--glass-bg: rgba(255, 255, 255, 0.25);
		--glass-border: rgba(255, 255, 255, 0.3);
		--glass-shadow: 0 8rpx 32rpx rgba(31, 38, 135, 0.15);
	}

	.me-container {
		padding: 40rpx 30rpx;
		// 设置背景图
		background-image: url('/static/background/3.jpg');
		background-size: cover;
		background-position: center;
		background-attachment: fixed;
		min-height: 100vh;
	}

	.welcome {
		font-size: 36rpx;
		font-weight: 600;
		color: #2c3e50;
		text-align: center;
		display: block;
		margin-bottom: 40rpx;
		// 文字阴影增加可读性
		text-shadow: 0 2rpx 4rpx rgba(255, 255, 255, 0.8);
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
		// 毛玻璃效果
		background: var(--glass-bg);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border: 1rpx solid var(--glass-border);
		box-shadow: var(--glass-shadow);
		transition: all 0.3s ease;
	}

	.grid-item:active {
		background-color: rgba(255, 255, 255, 0.35);
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

	/* 卡片样式 - 毛玻璃效果 */
	.uni-card {
		margin-bottom: 30rpx;
		// 移除默认背景，应用毛玻璃
		background: var(--glass-bg) !important;
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border: 1rpx solid var(--glass-border) !important;
		box-shadow: var(--glass-shadow);
		border-radius: 16rpx !important;
	}

	// 卡片标题样式调整
	.uni-card__header {
		color: #2c3e50 !important;
		font-weight: 600 !important;
	}
</style>
