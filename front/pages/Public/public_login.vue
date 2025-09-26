<template>
	<view class="login-container">

		<text class="welcome">登录</text>

		<view class="tab-group">
			<text :class="['tab-item', activeTab === 'phone' ? 'tab-active' : '']" @click="activeTab = 'phone'">
				手机号登录
			</text>
			<text :class="['tab-item', activeTab === 'email' ? 'tab-active' : '']" @click="activeTab = 'email'">
				邮箱登录
			</text>
		</view>

		<view v-if="activeTab === 'phone'" class="input-group">
			<text class="input-label">手机号</text>
			<input v-model="userform.phone" class="underline-input" placeholder="请输入手机号"
				placeholder-class="underline-placeholder" type="number" />
		</view>

		<view v-if="activeTab === 'email'" class="input-group">
			<text class="input-label">邮箱</text>
			<input v-model="userform.email" class="underline-input" placeholder="请输入邮箱"
				placeholder-class="underline-placeholder" type="email" />
		</view>

		<view class="input-group">
			<text class="input-label">密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="请输入密码"
				placeholder-class="underline-placeholder" password />
		</view>

		<button @click="login" class="login-btn">
			登录
		</button>

		<button @click="goRegister" class="register-btn">
			去注册
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
				activeTab: 'phone', // 默认手机登录
				userform: {
					id: "",
					username: "",
					phone: "",
					email: "",
					password: "",
					role: "user"
				},
				loading: false // 防止重复提交登录请求
			}
		},
		onShow() {
			// 自动填充本地缓存的登录信息（如果有）
			const userInfo = uni.getStorageSync("public_user")
			if (userInfo) {
				this.userform = userInfo

				// 根据填充的数据自动选择标签
				if (userInfo.email) {
					this.activeTab = 'email'
				} else if (userInfo.phone) {
					this.activeTab = 'phone'
				}
			}
		},

		methods: {
			async login() {
				// 前端简单验证
				if (this.activeTab === 'phone') {
					if (!this.userform.phone) {
						uni.showToast({
							title: "请输入手机号",
							icon: "none"
						})
						return
					}
					if (this.userform.phone.length !== 11) {
						uni.showToast({
							title: "手机号格式错误",
							icon: "none"
						})
						return
					}
				} else {
					if (!this.userform.email) {
						uni.showToast({
							title: "请输入邮箱",
							icon: "none"
						})
						return
					}
					if (!this.userform.email.includes('@')) {
						uni.showToast({
							title: "邮箱格式错误",
							icon: "none"
						})
						return
					}
				}

				if (!this.userform.password) {
					uni.showToast({
						title: "请输入密码",
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

				if (this.loading) {
					uni.showToast({
						title: "正在登录中，请稍候",
						icon: "none"
					})
					return
				}

				// 开始处理请求，禁止按钮功能
				this.loading = true

				const requestData = {
					password: this.userform.password,
					role: this.userform.role,
					identifier: ""
				};

				if (this.activeTab === 'phone') requestData.identifier = this.userform.phone;
				else requestData.identifier = this.userform.email;

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/auth/login',
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

					if (data.message === "登录成功") {
						data.userform.password = this.userform.password
						this.userform = data.userform
						
						console.log(data.token)
						
						// 用户信息写入本地缓存
						uni.setStorageSync("public_user", this.userform);
						uni.setStorageSync("public_token", data.token);

						this.loading = false;

						uni.showToast({
							title: "登录成功",
							icon: "none"
						})

						// 延迟0.5s后跳转，采用替换页面
						setTimeout(() => {
							uni.redirectTo({
								url: "/pages/Public/public_me"
							});
						}, 500);

					} else {
						this.loading = false;
						uni.showToast({
							title: data.error || "登录失败，未知错误",
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
			},

			goRegister() {
				uni.navigateTo({
					url: "/pages/Public/public_register"
				})
			}

		}

	}
</script>

<style>
	.login-container {
		padding: 60rpx 40rpx;
		background-color: #f8f9fa;
		min-height: 100vh;
		background:
			linear-gradient(rgba(255, 255, 255, 0.5), rgba(255, 255, 255, 0.5)),
			url('/static/background/3.jpg');

	}

	.welcome {
		font-size: 36rpx;
		font-weight: 600;
		color: #2c3e50;
		text-align: center;
		display: block;
		margin-bottom: 40rpx;
	}

	/* 切换标签样式 */
	.tab-group {
		display: flex;
		justify-content: center;
		margin-bottom: 40rpx;
		border-bottom: 1rpx solid #dcdfe6;
	}

	.tab-item {
		padding: 20rpx 40rpx;
		font-size: 28rpx;
		color: #666;
		cursor: pointer;
	}

	.tab-active {
		color: #3498db;
		border-bottom: 2rpx solid #3498db;
	}

	/* 复用注册页面的输入框样式 */
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

	.login-btn,
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

	.login-btn {
		background-color: #4a89b3;
		color: #fff;
		border: none;
	}

	.register-btn {
		margin-top: 30rpx;
		/* 注册按钮与登录按钮间距 */
		background-color: transparent;
		color: #4a89b3;
		border: 1rpx solid #4a89b3;
	}

	/* 按钮悬停效果 */
	.login-btn:active {
		background-color: rgba(52, 152, 219, 0.1);
	}

	.register-btn:active {
		background-color: rgba(52, 152, 219, 0.1);
	}
</style>