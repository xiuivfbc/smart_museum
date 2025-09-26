<template>
	<view class="transparent-profile">

		<view class="avatar-container">
			<image class="avatar" src="/static/head.jpg" mode="aspectFill" />
		</view>

		<text class="name">{{ userform.username || '未设置' }}</text>

		<!-- 获取二维码 -->
		<button @click="fetchQrcode" class="submit-btn" :disabled="loading">
			<text v-if="!loading">展示二维码</text>
			<text v-else>获取中...</text>
		</button>

		<view class="qrcode-container" v-if="showQrcode">
			<image class="qrcode-image" :src="qrcodeUrl" mode="widthFix" />
		</view>


		<view class="info-item" v-if="showQrcode">
			<text class="label">时间</text>
			<text class="value">{{ time.substring(0, 10) }}</text>
		</view>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				userform: {
					id: "",
					username: "",
					phone: "",
					email: "",
					role: ""
				},
				showQrcode: false, // 控制二维码显示
				qrcodeUrl: "", // 存储二维码图片URL
				loading: false, // 按钮加载状态
				time: ""
			}
		},
		onShow() {
			// 读取本地用户信息
			const userInfo = uni.getStorageSync("public_user")
			if (userInfo) {
				this.userform = userInfo
			}
		},
		methods: {
			async fetchQrcode() {
				if (this.loading) return

				this.loading = true
				this.showQrcode = false // 先隐藏
				const userInfo = uni.getStorageSync("public_user")
				console.log(userInfo.id)
				console.log(uni.getStorageSync("public_token"))

				const requestData = {
					id: userInfo.id
				};

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/ticket/list', // 替换为实际接口
					method: 'POST',
					data: requestData,
					header: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${uni.getStorageSync("public_token").trim()}`
					},
					timeout: 3000
				})


				if (err) {
					console.error('请求失败:', err);
					uni.showToast({
						title: '请求失败，请检查网络',
						icon: 'none'
					});
					this.loading = false;
					return;
				}

				console.log('后端返回的数据:', res.data);

				if (res.data.error) {
					uni.showToast({
						title: '未查询到票务',
						icon: 'none'
					})

					this.loading = false
					setTimeout(() => {
						uni.redirectTo({
							url: "/pages/Public/public_me"
						});
					}, 900);
					return
				}

				if (res.data && res.statusCode === 200) {
					const ticket = res.data.ticket

					this.qrcodeUrl = `http://127.0.0.1:8081/qrcodes/${ticket.path}`
					this.showQrcode = true
					this.time = ticket.time

				} else {
					uni.showToast({
						title: res.data.error || '获取二维码失败',
						icon: 'none'
					})
				}

				this.loading = false
				return
			}
		}
	}
</script>

<style lang="scss">
	.transparent-profile {
		width: 100%;
		padding: 40rpx 0;

		.avatar-container {
			display: flex;
			justify-content: center;
			margin-bottom: 60rpx;

			.avatar {
				width: 160rpx;
				height: 160rpx;
				border-radius: 50%;
				background-color: rgba(255, 255, 255, 0.2);
				border: 2rpx solid rgba(255, 255, 255, 0.3);
			}
		}

		/* 新增二维码样式 */
		.qrcode-container {
			margin: 40rpx auto;
			padding: 20rpx;
			background-color: #fff;
			border-radius: 16rpx;
			width: 60%;
			display: flex;
			justify-content: center;
			box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);

			.qrcode-image {
				width: 100%;
				height: auto;
			}
		}

		/* 按钮样式（与之前一致） */
		.submit-btn {
			width: 80%;
			margin: 80rpx auto 0;
			padding: 10rpx 0;
			font-size: 28rpx;
			border-radius: 100rpx;
			background-color: #4a89b3;
			color: #fff;
			border: none;
		}

		.submit-btn[disabled] {
			background-color: #ccc;
		}

		.name {
			font-size: 36rpx;
			font-weight: 600;
			color: #2c3e50;
			text-align: center;
			display: block;
			margin-bottom: 20rpx;
		}

		.divider {
			height: 1px;
			background-color: #eee;
			margin: 0 40rpx;
		}
	}

	.info-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 30rpx;
		width: 100%;
		box-sizing: border-box;

		.label {
			font-size: 32rpx;
			color: inherit;
		}

		.value {
			font-size: 32rpx;
			color: inherit;
		}
	}

	.divider {
		height: 1rpx;
		background: rgba(0, 0, 0, 0.1);
		margin: 0 30rpx;
	}
</style>