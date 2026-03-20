<template>
  <view class="container">
    <!-- Hero Section: Previous Style Points Card -->
    <view class="hero-card">
      <view class="bg-circle top-right"></view>
      <view class="bg-circle bottom-left"></view>

      <view class="user-header">
        <image class="avatar" :src="avatarUrl" mode="aspectFill" @click="changeAvatar" />
        <view class="user-main">
          <text class="nickname fredoka" @click="changeNickname">{{ userInfo.nickname || '点击设置昵称' }}</text>
          <view class="role-badge">{{ userInfo.role === 'parent' ? '家长管理端' : '宝贝端' }}</view>
        </view>
      </view>

      <!-- <view class="points-display" v-if="userInfo.role !== 'parent'">
        <text class="label">当前积分余额</text>
        <view class="amount-row">
          <text class="amount fredoka">{{ userInfo.points || 0 }}</text>
          <text class="unit">⭐</text>
        </view>
      </view>

      <text class="moto" v-if="userInfo.role !== 'parent'">加油！继续努力挣积分吧！</text>
      <text class="decoration" v-if="userInfo.role !== 'parent'">🐻</text> -->
    </view>

    <!-- Menu Section: Card Style -->
    <view class="menu-group">
      <view class="card menu-item" @click="goToRecords">
        <view class="menu-left">
          <view class="icon-box" style="background: #FFF9C4;">🎁</view>
          <text class="label">兑换记录</text>
        </view>
        <text class="arrow">➔</text>
      </view>

      <view class="card menu-item" v-if="userInfo.role === 'parent'" @click="goToPointsRecords">
        <view class="menu-left">
          <view class="icon-box" style="background: #E0F2F1;">📝</view>
          <text class="label">积分派发记录</text>
        </view>
        <text class="arrow">➔</text>
      </view>

      <view class="card menu-item" v-if="userInfo.role === 'parent'" @click="showBindCode">
        <view class="menu-left">
          <view class="icon-box" style="background: #E1F5FE;">🔗</view>
          <text class="label">获取邀请绑定码</text>
        </view>
        <text class="arrow">➔</text>
      </view>

      <view class="card menu-item" v-if="userInfo.role === 'child'" @click="showBindInput">
        <view class="menu-left">
          <view class="icon-box" style="background: #F3E5F5;">➕</view>
          <text class="label">绑定我的家长</text>
        </view>
        <text class="arrow">➔</text>
      </view>
    </view>

    <!-- Logout -->
    <view class="action-area">
      <button class="logout-btn" @click="handleLogout">退出登录 🚪</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { request } from '../../utils/request';

const userInfo = ref<any>({});
const avatarRefreshKey = ref(0);

const avatarUrl = computed(() => {
  avatarRefreshKey.value; // dependency
  if (!userInfo.value.avatar_url) return '/static/logo.png';
  // Add timestamp to force refresh
  const url = userInfo.value.avatar_url;
  const separator = url.includes('?') ? '&' : '?';
  return url + separator + 't=' + Date.now();
});
const isLoggedIn = ref(false);

const loadData = () => {
  const token = uni.getStorageSync('token');
  if (token) {
    isLoggedIn.value = true;
    const storedUser = uni.getStorageSync('userInfo');
    if (storedUser) {
      userInfo.value = JSON.parse(storedUser);
    }
  }
};

onMounted(loadData);

const changeNickname = () => {
  uni.showModal({
    title: '修改昵称',
    editable: true,
    placeholderText: '请输入新昵称',
    success: async (res) => {
      if (res.confirm && res.content) {
        try {
          const role = userInfo.value.role;
          const url = role === 'parent' ? '/parent/profile' : '/child/profile';
          await request({
            url: url,
            method: 'POST',
            data: { nickname: res.content }
          });
          userInfo.value.nickname = res.content;
          uni.setStorageSync('userInfo', JSON.stringify(userInfo.value));
          uni.showToast({ title: '修改成功', icon: 'success' });
        } catch (e) {
          uni.showToast({ title: '修改失败', icon: 'none' });
        }
      }
    }
  });
};

const changeAvatar = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempFilePath = res.tempFilePaths[0];
      uni.showLoading({ title: '上传中...' });

      try {
        const token = uni.getStorageSync('token');
        const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
        const url = userInfo.role === 'parent' ? '/parent/avatar' : '/child/avatar';

        // Upload file
        const uploadRes = await new Promise((resolve, reject) => {
          uni.uploadFile({
            url: 'http://localhost:18081' + url,
            filePath: tempFilePath,
            name: 'file',
            header: {
              'Authorization': 'Bearer ' + token
            },
            success: (res) => {
              if (res.statusCode === 200) {
                resolve(JSON.parse(res.data));
              } else {
                reject(res);
              }
            },
            fail: reject
          });
        });

        // Update local user info
        const newUserInfo = { ...userInfo, avatar_url: uploadRes.url };
        uni.setStorageSync('userInfo', JSON.stringify(newUserInfo));

        // Update global ref
        userInfo.value = newUserInfo;

        // Force refresh avatar
        avatarRefreshKey.value++;

        uni.hideLoading();
        uni.showToast({ title: '修改成功', icon: 'success' });
      } catch (e) {
        uni.hideLoading();
        uni.showToast({ title: '上传失败', icon: 'none' });
        console.error(e);
      }
    }
  });
};

const handleLogout = () => {
  uni.showModal({
    title: '退出确认',
    content: '确定要离开吗？',
    cancelText: '点错了',
    confirmText: '确定',
    confirmColor: '#FF6B35',
    success: (res) => {
      if (res.confirm) {
        uni.clearStorageSync();
        uni.reLaunch({ url: '/pages/index/index' });
      }
    }
  });
};

const goToRecords = () => uni.navigateTo({ url: '/pages/records/records' });

const goToPointsRecords = () => uni.navigateTo({ url: '/pages/points-records/points-records' });

const showBindCode = async () => {
  try {
    const res = await request({ url: '/parent/binding/code', method: 'POST' });
    uni.showModal({
      title: '专属绑定码',
      content: res.bind_code,
      showCancel: false,
      confirmText: '复制去发送'
    });
  } catch (e) {
    uni.showToast({ title: '获取失败', icon: 'none' });
  }
};

const showBindInput = () => {
  uni.showModal({
    title: '输入绑定码',
    editable: true,
    placeholderText: '输入家长提供的6位代码',
    success: async (res) => {
      if (res.confirm && res.content) {
        try {
          await request({
            url: '/child/binding/accept',
            method: 'POST',
            data: { bind_code: res.content }
          });
          uni.showToast({ title: '绑定成功！', icon: 'success' });
        } catch (e) {
          uni.showToast({ title: '验证失败', icon: 'none' });
        }
      }
    }
  });
};
</script>

<style lang="scss" scoped>
.container {
  padding-bottom: 40rpx;
}

.hero-card {
  margin: 30rpx 32rpx 40rpx;
  border-radius: 48rpx;
  padding: 40rpx;
  color: white;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #FF6B35, #FF8FAB, #C77DFF);
  box-shadow: 0 20rpx 50rpx rgba(255, 107, 53, 0.3);

  .bg-circle {
    position: absolute;
    background: rgba(255, 255, 255, 0.15);
    border-radius: 50%;
    &.top-right { width: 200rpx; height: 200rpx; top: -40rpx; right: -40rpx; }
    &.bottom-left { width: 120rpx; height: 120rpx; bottom: 20rpx; left: -20rpx; opacity: 0.1; }
  }

  .user-header {
    display: flex;
    align-items: center;
    margin-bottom: 40rpx;
    position: relative;
    z-index: 1;

    .avatar {
      width: 100rpx;
      height: 100rpx;
      border-radius: 50%;
      border: 4rpx solid rgba(255,255,255,0.4);
      margin-right: 20rpx;
    }

    .nickname {
      font-size: 36rpx;
      display: block;
    }

    .role-badge {
      font-size: 20rpx;
      background: rgba(255,255,255,0.2);
      padding: 4rpx 16rpx;
      border-radius: 20rpx;
      margin-top: 4rpx;
    }
  }

  .points-display {
    position: relative;
    z-index: 1;
    .label { font-size: 26rpx; opacity: 0.9; font-weight: 800; }
    .amount-row {
      display: flex;
      align-items: baseline;
      gap: 16rpx;
      margin: 10rpx 0;
      .amount { font-size: 100rpx; line-height: 1; }
      .unit { font-size: 32rpx; opacity: 0.9; }
    }
  }

  .moto { font-size: 28rpx; font-weight: 700; opacity: 0.9; position: relative; z-index: 1; }
  .decoration { position: absolute; right: 40rpx; bottom: 40rpx; font-size: 96rpx; opacity: 0.8; }
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 30rpx;
  margin-bottom: 24rpx;

  .menu-left {
    display: flex;
    align-items: center;
    .icon-box {
      width: 80rpx;
      height: 80rpx;
      border-radius: 24rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 40rpx;
      margin-right: 24rpx;
    }
    .label { font-size: 30rpx; font-weight: 800; color: #2D2D2D; }
  }
  .arrow { color: #DDD; font-weight: 900; }
}

.action-area {
  padding: 40rpx 32rpx;
  .logout-btn {
    background: linear-gradient(135deg, #FF6B35, #FF8FAB, #C77DFF);
    color: white;
    border-radius: 24rpx;
    font-weight: 800;
    font-size: 28rpx;
    &::after { border: none; }
  }
}
</style>
