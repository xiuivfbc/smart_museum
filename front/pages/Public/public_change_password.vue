<template>
	<view class="register-container">
		<text class="welcome">修改密码</text>

		<view class="input-group">
			<text class="input-label">原密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="请输入原密码" placeholder-class="underline-placeholder"
			 password />
		</view>


		<view class="input-group">
			<text class="input-label">新密码</text>
			<input v-model="new_password" class="underline-input" placeholder="请输入新密码" placeholder-class="underline-placeholder"
			 password />
		</view>

		<view class="input-group">
			<text class="input-label">确认新密码</text>
			<input v-model="confirmPassword" class="underline-input" placeholder="请再次输入新密码" placeholder-class="underline-placeholder"
			 password />
		</view>

		<button @click="upload" class="register-btn" :class="{'btn-loading': loading}" :disabled="loading">
			<text v-if="!loading">提交</text>
			<text v-else>提交中...</text>
		</button>
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
					password: "",
					role: "user"
				},
				new_password: "",
				confirmPassword: "",
				loading: false
			}
		},
		methods: {
			// async 修饰异步函数
			async upload() {
				if (!this.userform.password || this.userform.password.length < 6) {
					uni.showToast({
						title: "原密码格式错误",
						icon: "none"
					})
					return
				}

				if (!this.new_password || this.new_password.length < 6) {
					uni.showToast({
						title: "新密码至少六位",
						icon: "none"
					})
					return
				}

				// 新增：验证两次密码是否一致
				if (this.new_password !== this.confirmPassword) {
					uni.showToast({
						title: "两次输入的新密码不一致",
						icon: "none"
					})
					return
				}

				if (this.loading) {
					uni.showToast({
						title: "正在修改中，请稍候",
						icon: "none"
					})
					return
				}

				// 开始处理请求，禁止按钮功能
				this.loading = true


				const userInfo = uni.getStorageSync("public_user")

				const requestData = {
					password: this.userform.password,
					id: userInfo.id,
					new_password: this.new_password
				};


				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/auth/upload',
					method: 'POST',
					data: requestData,
					header: {
						'Content-Type': 'application/json'
					},
					timeout: 3000 // 最多等3s
				});

				console.log('后端返回的数据:', res.data);

				if (res.data.error) {
					if (res.data.error === '密码错误') {
						uni.showToast({
							title: "原密码错误",
							icon: "none"
						});
					} else {
						uni.showToast({
							title: "未知错误",
							icon: "none"
						});
					}

					this.loading = false;
					return;
				} else if (res.statusCode === 200) {

					const userInfo = uni.getStorageSync("public_user")

					userInfo.password = this.new_password
					uni.setStorageSync("public_user", userInfo);

					this.loading = false;

					uni.showToast({
						title: "修改成功",
						icon: "none"
					})
					this.userform = {
						id: "",
						username: "",
						phone: "",
						email: "",
						password: "",
						role: "user"
					};
					
					setTimeout(() => {
					uni.redirectTo({
						url: "/pages/Public/public_me"
					});
					}, 500);
					
				} else {
					this.loading = false;
					uni.showToast({
						title: `网络请求失败: ${res.statusCode}`,
						icon: "none"
					});
					return
				}
			}
		}
	}
</script>

<style>
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

	.register-btn {
		background-color: #4a89b3;
		color: #fff;
		border: none;
	}

	.register-btn:active {
		background-color: rgba(52, 152, 219, 0.1);
	}
</style>
