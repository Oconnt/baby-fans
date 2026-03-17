<template>
  <view class="register-page">
    <view class="header">
      <text class="title">注册账号</text>
    </view>

    <view class="card register-card">
      <view class="form-group">
        <text class="label">身份选择</text>
        <view class="role-selector">
          <view
            class="role-option"
            :class="{ active: form.role === 'parent' }"
            @click="form.role = 'parent'"
          >
            <text class="role-icon">👨‍👩‍👧</text>
            <text>家长</text>
          </view>
          <view
            class="role-option"
            :class="{ active: form.role === 'child' }"
            @click="form.role = 'child'"
          >
            <text class="role-icon">👶</text>
            <text>孩子</text>
          </view>
        </view>
      </view>

      <view class="form-group">
        <text class="label">昵称 (选填)</text>
        <input type="text" v-model="form.nickname" placeholder="请输入昵称" style="width: 100%; padding: 20rpx; border: 1px solid #eee; border-radius: 12rpx; background: white; color: #333;" />
      </view>

      <button class="btn-primary" @click="handleRegister">注册</button>

      <view class="login-link">
        <text @click="goToLogin">已有账号？立即登录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { request } from '../../utils/request';

const form = ref({
  name: '',
  password: '',
  nickname: '',
  role: 'parent' // 默认家长
});

const handleRegister = async () => {
  // No validation needed since we only have role (required) and nickname (optional)
  try {
    const res = await request({
      url: '/register',
      method: 'POST',
      data: form.value
    });
    uni.showToast({ title: '注册成功', icon: 'success' });
    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/register-success/register-success?code=${res.login_code}`
      });
    }, 1500);
  } catch (e: any) {
    uni.showToast({ title: e?.error || '注册失败', icon: 'none' });
  }
};

const goToLogin = () => {
  uni.redirectTo({ url: '/pages/login/login' });
};
</script>

<style lang="scss" scoped>
.register-page {
  min-height: 100vh;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 40rpx;
}

.header {
  margin-bottom: 60rpx;
  .title {
    font-size: 48rpx;
    font-weight: bold;
    color: #333;
  }
}

.register-card {
  width: 100%;
  padding: 40rpx;
  border-radius: 24rpx;
  background: white;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}

.form-group {
  margin-bottom: 30rpx;
  .label {
    font-size: 28rpx;
    color: #333;
    font-weight: bold;
    margin-bottom: 10rpx;
    display: block;
  }
}

.role-selector {
  display: flex;
  gap: 20rpx;
  .role-option {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20rpx;
    border: 2rpx solid #eee;
    border-radius: 16rpx;
    &.active {
      border-color: var(--orange);
      background: #fff3e0;
    }
    .role-icon { font-size: 48rpx; margin-bottom: 10rpx; }
  }
}

.input-box {
  border: 2rpx solid #eee;
  background: #ffffff;
  border-radius: 12rpx;
  padding: 20rpx;
  font-size: 28rpx;
  width: 100%;
  box-sizing: border-box;
  color: #333333;
  caret-color: #333333;
}

.btn-primary {
  background: var(--orange);
  color: white;
  border-radius: 40rpx;
  height: 80rpx;
  line-height: 80rpx;
  font-size: 32rpx;
  font-weight: bold;
  margin-top: 20rpx;
  &::after { border: none; }
}

.login-link {
  text-align: center;
  margin-top: 30rpx;
  text { color: #999; font-size: 24rpx; text-decoration: underline; }
}
</style>
