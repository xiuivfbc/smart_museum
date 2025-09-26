<template>
	<view class="register-container">
		<text class="welcome">注册</text>

		<view class="input-group">
			<text class="input-label">手机号</text>
			<input v-model="userform.phone" class="underline-input" placeholder="手机号/邮箱 请至少填写一个" placeholder-class="underline-placeholder"
			 type="number" />
		</view>

		<view class="input-group">
			<text class="input-label">邮箱</text>
			<input v-model="userform.email" class="underline-input" placeholder="手机号/邮箱 请至少填写一个" placeholder-class="underline-placeholder"
			 type="email" />
		</view>
		
		<view class="input-group">
			<text class="input-label">认证码</text>
			<input v-model="userform.active" class="underline-input" placeholder="请输入认证码"
				placeholder-class="underline-placeholder" password />
		</view>

		<view class="input-group">
			<text class="input-label">密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="密码至少六位" placeholder-class="underline-placeholder"
			 password />
		</view>

		<!-- 新增：确认密码输入框 -->
		<view class="input-group">
			<text class="input-label">确认密码</text>
			<input v-model="confirmPassword" class="underline-input" placeholder="请再次输入密码" placeholder-class="underline-placeholder"
			 password />
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
		components: {
			public_tabbar
		},
		data() {
			return {
				userform: {
					id: "",
					username: "",
					phone: "",
					email: "",
					password: "",
					role: "admin",
					active:"",
				},
				confirmPassword: "", // 新增：确认密码字段
				loading: false // 防止重复提交注册请求
			}
		},
		methods: {
			// async 修饰异步函数
			async register() {
				if (!this.userform.phone && !this.userform.email) {
					uni.showToast({
						title: "手机号/邮箱 请至少填写一个",
						icon: "none"
					})
					return
				}	

				if (this.userform.phone && this.userform.phone.length !== 11) {
					uni.showToast({
						title: "手机号格式错误",
						icon: "none"
					})
					return
				}

				if (this.userform.email && !this.userform.email.includes('@')) {
					uni.showToast({
						title: "邮箱格式错误",
						icon: "none"
					})
					return
				}
				
				if (!this.userform.active) {
					uni.showToast({
						title: "认证码为空",
						icon: "none"
					})
					return
				}
				
				if (!this.userform.password || this.userform.password.length < 6) {
					uni.showToast({
						title: "密码至少六位",
						icon: "none"
					})
					return
				}

				// 新增：验证两次密码是否一致
				if (this.userform.password !== this.confirmPassword) {
					uni.showToast({
						title: "两次输入的密码不一致",
						icon: "none"
					})
					return
				}

				if (this.loading) {
					uni.showToast({
						title: "正在注册中，请稍候",
						icon: "none"
					})
					return
				}

				// 开始处理请求，禁止按钮功能
				this.loading = true


				const requestData = {
					password: this.userform.password,
					role: this.userform.role,
					username: "Expert",
					identifier: this.userform.active,
				};

				if (this.userform.phone) requestData.phone = this.userform.phone;
				if (this.userform.email) requestData.email = this.userform.email;


				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/auth/register',
					method: 'POST',
					data: requestData,
					header: {
						'Content-Type': 'application/json'
					},
					timeout: 3000 // 最多等3s
				});

				if (err) {
					console.error('请求错误:', err);

					if (err.errMsg && err.errMsg.includes('timeout')) {
						uni.showToast({
							title: "请求超时，请检查网络",
							icon: "none"
						});
					} else {
						uni.showToast({
							title: "网络异常，请检查连接",
							icon: "none"
						});
					}

					this.loading = false;
					return;
				}

				console.log('后端返回的数据:', res.data);

				// 成功响应
				if (res.statusCode === 200) {
					const data = res.data;

					if (data.message === "注册成功") {
						
						// 缓存用户数据，login 时就只需要输入密码
						this.userform.password = ""
						uni.setStorageSync("expert_user", this.userform)
						uni.showToast({
							title: "注册成功，正在跳转至登录页",
							icon: "none"
						});
						this.loading = false;
						// 延迟跳转，让用户看到成功的提示
						setTimeout(() => {
							uni.redirectTo({
								url: "/pages/Expert/expert_login"
							});
						}, 700);

					} else {
						this.loading = false;
						uni.showToast({
							title: data.message || "注册失败，请重试",
							icon: "none"
						});
					}
				} else {
					this.loading = false;
					uni.showToast({
						title: `网络请求失败: ${res.statusCode}`,
						icon: "none"
					});
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
		