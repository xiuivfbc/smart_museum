<template>
	<view class="profile-container">
		<text class="welcome">改签</text>

		<!-- 日期选择 -->
		<view class="input-group">
			<text class="input-label">新日期</text>
			<uni-datetime-picker type="date" v-model="date" :start="minDate" :end="maxDate" @change="bindDateChange"
				class="underline-input" />
		</view>

		<!-- 密码验证 -->
		<view class="input-group">
			<text class="input-label">密码</text>
			<input v-model="password" class="underline-input" placeholder="请输入密码"
				placeholder-class="underline-placeholder" password type="text" @input="validatePassword" />

		</view>


		<button @click="handleReschedule" class="submit-btn" :class="{'btn-loading': loading}"
			:disabled="loading || passwordError">
			<text v-if="!loading">确认改签</text>
			<text v-else>处理中...</text>
		</button>
	</view>
</template>

<script>
	import uniDatetimePicker from '@/uni_modules/uni-datetime-picker/components/uni-datetime-picker/uni-datetime-picker.vue'
	export default {
		components: {
			uniDatetimePicker
		},
		data() {
			return {
				date: '',
				password: '',
				passwordError: '',
				minDate: this.getToday(),
				maxDate: this.getDateAfterDays(14),
				loading: false,
				rfc3339Time: ''
			}
		},
		methods: {
			validatePassword() {
				if (this.password.length < 6) {
					this.passwordError = true
				} else {
					this.passwordError = false
				}
			},
			getToday() {
				const today = new Date()
				return today.toISOString().split('T')[0]
			},
			getDateAfterDays(days) {
				const date = new Date()
				date.setDate(date.getDate() + days)
				return date.toISOString().split('T')[0]
			},
			bindDateChange(e) {
				this.date = e
				this.rfc3339Time = `${e}T00:00:00Z`
			},
			async handleReschedule() {
				if (!this.password || this.passwordError) {
					uni.showToast({
						title: this.passwordError || "请输入有效密码",
						icon: "none"
					})
					return
				}
				if (!this.date) {
					uni.showToast({
						title: "请选择新日期",
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

				const requestData = {
					time: this.rfc3339Time,
					id: userInfo.id
				};

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/ticket/update',
					method: 'PUT',
					data: requestData,
					header: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${uni.getStorageSync("public_token").trim()}`
					},
					timeout: 3000
				});

				if (res.statusCode === 200) {
					uni.showToast({
						title: "改签成功",
						icon: "none"
					});
					setTimeout(() => {
						uni.redirectTo({
							url: "/pages/Public/public_me"
						});
					}, 500);
				} else {
					uni.showToast({
						title: res.data?.error || "改签失败",
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

	.uni-datetime-picker .uni-date-editor--datetime {
		border: none !important;
		border-bottom: 1px solid #dcdfe6 !important;
		padding: 20rpx 0 !important;
	}

	.uni-datetime-picker .uni-input-input {
		font-size: 30rpx !important;
		color: #2c3e50 !important;
	}
</style>