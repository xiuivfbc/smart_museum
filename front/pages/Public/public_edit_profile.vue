<template>
	<view class="profile-container">
		<text class="welcome">更新信息</text>

		<!-- 切换三个选项卡 -->
		<view class="tab-group">
			<text :class="['tab-item', activeTab === 'phone' ? 'tab-active' : '']" @click="activeTab = 'phone'">
				绑定手机号
			</text>
			<text :class="['tab-item', activeTab === 'email' ? 'tab-active' : '']" @click="activeTab = 'email'">
				绑定邮箱
			</text>
			<text :class="['tab-item', activeTab === 'nickname' ? 'tab-active' : '']" @click="activeTab = 'nickname'">
				修改昵称
			</text>
		</view>

		<!-- 手机号绑定 -->
		<view v-if="activeTab === 'phone'" class="input-group">
			<text class="input-label">手机号</text>
			<input v-model="userform.phone" class="underline-input" placeholder="请输入11位手机号" type="number" maxlength="11" />

			<text class="input-label">当前密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="请输入密码" password />
		</view>

		<!-- 邮箱绑定 -->
		<view v-if="activeTab === 'email'" class="input-group">
			<text class="input-label">邮箱</text>
			<input v-model="userform.email" class="underline-input" placeholder="请输入有效邮箱" type="email" />

			<text class="input-label">当前密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="请输入密码" password />
		</view>

		<!-- 昵称修改 -->
		<view v-if="activeTab === 'nickname'" class="input-group">
			<text class="input-label">新昵称</text>
			<input v-model="userform.username" class="underline-input" placeholder="2-10个字符" maxlength="10" />

			<text class="input-label">当前密码</text>
			<input v-model="userform.password" class="underline-input" placeholder="请输入密码" password />
		</view>

		<!-- 提交按钮 -->
		<button @click="handleSubmit" class="submit-btn" :class="{'btn-loading': loading}" :disabled="loading">
			<text v-if="!loading">提交</text>
			<text v-else>提交中...</text>
		</button>

	</view>
</template>

<script>
	export default {
		data() {
			return {
				activeTab: 'phone',
				loading: false, // 修正：将 isLoading 改为 loading，以匹配模板中的使用
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
		methods: {
			// 提交修改
			async handleSubmit() {

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
				} else if (this.activeTab === 'email') {
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
				} else {
					if (!this.userform.username) {
						uni.showToast({
							title: "请输入昵称",
							icon: "none"
						})
						return
					}
					if (this.userform.username.length < 2 || this.userform.username.length > 10) {
						uni.showToast({
							title: "昵称长度违规",
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
						title: "正在更新中，请稍候",
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
					username: "",
					phone: "",
					email: ""
				};

				if (this.activeTab === 'phone') {
					requestData.phone = this.userform.phone
				} else if (this.activeTab === 'email') {
					requestData.email = this.userform.email
				} else {
					requestData.username = this.userform.username
				}

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
							title: "密码错误",
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

					if (this.activeTab === 'phone') {
						userInfo.phone = this.userform.phone
					} else if (this.activeTab === 'email') {
						userInfo.email = this.userform.email
					} else {
						userInfo.username = this.userform.username
					}

					uni.setStorageSync("public_user", userInfo);
					
					this.loading = false;

					uni.showToast({
						title: "更新成功",
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
		margin-bottom: 40rpx;
	}

	/* 三选项卡样式 */
	.tab-group {
		display: flex;
		justify-content: space-between;
		margin-bottom: 40rpx;
		border-bottom: 1rpx solid #dcdfe6;
	}

	.tab-item {
		flex: 1;
		text-align: center;
		padding: 20rpx 0;
		font-size: 28rpx;
		color: #666;
	}

	.tab-active {
		color: #3498db;
		border-bottom: 2rpx solid #3498db;
	}

	/* 输入区域 */
	.input-group {
		margin-top: 20rpx;
	}

	.input-label {
		font-size: 28rpx;
		color: #34495e;
		margin: 30rpx 0 16rpx;
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

	/* 提交按钮 */
	.submit-btn {
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
		background-color: #4a89b3;
		color: #fff;
		border: none;
	}

	.submit-btn[disabled] {
		background-color: #ccc;
	}
</style>
