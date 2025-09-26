<template>
	<view class="profile-container">
		<text class="welcome">退票</text>

		<!-- 仅保留密码输入 -->
		<view class="input-group">
			<text class="input-label">密码</text>
			<input v-model="password" class="underline-input" placeholder="请输入密码"
				placeholder-class="underline-placeholder" password type="text" @input="validatePassword" />

		</view>


		<button @click="handleSubmit" class="submit-btn" :disabled="loading || passwordError">
			<text v-if="!loading">提交</text>
			<text v-else>提交中...</text>
		</button>
	</view>
</template>

<script>
	export default {
		data() {
			return {
				password: '',
				passwordError: false,
				loading: false
			}
		},
		methods: {
			validatePassword() {
				this.passwordError = this.password.length < 6
			},
			async handleSubmit() {
				if (!this.password || this.passwordError) {
					uni.showToast({
						title: "请输入有效密码",
						icon: "none"
					})
					return
				}

				const userInfo = uni.getStorageSync("public_user")

				if (this.password !== userInfo.password) {
					uni.showToast({
						title: "密码错误",
						icon: "none"
					})
					return
				}

				this.loading = true

				console.log(userInfo.id)
				const requestData = {
					id: userInfo.id
				};

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/ticket/delete',
					method: 'DELETE',
					data: requestData,
					header: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${uni.getStorageSync("public_token").trim()}`
					},
					timeout: 3000
				});

				if (res.statusCode === 200) {
					uni.showToast({
						title: "退票成功",
						icon: "none"
					});
					setTimeout(() => {
						uni.redirectTo({
							url: "/pages/Public/public_me"
						});
					}, 500);
				} else {
					uni.showToast({
						title: res.data?.error || "退票失败",
						icon: "none"
					});
					this.loading = false
				}
			}
		}
	}
</script>

<style>
	.profile-container {
		padding: 60rpx 40rpx;
		background-color: #f8f9fa;
		min-height: 100vh;
	}

	.welcome {
		font-size: 36rpx;
		font-weight: 600;
		color: #2c3e50;
		text-align: center;
		display: block;
		margin-bottom: 60rpx;
	}

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
	}

	.underline-input:focus {
		border-bottom-color: #3498db;
		outline: none;
	}

	.underline-placeholder {
		color: #a0a0a0;
		font-size: 28rpx;
	}

	.error-tip {
		color: #e74c3c;
		font-size: 24rpx;
		margin-top: 8rpx;
		display: block;
	}

	.note {
		font-size: 24rpx;
		color: #7f8c8d;
		text-align: center;
		display: block;
		margin: 20rpx 0;
	}

	.submit-btn {
		width: 80%;
		margin: 60rpx auto 0;
		padding: 12rpx 0;
		font-size: 28rpx;
		border-radius: 100rpx;
		background-color: #4a89b3;
		color: #fff;
		border: none;
		transition: all 0.3s;
	}

	.submit-btn[disabled] {
		background-color: #ccc;
	}

	.submit-btn:active {
		background-color: #3a76a3;
	}
</style>