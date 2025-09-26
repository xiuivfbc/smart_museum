<template>
	<view class="container">
		<view v-if="isLogin">
			<!-- 背景模糊层 -->
			<view class="blur-background"></view>

			<!-- 顶部大卡片 - 监测控制 -->
			<view class="big-card frost-glass">
				<view class="dual-column-layout">

					<!-- 右侧内容列 (2/3宽度) -->
					<view class="content-column">

						<text class="column-title"> 展厅区域控制面板</text>
						<button class="column-btn" :class="{'active-state': isRequesting}" @click="GetEnvironment">
							{{ isRequesting ? '实时监测中' : '开启监测' }}
						</button>
					</view>
					<!-- 左侧图标列 (1/3宽度) -->
					<view class="icon-column">

						<image class="column-icon" src="/static/icons/hall.png" mode="aspectFit"></image>
					</view>


				</view>
			</view>

			<!-- 三指标卡片 - 调整了间距 -->
			<view class="metrics-row">
				<view class="metric-card frost-glass centered-content" style="margin-right: 10rpx;">
					<text class="metric-label">温度</text>
					<text class="metric-value">{{ environment.temperature || '--' }} °C</text>
				</view>

				<view class="metric-card frost-glass centered-content" style="margin:0 10rpx;">
					<text class="metric-label">湿度</text>
					<text class="metric-value">{{ environment.humidity || '--' }} %</text>
				</view>

				<view class="metric-card frost-glass centered-content" style="margin-left: 10rpx;">
					<text class="metric-label">气体浓度</text>
					<text class="metric-value">{{ environment.gas || '--' }} ppm</text>
				</view>
			</view>


			<view class="big-card frost-glass">
				<!-- 第一行：60%高度，分左右两列 -->
				<view class="top-row">
					<!-- 左侧：标题+图标 -->
					<view class="left-section">
						<text class="card-main-title">风扇调节系统</text>
						<image class="fan-icon" src="/static/icons/fan.png" mode="aspectFit"></image>
					</view>

					<!-- 右侧：控制区域 -->
					<view class="right-section">
						<!-- 开启阈值 -->
						<text class="control-label">开启阈值</text>

						<uni-number-box @change="yuzhichange" class="threshold-control" background="#ffffff"
							color="#ac7d61" :min="10" :max="60" v-model="yuzhi" />


						<!-- 自动控制 -->
						<text class="control-label">自动控制</text>
						<view class="toggle-container">
							<text class="toggle-text" :style="{color: !isfanon ? '#294b66' : '#999'}">关闭</text>
							<switch class="toggle-switch" :checked="isfanon" @change="togglefan" color="#294b66" />
							<text class="toggle-text" :style="{color: isfanon ? '#294b66' : '#999'}">开启</text>
						</view>
					</view>
				</view>

				<!-- 第二行：40%高度，滑动条 -->
				<view class="slider-row">
					<text class="control-label-">风速档位</text>
					<view class="slider-track">

						<slider class="fan-slider" v-model="environment.fan_level" min="0" max="3" step="1"
							activeColor="#ac7d61" @change="onSliderChange" />
					</view>
					<view class="slider-labels">
						<text v-for="i in 4" :key="i" class="slider-label">{{i-1}}</text>
					</view>
				</view>


			</view>


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
					password: "",
					role: "expert"
				},
				environment: {
					auto_fan: "",
					auto_light: "",
					dark: "",
					fan_level: 0,
					gas: "",
					human: "",
					humidity: "",
					id: "",
					light: "",
					temperature: "",
				},
				isLogin: false,
				isRequesting: false,
				requestInterval: null,

				isfanon: true,
				yuzhi: 25,
			}
		},
		watch: {
			environment: {
				deep: true,
				handler(newVal) {
					this.checkConditionautofan(newVal.temperature, newVal.auto_fan);

				}
			}
		},
		methods: {
			async checkConditionautofan(temperature, auto_fan) {

				if (String(auto_fan) === '0' && this.environment.fan_level === 0 && !this.toastLock) {
					this.toastLock = true;

					if (temperature > this.yuzhi) {
						uni.showModal({
							title: '环境提示',
							content: '环境温度过高，请打开风扇',
							confirmText: '确认',
							showCancel: false, // 不显示取消按钮
							success: (res) => {
								if (res.confirm) {
									this.toastLock = false; // 用户点击确认后解锁
								}
							}
						});
					}
				}

			},
			async startMonitoring() {
				this.requestInterval = setInterval(this.fetchEnvironmentData, 2000);
				await this.fetchEnvironmentData(); // 立即获取一次数据
			},
			async yuzhichange(e) {

				uni.setStorageSync("hall_yuzhi", this.yuzhi);

				const requestData = {
					device: "device1",
				};

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/device/get',
					method: 'POST',
					data: requestData,
					header: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
					},
					timeout: 3000
				});
				console.log('设备现状:', res.data)

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

				if (res.statusCode === 200) {
					this.environment.before = res.data.value

					const dataString = String(res.data.value); // 确保转为字符串
					const prefix = dataString.substring(0, 6);
					const modifiedPart = String(this.yuzhi); // 要替换的内容
					const suffix = dataString.substring(8);

					this.environment.now = prefix + modifiedPart + suffix;

					const request = {
						device: "device1",
						value: this.environment.now,
					};

					console.log('发送的指令:', request)

					const [error, resu] = await uni.request({
						url: 'http://127.0.0.1:8081/device/update',
						method: 'POST',
						data: request,
						header: {
							'Content-Type': 'application/json',
							'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
						},
						timeout: 3000
					});
				} else {
					uni.showToast({
						title: `网络请求失败: ${res.statusCode}`,
						icon: "none"
					});
				}
			},

			gohall() {
				uni.navigateTo({
					url: '/pages/Expert/expert_hall'
				});
			},
			gocage() {
				uni.navigateTo({
					url: '/pages/Expert/expert_cage'
				});
			},

			async onSliderChange(e) {
				this.environment.fan_level = e.detail.value;
				console.log('当前风扇等级', this.environment.fan_level)

				const requestData = {
					device: "device1",
				};

				const [err, res] = await uni.request({
					url: 'http://127.0.0.1:8081/device/get',
					method: 'POST',
					data: requestData,
					header: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
					},
					timeout: 3000
				});
				console.log('设备现状:', res.data)

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

				if (res.statusCode === 200) {
					this.environment.before = res.data.value

					const dataString = String(res.data.value); // 确保转为字符串
					const prefix = dataString.substring(0, 1);
					const modifiedPart = this.environment.fan_level; // 要替换的内容
					const suffix = dataString.substring(2);

					this.environment.now = prefix + modifiedPart + suffix;

					const request = {
						device: "device1",
						value: this.environment.now,
					};

					console.log('发送的指令:', request)

					const [error, resu] = await uni.request({
						url: 'http://127.0.0.1:8081/device/update',
						method: 'POST',
						data: request,
						header: {
							'Content-Type': 'application/json',
							'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
						},
						timeout: 3000
					});
				} else {
					uni.showToast({
						title: `网络请求失败: ${res.statusCode}`,
						icon: "none"
					});
				}
			},

			async GetEnvironment() {
				if (this.isRequesting) {
					clearInterval(this.requestInterval);
					uni.setStorageSync("hall_env_status", false);
					this.isRequesting = false;
					this.environment.gas = ""
					this.environment.temperature = ""
					this.environment.humidity = ""
				} else {
					this.isRequesting = true;
					uni.setStorageSync("hall_env_status", true);
					// 立即获取一次数据
					await this.fetchEnvironmentData();
					// 然后设置定时器
					this.requestInterval = setInterval(this.fetchEnvironmentData, 2000);
				}
			},

			async fetchEnvironmentData() {
				const requestData = {
					room: "room1",
				};

				try {
					const [err, res] = await uni.request({
						url: 'http://127.0.0.1:8081/environments/name',
						method: 'POST',
						data: requestData,
						header: {
							'Content-Type': 'application/json'
						},
						timeout: 3000
					});

					if (err) {
						console.error("请求出错:", err);
						this.stopMonitoring();
					} else {
						this.environment = res.data.data;
						console.log("hall当前环境:", this.environment);
					}
				} catch (error) {
					console.error("请求异常:", error);
					this.stopMonitoring();
				}
			},

			stopMonitoring() {
				this.isRequesting = false;
				clearInterval(this.requestInterval);
				uni.showToast({
					title: '监测已停止',
					icon: 'none'
				});
			},

			async togglefan(e) {
				this.isfanon = e.detail.value;
				uni.setStorageSync("hall_isfanon", this.isfanon);
				uni.showToast({
					title: this.isfanon ? '风扇自动控制已启用' : '风扇自动控制已禁用',
					icon: 'none'
				});

				if (this.isfanon) {

					const requestData = {
						device: "device1",
					};

					const [err, res] = await uni.request({
						url: 'http://127.0.0.1:8081/device/get',
						method: 'POST',
						data: requestData,
						header: {
							'Content-Type': 'application/json',
							'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
						},
						timeout: 3000
					});
					console.log('设备现状:', res.data)

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

					if (res.statusCode === 200) {
						this.environment.before = res.data.value

						const dataString = String(res.data.value); // 确保转为字符串
						const prefix = '1';
						const modifiedPart = dataString.substring(1);; // 要替换的内容

						this.environment.now = prefix + modifiedPart;

						const request = {
							device: "device1",
							value: this.environment.now,
						};

						console.log('发送的指令:', request)

						const [error, resu] = await uni.request({
							url: 'http://127.0.0.1:8081/device/update',
							method: 'POST',
							data: request,
							header: {
								'Content-Type': 'application/json',
								'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
							},
							timeout: 3000
						});
					} else {
						uni.showToast({
							title: `网络请求失败: ${res.statusCode}`,
							icon: "none"
						});
					}
				} else {
					const requestData = {
						device: "device1",
					};

					const [err, res] = await uni.request({
						url: 'http://127.0.0.1:8081/device/get',
						method: 'POST',
						data: requestData,
						header: {
							'Content-Type': 'application/json',
							'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
						},
						timeout: 3000
					});
					console.log('设备现状:', res.data)

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

					if (res.statusCode === 200) {

						const dataString = String(res.data.value); // 确保转为字符串
						const prefix = '0';
						const modifiedPart = dataString.substring(1);; // 要替换的内容

						this.environment.now = prefix + modifiedPart;


						const request = {
							device: "device1",
							value: this.environment.now,
						};

						console.log('发送的指令', request)

						const [error, resu] = await uni.request({
							url: 'http://127.0.0.1:8081/device/update',
							method: 'POST',
							data: request,
							header: {
								'Content-Type': 'application/json',
								'Authorization': `Bearer ${uni.getStorageSync("expert_token").trim()}`
							},
							timeout: 3000
						});
					} else {
						uni.showToast({
							title: `网络请求失败: ${res.statusCode}`,
							icon: "none"
						});
					}
				}
			}
		},
		onHide() {
			if (this.requestInterval) {
				clearInterval(this.requestInterval);
				this.requestInterval = null;
			}
		},
		onUnload() {
			if (this.requestInterval) {
				clearInterval(this.requestInterval);
				this.requestInterval = null;
			}
		},
		onShow() {
			const token = uni.getStorageSync("expert_token");
			const userInfo = uni.getStorageSync("expert_user");

			const hall_env = uni.getStorageSync("hall_env_status");

			const curfan = uni.getStorageSync("hall_isfanon");

			const curyuzhi = uni.getStorageSync("hall_yuzhi");

			if (curyuzhi) {
				this.yuzhi = curyuzhi;
			}

			if (hall_env === true || hall_env === 'true') {
				this.isRequesting = true
			}
			if (this.isRequesting && !this.requestInterval) {
				console.log('需要重启监测');
				this.startMonitoring(); // 新增的独立方法
			}
			if (curfan === "true" || curfan === true) {
				this.isfanon = curfan
			}

			if (token) {
				this.isLogin = true;
				this.userform = userInfo;
			} else {
				uni.redirectTo({
					url: "/pages/Expert/expert_login"
				});
			}
		}
	};
</script>

<style scoped>
	.top-row {
		display: flex;
		height: 20vh;
		/* 占60%高度 */
	}

	.left-section {
		width: 60%;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding-left: 10rpx;
	}

	.right-section {
		width: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		padding-left: 50rpx;
		justify-content: space-evenly;
		/* 均匀分布空间 */
	}

	.fan-icon {
		width: 170rpx;
		height: 170rpx;
		margin-top: 20rpx;
	}

	.control-label {
		margin: 30rpx 0;
		margin-left: 49rpx;
		/* 标签左间距 */
		font-size: 30rpx;
		/* 字体大小 */
		font-weight: 550;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		/* 字体颜色 */
		margin-bottom: 8rpx;
		/* 与value的间距 */
	}

	.control-label- {
		display: block;
		/* 确保占据整行 */
		margin: 20rpx 0;
		text-align: center;
		/* 文字水平居中 */
		font-size: 30rpx;
		/* 字体大小 */
		font-weight: 550;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		/* 字体颜色 */
		margin-bottom: 0rpx;
		/* 与value的间距 */
	}

	.threshold-control {
		transform: scale(1.3);
		/* 整体放大 */
		margin: 15rpx 0 25rpx 43rpx;
	}

	/* 第二行滑动条 */
	.slider-row {
		height: 13vh;
		/* 占40%高度 */
		padding: 10rpx 28rpx;
	}

	.slider-title {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 30rpx;
		display: block;
		text-align: center;
	}

	.slider-track {
		position: relative;
		height: 2px;
		/* 轨道高度 */
		background: #e0e0e0;
		/* 轨道颜色 */
		border-radius: 4rpx;
		margin: 20rpx 17rpx;
		/* 调整上下左右间距 */
		margin-bottom: 30rpx;
	}

	/* 滑块控件 */
	.fan-slider {
		position: absolute;
		width: 100%;
		top: -16rpx;
		/* 关键调整：让滑块按钮中心对齐轨道 */
		margin: 0;
		padding: 0;

		/* 移除默认边距 */
		::v-deep .uni-slider-handle-wrapper {
			transform: translateY(0);
			/* 处理小程序组件内部偏移 */
		}
	}

	.slider-labels {
		display: flex;
		justify-content: space-between;
		margin-top: 0rpx;
	}

	.slider-label {
		font-size: 28rpx;
		color: #294b66;
		width: 50rpx;
		text-align: center;
	}

	/* 大卡片整体样式 */
	.big-card {
		padding: 40rpx;
		border-radius: 24rpx;
		margin-bottom: 40rpx;
	}

	/* 主标题样式 */
	.card-main-title {
		display: block;
		text-align: center;
		font-size: 36rpx;
		/* 字体大小 */
		font-weight: 600;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		margin-bottom: 18rpx;
	}

	/* 图标行样式 */
	.card-icon-row {
		display: flex;
		justify-content: center;
		margin-bottom: 8rpx;
	}

	.center-icon {
		width: 140rpx;
		height: 140rpx;
	}

	/* 底部双列布局 */
	.dual-column-bottom {
		display: flex;
		justify-content: space-between;
		width: 100%;
		margin-top: 10rpx;
	}

	.left-column,
	.right-column {
		display: flex;
		flex-direction: column;
		align-items: center;
		/* 内容居中 */
		width: 48%;
		/* 留出间距 */
	}



	.bolabel {
		font-size: 30rpx;
		/* 字体大小 */
		font-weight: 500;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		/* 字体颜色 */
		margin-bottom: 7rpx;
		/* 与value的间距 */
	}

	/* 容器样式 - 调整内边距和背景 */
	.container {
		padding: 30rpx;
		/* 整体内边距调整 */
		padding-bottom: 120rpx;
		/* 底部留更多空间 */
		min-height: 100vh;
		position: relative;
		overflow: hidden;
		background:
			linear-gradient(rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.1)),
			url('/static/background/3.jpg');
	}

	.container::before {
		content: "";
		display: block;
		height: 70rpx;
	}

	/* 磨砂玻璃效果 - 关键样式 */
	.frost-glass {
		background: rgba(255, 255, 255, 0.65);
		/* 调整透明度 */
		backdrop-filter: blur(5px);
		/* 模糊程度 */
		-webkit-backdrop-filter: blur(15px);
		border: 1px solid rgba(255, 255, 255, 0.9);
		/* 更细的边框 */
		box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.25);
		/* 更柔和的阴影 */
	}

	/* 大卡片样式 - 调整圆角和内边距 */
	.big-card {
		border-radius: 24rpx;
		/* 圆角大小 */
		padding: 35rpx;
		/* 内边距 */
		margin-bottom: 40rpx;
		/* 下边距 */
	}

	/* 三指标卡片行 - 调整间距 */
	.metrics-row {
		display: flex;
		justify-content: space-between;
		margin-bottom: 40rpx;
	}

	.metric-label,
	.metric-value {
		width: 100%;
		text-align: center;
		display: block
	}

	.metric-label {
		font-size: 28rpx;
		/* 字体大小 */
		font-weight: 550;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		/* 字体颜色 */
		margin-bottom: 8rpx;
		/* 与value的间距 */
	}

	.metric-value {
		font-size: 30rpx;
		/* 比label稍大 */
		font-weight: bold;
		/* 加粗 */
		color: #2d2d2d;
		/* 深色更突出 */
	}

	/* 单个指标卡片 - 宽度和间距调整 */
	.metric-card {
		border-radius: 20rpx;
		padding: 35rpx 10rpx; // 前者控制三个小卡片高度
		width: calc(33.33% - 10rpx);
		/* 计算宽度考虑间距 */
		text-align: center;
		flex-direction: column;
	}

	.toggle-text {
		font-size: 26rpx;
		transition: color 0.3s;
	}

	/* 底部双卡片行 */
	.bottom-row {
		display: flex;
		justify-content: space-between;
		margin-bottom: 30rpx;
	}

	/* 底部卡片样式 */
	.bottom-card {
		padding: 30rpx 10rpx;
		border-radius: 20rpx;
		padding: 25rpx 0;
		width: 47.5%;
		text-align: center;
	}


	/* 新增的正方形卡片行 */
	.square-row {
		display: flex;
		justify-content: space-between;
	}

	/* 正方形卡片样式 */
	.square-card {
		width: 48%;
		/* 宽度 */
		aspect-ratio: 1/1;
		/* 保持1:1正方形 */
		border-radius: 20rpx;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 20rpx;
	}

	/* 正方形卡片文字 */
	.square-text {
		font-size: 32rpx;
		color: #333;
		margin-bottom: 20rpx;
		text-align: center;
	}

	/* 正方形卡片图标 - 可点击 */
	.square-icon {
		width: 80rpx;
		height: 80rpx;
		transition: transform 0.3s;
	}

	.toggle-text {
		font-size: 28rpx;
		/* 文字大小 */
		font-weight: 550;
		/* 文字粗细（400=正常，500=中等，600=半粗） */
		color: #999;
		/* 默认颜色 */
		transition: all 0.6s;
		/* 添加过渡效果 */
	}

	.entext {
		font-size: 30rpx;
		/* 文字大小 */
		font-weight: 550;
		/* 文字粗细（400=正常，500=中等，600=半粗） */
		color: #999;

	}

	.toggle-text:first-child {
		margin-right: 10rpx;
		/* "闭馆"与开关的间距 */
	}

	.toggle-text:last-child {
		margin-left: 1rpx;
		/* "开馆"与开关的间距 */
	}

	.square-icon:active {
		transform: scale(0.9);
		/* 点击效果 */
	}

	.centered-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
	}


	/* 调整图标大小 */
	.card-icon {
		width: 120rpx;
		/* 增大图标 */
		height: 120rpx;

	}

	.goicon {
		width: 100rpx;
		/* 增大图标 */
		height: 110rpx;

	}

	.metric-icon {
		width: 70rpx;
		/* 增大指标图标 */
		height: 70rpx;
		margin: 0 auto 15rpx;
	}

	.bottom-icon {
		width: 90rpx;
		/* 增大底部图标 */
		height: 90rpx;
		margin-bottom: 20rpx;
	}

	/* 双列布局容器 */
	.dual-column-layout {
		display: flex;
		width: 100%;
	}

	/* 图标列样式 */
	.icon-column {
		width: 33.33%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.column-icon {
		width: 130rpx;
		height: 130rpx;
	}

	/* 内容列样式 */
	.content-column {
		width: 66.67%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		/* 增加左侧内边距 */
		text-align: center;
	}

	.column-title {
		font-size: 39rpx;
		font-weight: 600;
		color: #565656;
		margin-bottom: 10rpx;
	}

	/* use */

	.column-btn {
		background-color: #ac7d61;
		color: white;
		border-radius: 50rpx;
		font-size: 27rpx;
		padding: 0rpx 0rpx;
		width: 370rpx;
		border: none;
		transition: all 0.3s;
	}

	.column-btn.active-state {
		background-color: #cdcdcd;
	}
</style>