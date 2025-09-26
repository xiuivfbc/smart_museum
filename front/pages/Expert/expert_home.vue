<template>
	<view class="container">
		<view v-if="isLogin">
			<view class="blur-background"></view>

			<view class="big-card frost-glass">
				<view class="dual-column-layout">
					<view class="icon-column">
						<image class="column-icon" src="/static/icons/monitor.png" mode="aspectFit"></image>
					</view>

					<view class="content-column">
						<text class="column-title">环境监测控制面板</text>
						<button class="column-btn" :class="{'active-state': isRequesting}" @click="GetEnvironment">
							{{ isRequesting ? '实时监测中' : '开启监测' }}
						</button>
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

			<!-- 一键开闭馆卡片 -->
			<view class="bottom-row">
				<view class="card-content bottom-card frost-glass">
					<view class="card-right" style="margin-top: 0rpx;">
						<text class="moretext">场馆状态</text>
					</view>
					<image class="card-icon" style="margin-top: 10rpx;"
						:src="isFacilityOn ? '/static/icons/omuseum.png' : '/static/icons/cmuseum.png'"
						mode="aspectFit"></image>
					<view class="card-right" style="margin-top: 10rpx;">
						<view class="toggle-container">
							<text class="toggle-text" :style="{color: !isFacilityOn ? '#294b66' : '#999'}">闭馆</text>
							<switch class="toggle-switch" :checked="isFacilityOn" @change="toggleFacility"
								color="#294b66" />
							<text class="toggle-text" :style="{color: isFacilityOn ? '#294b66' : '#999'}">开馆</text>
						</view>
					</view>
				</view>

				<view class="card-content bottom-card frost-glass">
					<view class="card-right" style="margin-top: 0rpx;">
						<text class="moretext">控制模式</text>
					</view>
					<image class="card-icon" style="margin-top: 10rpx;"
						:src="isauto ? '/static/icons/uauto.png' : '/static/icons/auto.png'" mode="aspectFit"></image>
					<view class="card-right" style="margin-top: 10rpx;">
						<view class="toggle-container">
							<text class="toggle-text" :style="{color: !isauto ? '#294b66' : '#999'}">自动</text>
							<switch class="toggle-switch" :checked="isauto" @change="changeauto" color="#294b66" />
							<text class="toggle-text" :style="{color: isauto ? '#294b66' : '#999'}">手动</text>
						</view>
					</view>
				</view>
			</view>

			<view class="bottom-row">
				<view class="card-content bottom-card frost-glass" @click="gohall">
					<image class="goicon" src="/static/icons/hall.png" mode="aspectFit"></image>
					<view class="card-right" style="margin-top: 10rpx;">
						<text class="entext">前往展厅</text>
					</view>
				</view>

				<view class="card-content bottom-card frost-glass" @click="gocage">
					<image class="goicon" src="/static/icons/cage.png" mode="aspectFit"></image>
					<view class="card-right" style="margin-top: 10rpx;">
						<text class="entext">前往仓库</text>
					</view>
				</view>
			</view>

			<expert_tabbar />
		</view>
	</view>
</template>

<script>
	import expert_tabbar from '../Common/expert_tabbar.vue';

	export default {
		components: {
			expert_tabbar
		},
		data() {
			return {
				nvueWidth: 730,
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
					fan_level: "",
					gas: "",
					human: "",
					humidity: "",
					id: "",
					light: "",
					temperature: "",
					before: "",
					now: "",
				},
				isLogin: false,
				isRequesting: false,
				homerequest: null,
				isFacilityOn: false,
				isauto: false,
				
				
			}
		},
		watch: {
			environment: {
				deep: true, // 监听所有子属性变化
				handler(newVal) {
					this.checkConditionthief(newVal.light, newVal.human);
					this.checkConditionled(newVal.dark);
				}
			}
		},
		methods: {
			checkConditionthief(light, human) {

				if (String(light) === '1' && String(human) === '1' && !this.toastLock) {
					this.toastLock = true;

					uni.showModal({
						title: '安全警告',
						content: '检测到异常人员活动',
						confirmText: '确认',
						showCancel: false, // 不显示取消按钮
						success: (res) => {
							if (res.confirm) {
								this.toastLock = false; // 用户点击确认后解锁
							}
						}
					});

				}
			},
			checkConditionled(dark) {
			
				if (String(this.environment.auto_light) === '1' && String(dark) === '1' && !this.toastLock) {
					this.toastLock = true;
			
					uni.showModal({
						title: '环境提示',
						content: '光线过暗，灯光自动开启',
						confirmText: '确认',
						showCancel: false, // 不显示取消按钮
						success: (res) => {
							if (res.confirm) {
								this.toastLock = false; // 用户点击确认后解锁
							}
						}
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
			async GetEnvironment() {
				if (this.isRequesting) {
					clearInterval(this.homerequest);
					this.isRequesting = false;
					uni.setStorageSync("home_env_status", this.isRequesting);
					this.environment.gas = ""
					this.environment.temperature = ""
					this.environment.humidity = ""
				} else {
					this.isRequesting = true;
					uni.setStorageSync("home_env_status", this.isRequesting);
					// 立即获取一次数据
					await this.fetchEnvironmentData();
					// 然后设置定时器
					this.homerequest = setInterval(this.fetchEnvironmentData, 2000);
				}
			},
			async startMonitoring() {
			  this.homerequest = setInterval(this.fetchEnvironmentData, 2000);
			  await this.fetchEnvironmentData(); // 立即获取一次数据
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
						console.log("home当前环境:", this.environment);
					}
				} catch (error) {
					console.error("请求异常:", error);
					this.stopMonitoring();
				}
			},
			stopMonitoring() {
				this.isRequesting = false;
				clearInterval(this.homerequest);
				uni.showToast({
					title: '监测已停止',
					icon: 'none'
				});
			},
			async toggleFacility(e) {
				this.isFacilityOn = e.detail.value;
				console.log('场馆状态', this.isFacilityOn)
				uni.setStorageSync("home_isclose", this.isFacilityOn);
				uni.showToast({
					title: this.isFacilityOn ? '场馆已开启' : '场馆已关闭',
					icon: 'none'
				});
				if (this.isFacilityOn) {

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
						const prefix = dataString.substring(0, 4); // 前4位不变
						const modifiedPart = '20'; // 要替换的内容
						const tmp = dataString.substring(6, 8);
						const ttt = '1'; // 要替换的内容
						const suffix = dataString.substring(9);

						this.environment.now = prefix + modifiedPart + tmp + ttt + suffix;

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

						const dataString = String(res.data.value);
						const prefix = dataString.substring(0, 3); // 前4位不变
						const modifiedPart = '011'; // 要替换的内容
						const tmp = dataString.substring(6, 8);
						const ttt = '0'; // 要替换的内容
						const suffix = dataString.substring(9);

						this.environment.now = prefix + modifiedPart + tmp + ttt + suffix;

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
			},
			async changeauto(e) {
				this.isauto = e.detail.value;
				uni.setStorageSync("home_isauto", this.isauto);
				uni.showToast({
					title: this.isauto ? '手动控制已开启' : '自动控制已开启',
					icon: 'none'
				});

				if (!this.isauto) {

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
						const prefix = '1'; // 前4位不变
						const modifiedPart = dataString.substring(1, 3);; // 要替换的内容
						const tmp = '1';

						const suffix = dataString.substring(4);

						this.environment.now = prefix + modifiedPart + tmp + suffix;

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
						const prefix = '01'; // 前4位不变
						const modifiedPart = dataString.substring(2, 3);; // 要替换的内容
						const tmp = '0';

						const suffix = dataString.substring(4);

						this.environment.now = prefix + modifiedPart + tmp + suffix;

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
				}


				// 这里可以添加实际控制设备的逻辑
			}
		},

		onShow() {
			const token = uni.getStorageSync("expert_token");
			const userInfo = uni.getStorageSync("expert_user");

			const status = uni.getStorageSync("home_env_status");

			const curfa = uni.getStorageSync("home_isclose");
			const curauto = uni.getStorageSync("home_isauto");

			if (status === true || status === 'true') {
				this.isRequesting = true
			}
			
			this.isFacilityOn = curfa === true || curfa === "true";

			this.isauto = curauto === true || curauto === "true";

			if (token) {
				this.isLogin = true;
				this.userform = userInfo;
			} else {
				uni.redirectTo({
					url: "/pages/Expert/expert_login"
				});
			}
			
			if (this.isRequesting && !this.homerequest) {
			    console.log('需要重启监测');
			    this.startMonitoring(); // 新增的独立方法
			  }
		},
		
		onHide() {
			if (this.homerequest) {
				clearInterval(this.homerequest);
				this.homerequest = null;
			}
		},
		onUnload() {
			if (this.homerequest) {
				clearInterval(this.homerequest);
				this.homerequest = null;
			}
		},
	};
</script>

<style scoped>
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
		padding: 40rpx;
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

	.moretext {
		display: block;
		text-align: center;
		font-size: 33rpx;
		/* 字体大小 */
		font-weight: 600;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		margin-bottom: 6rpx;

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
		font-size: 38rpx;
		font-weight: 600;
		color: #333;
		margin-bottom: 40rpx;
	}

	.column-btn {
		background-color: #ac7d61;
		color: white;
		border-radius: 50rpx;
		font-size: 27rpx;
		padding: 6rpx 0rpx;
		width: 400rpx;
		border: none;
		transition: all 0.3s;
	}

	.column-btn.active-state {
		background-color: #cdcdcd;
	}
</style>