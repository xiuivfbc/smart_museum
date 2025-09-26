<template>
	<view class="transparent-profile">
		<!-- 头像区域 -->
		<view class="avatar-container">
			<image class="avatar" src="/static/head.jpg" mode="aspectFill" />
		</view>

		<!-- 信息列表 -->
		<view class="info-list">
			<view class="info-item">
				<text class="label">昵称</text>
				<text class="value">{{ userform.username || '未设置' }}</text>
			</view>
			<view class="divider"></view>

			<view class="info-item">
				<text class="label">手机号</text>
				<text class="value">{{ userform.phone || '未绑定' }}</text>
			</view>
			<view class="divider"></view>

			<view class="info-item">
				<text class="label">邮箱</text>
				<text class="value">{{ userform.email || '未绑定' }}</text>
			</view>
			<view class="divider"></view>
			<view class="info-item">
				<text class="label">ID</text>
				<text class="value">{{ userform.id || '未设置' }}</text>
			</view>
			<view class="divider"></view>
			<view class="info-item">
				<text class="label">角色</text>
				<text class="value">{{ userform.role || '未设置' }}</text>
			</view>
			<view class="divider"></view>
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
				}
			}
		},
		onShow() {
			// 读取本地用户信息
			const userInfo = uni.getStorageSync("expert_user")
			if (userInfo) {
				this.userform.phone = userInfo.phone || ""
				this.userform.email = userInfo.email || ""
				this.userform.username = userInfo.username || ""
				this.userform.id = userInfo.id || ""
				
				if (userInfo.role === "user") this.userform.role = "公众"
				if (userInfo.role === "expert") this.userform.role = "专家" 
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

		.info-list {
			width: 100%;

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
		}
	}
</style>
