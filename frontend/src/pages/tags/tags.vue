<template>
  <view class="container">

    <!-- Parent View: Resource Management Hub -->
    <view v-if="userRole === 'parent'">
      <view class="hub-header">
        <text class="hub-title">资源管理</text>
        <text class="hub-subtitle">管理标签、任务模版和任务</text>
      </view>

      <!-- Navigation Cards -->
      <view class="nav-cards">
        <view class="nav-card card" @click="goToPage('/pages/label-list/label-list')">
          <view class="nav-icon tag-icon">🏷️</view>
          <text class="nav-title">标签管理</text>
        </view>

        <view class="nav-card card" @click="goToPage('/pages/task-template-list/task-template-list')">
          <view class="nav-icon task-icon">📋</view>
          <text class="nav-title">任务模版</text>
        </view>

        <view class="nav-card card" @click="goToPage('/pages/task-list/task-list')">
          <view class="nav-icon list-icon">📝</view>
          <text class="nav-title">任务管理</text>
        </view>
      </view>
    </view>

    <!-- Child View: Points History -->
    <view v-if="userRole === 'child'" class="history-view">
      <view v-for="item in records" :key="item.id" class="record-item card">
        <view class="record-info">
          <text class="reason">{{ item.reason }}</text>
          <text class="operator">{{ item.operator ? (item.operator.nickname || item.operator.name) : '' }}</text>
          <text class="time">{{ formatTime(item.created_at) }}</text>
        </view>
        <text class="amount" :class="item.amount >= 0 ? 'plus' : 'minus'">
          {{ item.amount >= 0 ? '+' : '' }}{{ item.amount }}
        </text>
      </view>
      <view v-if="records.length === 0" class="empty-state">
        <text>暂无记录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import { request } from '../../utils/request';

const userRole = ref('');
const records = ref<any[]>([]);

const updateRoleAndTitle = () => {
  const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
  userRole.value = userInfo.role;

  if (userRole.value === 'parent') {
    uni.setNavigationBarTitle({ title: '资源管理' });
  } else {
    uni.setNavigationBarTitle({ title: '积分修改记录' });
    fetchHistory();
  }
};

onMounted(() => {
  const token = uni.getStorageSync('token');
  if (!token) {
    uni.reLaunch({ url: '/pages/login/login' });
    return;
  }
  updateRoleAndTitle();
});

onShow(() => {
  updateRoleAndTitle();
});

const onTabItemTap = () => {
  updateRoleAndTitle();
};

const fetchHistory = async () => {
  try {
    const res = await request({ url: '/child/points/history', method: 'GET' });
    records.value = res || [];
  } catch (e) {
    console.error(e);
  }
};

const formatTime = (timeStr: string) => {
  if (!timeStr) return '';
  const date = new Date(timeStr);
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`;
};

const goToPage = (url: string) => {
  // Only tabbar pages can use switchTab
  const tabbarPages = ['/pages/home-parent/home-parent', '/pages/tags/tags', '/pages/shop/shop', '/pages/mine/mine'];
  if (tabbarPages.includes(url)) {
    uni.switchTab({ url });
  } else {
    uni.navigateTo({ url });
  }
};
</script>

<style lang="scss" scoped>
.container {
  padding: 30rpx;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.hub-header {
  padding: 40rpx;
  margin-bottom: 30rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
  border-radius: 24rpx;
  .hub-title { font-size: 40rpx; font-weight: bold; color: white; display: block; }
  .hub-subtitle { font-size: 24rpx; color: rgba(255,255,255,0.9); margin-top: 8rpx; display: block; }
}

.nav-cards {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.nav-card {
  padding: 40rpx;
  display: flex;
  align-items: center;
  .nav-icon {
    font-size: 60rpx;
    margin-right: 30rpx;
  }
  .nav-title {
    font-size: 32rpx; font-weight: bold; color: #333; display: block;
  }
  .nav-desc {
    font-size: 24rpx; color: #999; margin-top: 8rpx; display: block;
  }
}

.card {
  background: white;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}

.history-view {
  display: flex; flex-direction: column; gap: 20rpx;
  .record-item {
    padding: 24rpx;
    display: flex; justify-content: space-between; align-items: center;
    .record-info {
      flex: 1;
      .reason { font-size: 28rpx; color: #333; display: block; font-weight: bold; }
      .operator { font-size: 24rpx; color: #999; display: block; margin-top: 4rpx; }
      .time { font-size: 22rpx; color: #bbb; display: block; margin-top: 4rpx; }
    }
    .amount { font-size: 32rpx; font-weight: bold; }
    .plus { color: #FF6B35; }
    .minus { color: #666; }
  }
}

.empty-state {
  text-align: center;
  color: #999;
  font-size: 28rpx;
  padding: 80rpx;
}
</style>
