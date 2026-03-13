<template>
  <view class="login-page">
    <view class="header">
      <text class="title fredoka">Baby-Fans</text>
      <text class="subtitle">宝贝积分管理系统 v1.2</text>
    </view>

    <!-- Unified Login Card -->
    <view class="card login-card">
      <text class="prompt">使用登录码进入</text>
      <input class="input-box" v-model="loginCode" type="number" placeholder="6 位数字" maxlength="6" />
      <button class="btn-primary" @click="handleCodeLogin">立即登录 ➔</button>

      <view class="divider">
        <view class="line"></view>
        <text class="divider-text">其他登录方式</text>
        <view class="line"></view>
      </view>

      <!-- WeChat Login Button -->
      <button class="btn-wechat" @click="handleWechatLogin">
        <text class="icon">💚</text>
        <text>微信一键登录</text>
      </button>

      <view class="demo-hints">
        <text class="hint-label">体验账号:</text>
        <text class="hint-btn parent" @click="loginCode = '888888'">家长端</text>
        <text class="hint-btn child" @click="loginCode = '666666'">孩子端</text>
      </view>
    </view>

    <text class="decoration decor1">🍭</text>
    <text class="decoration decor2">🎨</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { request } from '../../utils/request';

const loginCode = ref('');

const handleCodeLogin = async () => {
  if (loginCode.value.length < 6) return uni.showToast({ title: '码太短啦', icon: 'none' });

  try {
    const res = await request({
      url: `/login/code?code=${loginCode.value}`,
      method: 'GET'
    });

    uni.setStorageSync('token', res.token);
    uni.setStorageSync('userInfo', JSON.stringify({
      id: res.user_id,
      role: res.role,
      nickname: res.role === 'parent' ? '超级家长' : '勤劳宝宝',
      points: 100
    }));

    uni.showToast({ title: '欢迎回来 ✨', icon: 'success' });
    setTimeout(() => {
      // Navigate to the 'Home' tab
      uni.switchTab({ url: '/pages/index/index' });
    }, 1000);
  } catch (e) {
    uni.showToast({ title: '登录码不对哦', icon: 'none' });
  }
};

const handleWechatLogin = () => {
  uni.login({
    provider: 'weixin',
    success: async (loginRes) => {
      try {
        const res = await request({
          url: '/api/v1/auth/wechat/login',
          method: 'POST',
          data: {
            code: loginRes.code,
            role: 'parent'
          }
        });
        uni.setStorageSync('token', res.token);
        uni.setStorageSync('userInfo', JSON.stringify(res.user));
        uni.showToast({ title: '微信登录成功', icon: 'success' });
        setTimeout(() => {
          uni.switchTab({ url: '/pages/index/index' });
        }, 1000);
      } catch (e) {
        uni.showToast({ title: '微信授权失败', icon: 'none' });
      }
    }
  });
};
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(160deg, #FF6B35 0%, #FF8FAB 40%, #C77DFF 80%, #4CC9F0 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 40rpx;
  color: white;
  position: relative;
  overflow: hidden;
}

.header {
  text-align: center;
  margin-top: 60rpx;
  margin-bottom: 80rpx;
  .title {
    font-size: 88rpx;
    display: block;
    text-shadow: 0 8rpx 20rpx rgba(0,0,0,0.1);
  }
  .subtitle {
    font-size: 24rpx;
    font-weight: 700;
    opacity: 0.8;
  }
}

.login-card {
  width: 100%;
  margin: 0;
  padding: 60rpx 48rpx;
  color: var(--text);
  border-radius: 48rpx;
  box-shadow: 0 20rpx 60rpx rgba(0,0,0,0.15);

  .prompt {
    display: block;
    text-align: center;
    font-size: 28rpx;
    color: var(--text2);
    font-weight: 800;
    margin-bottom: 40rpx;
  }

  .input-box {
    border: 4rpx solid #F8F8F8;
    background: #FAFAFA;
    border-radius: 24rpx;
    padding: 30rpx;
    margin-bottom: 40rpx;
    text-align: center;
    font-size: 48rpx;
    font-family: 'Fredoka One';
    color: var(--orange);
  }

  .divider {
    display: flex;
    align-items: center;
    margin: 40rpx 0;
    .line { flex: 1; height: 2rpx; background: #EEE; }
    .divider-text { font-size: 22rpx; color: #BBB; padding: 0 20rpx; font-weight: 700; }
  }

  .btn-wechat {
    background: #fdfdfd;
    border: 2rpx solid #EEE;
    color: #333;
    border-radius: 20rpx;
    font-size: 28rpx;
    font-weight: 800;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10rpx 0;
    &::after { border: none; }
    .icon { margin-right: 12rpx; font-size: 32rpx; }
  }

  .demo-hints {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 20rpx;
    margin-top: 50rpx;
    .hint-label { font-size: 20rpx; color: #CCC; }
    .hint-btn {
      font-size: 20rpx;
      font-weight: 800;
      padding: 6rpx 16rpx;
      border-radius: 10rpx;
      &.parent { color: var(--purple); background: #F3E5F5; }
      &.child { color: var(--orange); background: #FFF3E0; }
    }
  }
}

.decoration {
  position: absolute;
  font-size: 100rpx;
  opacity: 0.3;
  &.decor1 { top: 12%; left: 8%; transform: rotate(-20deg); }
  &.decor2 { bottom: 12%; right: 8%; transform: rotate(15deg); }
}
</style>
