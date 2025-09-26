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

						<text class="column-title"> 仓库区域控制面板</text>
						<button class="column-btn" :class="{'active-state': isRequesting}" @click="GetEnvironment">
							{{ isRequesting ? '实时监测中' : '开启监测' }}
						</button>
					</view>
					<!-- 左侧图标列 (1/3宽度) -->
					<view class="icon-column">

						<image class="column-icon" src="/static/icons/cage.png" mode="aspectFit"></image>
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


			<view class="bottom-row">
				<view class="card-content bottom-card frost-glass">
					<view class="card-right" style="margin-top: 0rpx;">
						<text class="moretext">红外检测</text>
					</view>
					<image class="card-icon"
						:src="islexon ? '/static/icons/hongwaion.png' : '/static/icons/hongwaioff.png'"
						mode="aspectFit"></image>
					<view class="card-right" style="margin-top: 10rpx;">
						<view class="toggle-container">
							<text class="toggle-text" :style="{color: !islexon ? '#294b66' : '#999'}">禁用</text>
							<switch class="toggle-switch" :checked="islexon" @change="togglelex" color="#294b66" />
							<text class="toggle-text" :style="{color: islexon ? '#294b66' : '#999'}">启用</text>
						</view>
					</view>
				</view>

				<view class="card-content bottom-card frost-glass">
					<view class="card-right" style="margin-top: 0rpx;">
						<text class="moretext">烟雾检测</text>
					</view>
					<image class="card-icon"
						:src="issmokeon ? '/static/icons/smokeon.png' : '/static/icons/smokeoff.png'" mode="aspectFit">
					</image>
					<view class="card-right" style="margin-top: 10rpx;">
						<view class="toggle-container">
							<text class="toggle-text" :style="{color: !issmokeon ? '#294b66' : '#999'}">禁用</text>
							<switch class="toggle-switch" :checked="issmokeon" @change="togglesmoke" color="#294b66" />
							<text class="toggle-text" :style="{color: issmokeon ? '#294b66' : '#999'}">启用</text>
						</view>
					</view>
				</view>
			</view>

	

			<view class="big-card frost-glass">

				<text class="card-main-title">湿度调节系统</text>
			

				<view class="card-icon-row">
					<image class="center-icon" src="/static/icons/wateron.png" mode="aspectFit"></image>
				</view>

				<view class="dual-column-bottom">

					<view class="left-column">
			
						<text class="bolabel">开启阈值</text>
						<uni-card :is-shadow="true" margin=1rpx padding=9rpx>
							<uni-number-box @change="yuzhichange"
							 background="#ffffff" color="#ac7d61" :min="0" :max="99"
								v-model="yuzhi" />
						</uni-card>
			
					</view>
			

					<view class="right-column">
						<text class="bolabel">自动控制</text>
						<view class="toggle-container" style="margin-top: 6rpx;">
							<text class="toggle-text" :style="{color: !iswateron ? '#294b66' : '#999'}">禁用</text>
							<switch class="toggle-switch" :checked="iswateron" @change="togglewater" color="#294b66" />
							<text class="toggle-text" :style="{color: iswateron ? '#294b66' : '#999'}">启用</text>
						</view>
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
				},
				isLogin: false,
				isRequesting: false,
				requestInterval: null,

				islexon: false,
				iswateron: false,
				isledon: false,
				issmokeon: false,
				
				yuzhi:30
			}
		},
		watch: {
			environment: {
				deep: true, 
				handler(newVal) {
					this.checkConditiongas(newVal.gas);
				}
			}
		},
		methods: {
			async checkConditiongas(gas){
				if (gas > 700 && !this.toastLock) {
					this.toastLock = true;

						uni.showModal({
							title: '环境提示',
							content: '气体浓度过高，蜂鸣器已报警',
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
			async yuzhichange(e) {
			
				uni.setStorageSync("cage_yuzhi", this.yuzhi);
			
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
					const prefix = dataString.substring(0, 12);
					const modifiedPart = String(this.yuzhi); // 要替换的内容

			
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
			async startMonitoring() {
			  this.requestInterval = setInterval(this.fetchEnvironmentData, 2000);
			  await this.fetchEnvironmentData(); // 立即获取一次数据
			},
			async GetEnvironment() {
				if (this.isRequesting) {
					clearInterval(this.requestInterval);
					uni.setStorageSync("cage_env_status", false);
					this.isRequesting = false;
					this.environment.gas = ""
					this.environment.temperature = ""
					this.environment.humidity = ""
				} else {
					this.isRequesting = true;
					uni.setStorageSync("cage_env_status", true);
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
						console.log("请求成功:", res.data);
						this.environment = res.data.data;
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

			async togglewater(e) {
				this.iswateron = e.detail.value;
				uni.setStorageSync("cage_iswater", this.iswateron);

				uni.showToast({
					title: this.iswateron ? '自动控制已开启' : '自动控制已关闭',
					icon: 'none'
				});


				if (this.iswateron) {

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
						const prefix = dataString.substring(0, 10);
						const modifiedPart = '1'; // 要替换的内容
						const suffix = dataString.substring(11);

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
						const prefix = dataString.substring(0, 10);
						const modifiedPart = '0'; // 要替换的内容
						const suffix = dataString.substring(11);

						this.environment.now = prefix + modifiedPart + suffix;

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
			async togglelex(e) {
				this.islexon = e.detail.value;
				uni.setStorageSync("cage_islex", this.islexon);

				uni.showToast({
					title: this.islexon ? '红外检测已开启' : '红外检测已关闭',
					icon: 'none'
				});


				if (this.islexon) {

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
						const prefix = dataString.substring(0, 5);
						const modifiedPart = '1'; // 要替换的内容
						const suffix = dataString.substring(6);

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
						const prefix = dataString.substring(0, 5);
						const modifiedPart = '0'; // 要替换的内容
						const suffix = dataString.substring(6);

						this.environment.now = prefix + modifiedPart + suffix;

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
			async togglesmoke(e) {
				this.issmokeon = e.detail.value;
				uni.setStorageSync("cage_issmoke", this.issmokeon);
				uni.showToast({
					title: this.issmokeon ? '烟雾报警已启用' : '烟雾报警已禁用',
					icon: 'none'
				});

				if (this.issmokeon) {

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
						const prefix = dataString.substring(0, 2);
						const modifiedPart = '1'; // 要替换的内容
						const suffix = dataString.substring(3);

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
						const prefix = dataString.substring(0, 2);
						const modifiedPart = '0'; // 要替换的内容
						const suffix = dataString.substring(3);

						this.environment.now = prefix + modifiedPart + suffix;


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
			
			const cage_env = uni.getStorageSync("cage_env_status");

			const curlex = uni.getStorageSync("cage_islex");
			const cursmoke = uni.getStorageSync("cage_issmoke");
			const curwater= uni.getStorageSync("cage_iswater");

			const curyuzhi = uni.getStorageSync("cage_yuzhi");
			
			if (curyuzhi) {
				this.yuzhi = curyuzhi;
			}
			
			if (cage_env === true || cage_env === 'true') {
				this.isRequesting = true
			}
			if (this.isRequesting && !this.requestInterval) {
			    console.log('需要重启监测');
			    this.startMonitoring(); // 新增的独立方法
			  }
			this.islexon = curlex === true || curlex === "true";
			
			this.issmokeon = cursmoke === true || cursmoke === "true";
			
			this.iswateron= curwater === true || curwater === "true";

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

	.moretext {
		display: block;
		text-align: center;
		font-size: 33rpx;
		/* 字体大小 */
		font-weight: 600;
		/* 字体粗细（400=正常，600=半粗，700=粗体） */
		color: #666;
		margin-bottom: 14rpx;
	
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