<template>
  <view class="success-page">
    <view class="card success-card">
      <text class="icon">🎉</text>
      <text class="title">注册成功</text>
      <text class="desc">您的登录码为</text>

      <view class="code-box">
        <text class="code">{{ loginCode }}</text>
      </view>

      <button class="btn-copy" @click="handleCopy">复制登录码</button>
      <view class="link" @click="goToLogin">返回登录页</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';

const loginCode = ref('');

onLoad((options) => {
  if (options.code) {
    loginCode.value = options.code;
  }
});

const handleCopy = () => {
  uni.setClipboardData({
    data: loginCode.value,
    success: () => {
      uni.showToast({ title: '已复制', icon: 'success' });
    }
  });
};

const goToLogin = () => {
  uni.redirectTo({
    url: `/pages/login/login?code=${loginCode.value}`
  });
};
</script>

<style lang="scss" scoped>
.success-page {
  min-height: 100vh;
  background: linear-gradient(160deg, #FF6B35 0%, #FF8FAB 40%, #C77DFF 80%, #4CC9F0 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-card {
  width: 80%;
  padding: 60rpx 40rpx;
  border-radius: 32rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 10rpx 30rpx rgba(0,0,0,0.2);

  .icon { font-size: 100rpx; margin-bottom: 20rpx; }
  .title { font-size: 48rpx; font-weight: bold; color: #333; margin-bottom: 10rpx; }
  .desc { font-size: 28rpx; color: #666; margin-bottom: 30rpx; }

  .code-box {
    background: #f8f9fa;
    padding: 20rpx 40rpx;
    border-radius: 16rpx;
    margin-bottom: 40rpx;
    border: 2rpx dashed #ccc;
    .code { font-size: 60rpx; font-weight: bold; color: #FF6B35; letter-spacing: 8rpx; }
  }

  .btn-copy {
    width: 100%;
    background: #FF6B35;
    color: white;
    border-radius: 40rpx;
    height: 80rpx;
    line-height: 80rpx;
    font-weight: bold;
    margin-bottom: 20rpx;
    &::after { border: none; }
  }

  .link {
    color: #999;
    font-size: 24rpx;
    text-decoration: underline;
  }
}
</style>
