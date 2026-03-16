<template>
  <view class="container">
    <view class="header card">
      <text class="title">积分派发记录</text>
      <text class="subtitle">查看给孩子的积分调整记录</text>
    </view>

    <scroll-view
      scroll-y
      class="list-scroll"
    >
      <view v-if="records.length > 0" class="record-list">
        <view v-for="item in records" :key="item.id" class="record-card">
          <view class="record-info">
            <view class="header">
              <text class="child-name">{{ item.user?.nickname || '孩子' }}</text>
              <text :class="['amount', item.amount >= 0 ? 'plus' : 'minus']">
                {{ item.amount >= 0 ? '+' : '' }}{{ item.amount }}
              </text>
            </view>
            <view class="body">
              <text class="reason">原因: {{ item.reason }}</text>
              <text class="time">{{ formatTime(item.created_at) }}</text>
            </view>
          </view>
        </view>
      </view>

      <view v-else class="empty-state">
        <text class="icon">📝</text>
        <text class="text">暂无积分派发记录</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { request } from '../../utils/request';

interface PointsRecord {
  id: number;
  amount: number;
  reason: string;
  created_at: string;
  user?: {
    nickname: string;
    name: string;
  };
}

const records = ref<PointsRecord[]>([]);

const fetchRecords = async () => {
  try {
    const res = await request<PointsRecord[]>({
      url: '/parent/points/records',
      method: 'GET'
    });
    records.value = res || [];
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
};

const formatTime = (timeStr: string) => {
  if (!timeStr) return '';
  const date = new Date(timeStr);
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`;
};

onMounted(() => {
  fetchRecords();
});
</script>

<style lang="scss" scoped>
.container {
  height: 100vh;
  background-color: #f8f8f8;
  display: flex;
  flex-direction: column;
}

.header {
  padding: 40rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
  color: white;
  .title { font-size: 40rpx; font-weight: bold; display: block; }
  .subtitle { font-size: 24rpx; opacity: 0.9; margin-top: 8rpx; display: block; }
}

.list-scroll {
  flex: 1;
  height: 0; // Important for flex child
}

.record-list {
  padding: 20rpx;
}

.record-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  display: flex;
  box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.05);

  .record-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      background: transparent;
      padding: 0;
      margin-bottom: 0;

      .child-name {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
      }

      .amount {
        font-size: 36rpx;
        font-weight: bold;
        &.plus { color: #FF6B35; }
        &.minus { color: #666; }
      }
    }

    .body {
      display: flex;
      justify-content: space-between;
      margin-top: 12rpx;
      .reason {
        font-size: 26rpx;
        color: #666;
      }
      .time {
        font-size: 24rpx;
        color: #999;
      }
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 200rpx;

  .icon { font-size: 100rpx; margin-bottom: 20rpx; }
  .text { font-size: 28rpx; color: #999; }
}

.card {
  background: white;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}
</style>
