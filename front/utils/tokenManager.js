// Token管理工具函数
class TokenManager {
    // 获取access token
    static getAccessToken() {
        return uni.getStorageSync('access_token');
    }

    // 获取refresh token
    static getRefreshToken() {
        return uni.getStorageSync('refresh_token');
    }

    // 设置tokens
    static setTokens(accessToken, refreshToken) {
        uni.setStorageSync('access_token', accessToken);
        uni.setStorageSync('refresh_token', refreshToken);
    }

    // 清除所有tokens
    static clearTokens() {
        uni.removeStorageSync('access_token');
        uni.removeStorageSync('refresh_token');
        uni.removeStorageSync('public_user');
        uni.removeStorageSync('public_token'); // 清除旧版本token
    }

    // 刷新token
    static async refreshAccessToken() {
        const refreshToken = this.getRefreshToken();
        if (!refreshToken) {
            throw new Error('No refresh token available');
        }

        try {
            const [err, res] = await uni.request({
                url: 'http://127.0.0.1:8081/auth/refresh',
                method: 'POST',
                data: {
                    refresh_token: refreshToken
                },
                header: {
                    'Content-Type': 'application/json'
                },
                timeout: 5000
            });

            if (err) {
                throw new Error('Network error during token refresh');
            }

            if (res.statusCode === 200 && res.data.access_token) {
                this.setTokens(res.data.access_token, res.data.refresh_token || refreshToken);
                return res.data.access_token;
            } else {
                throw new Error('Token refresh failed');
            }
        } catch (error) {
            console.error('Token refresh error:', error);
            // 刷新失败，清除所有token并跳转到登录页
            this.clearTokens();
            uni.redirectTo({
                url: '/pages/Public/public_login'
            });
            throw error;
        }
    }

    // 带自动刷新的请求函数
    static async requestWithAuth(options) {
        let accessToken = this.getAccessToken();

        // 添加Authorization头
        options.header = options.header || {};
        options.header['Authorization'] = `Bearer ${accessToken}`;

        const [err, res] = await uni.request(options);

        // 如果是401错误（token过期），尝试刷新token并重试
        if (res && res.statusCode === 401) {
            try {
                accessToken = await this.refreshAccessToken();
                options.header['Authorization'] = `Bearer ${accessToken}`;
                return await uni.request(options);
            } catch (refreshError) {
                // 刷新失败，直接返回原始错误
                return [err, res];
            }
        }

        return [err, res];
    }

    // 登出函数
    static async logout() {
        const refreshToken = this.getRefreshToken();

        if (refreshToken) {
            try {
                await uni.request({
                    url: 'http://127.0.0.1:8081/auth/logout',
                    method: 'POST',
                    data: {
                        refresh_token: refreshToken
                    },
                    header: {
                        'Content-Type': 'application/json'
                    },
                    timeout: 5000
                });
            } catch (error) {
                console.error('Logout request failed:', error);
            }
        }

        // 清除本地存储的token
        this.clearTokens();

        // 跳转到登录页
        uni.redirectTo({
            url: '/pages/Public/public_login'
        });
    }
}

module.exports = TokenManager;