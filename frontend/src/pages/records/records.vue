<template>
  <view class="container">
    <scroll-view
      scroll-y
      class="list-scroll"
      @scrolltolower="loadMore"
      refresh-with-animation
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
    >
      <view v-if="records.length > 0" class="record-list">
        <view v-for="item in records" :key="item.id" class="record-card">
          <image class="item-image" :src="item.item?.image_path || '/static/logo.png'" mode="aspectFill" />
          <view class="item-info">
            <view class="header">
              <text class="name">{{ item.item?.name || '未知商品' }}</text>
              <text :class="['status', item.status]">{{ formatStatus(item.status) }}</text>
            </view>
            <view class="body">
              <text class="user">兑换人: {{ item.user?.nickname || '匿名孩子' }}</text>
              <text class="time">{{ formatTime(item.created_at) }}</text>
            </view>
            <!-- Action buttons for Parent -->
            <view v-if="userRole === 'parent' && item.status === 'pending'" class="actions">
              <button class="action-btn confirm" @click="confirmRedemption(item.id)">确认兑现</button>
              <button class="action-btn cancel" @click="cancelRedemption(item.id)">取消</button>
            </view>
          </view>
        </view>
        <view class="loading-more" v-if="loading">加载中...</view>
        <view class="no-more" v-if="!hasMore && records.length > 0">没有更多了</view>
      </view>

      <view v-else-if="!loading" class="empty-state">
        <text class="icon">🎁</text>
        <text class="text">暂无兑换记录</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { request } from '../../utils/request';

interface Redemption {
  id: number;
  status: string;
  created_at: string;
  item?: {
    name: string;
    image_path: string;
  };
  user?: {
    nickname: string;
  };
}

const records = ref<Redemption[]>([]);
const userRole = ref('');
const loading = ref(false);
const isRefreshing = ref(false);
const hasMore = ref(true);
const page = ref(1);

const fetchRecords = async (isRefresh = false) => {
  if (loading.value) return;
  loading.value = true;

  try {
    const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
    userRole.value = userInfo.role;
    const endpoint = userInfo.role === 'parent' ? '/parent/redemptions' : '/child/redemptions';

    const res = await request<Redemption[]>({
      url: endpoint,
      method: 'GET'
    });

    if (isRefresh) {
      records.value = res;
      isRefreshing.value = false;
    } else {
      // For demo simplicity, we assume single page if backend doesn't support pagination yet
      records.value = [...records.value, ...res];
      hasMore.value = false;
    }
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' });
  } finally {
    loading.value = false;
  }
};

const onRefresh = () => {
  isRefreshing.value = true;
  page.value = 1;
  hasMore.value = true;
  fetchRecords(true);
};

const loadMore = () => {
  if (hasMore.value && !loading.value) {
    page.value++;
    fetchRecords();
  }
};

const formatStatus = (status: string) => {
  const map: Record<string, string> = {
    'pending': '待兑现',
    'completed': '已兑现',
    'confirmed': '已兑现',
    'cancelled': '已取消'
  };
  return map[status] || status;
};

const confirmRedemption = async (id: number) => {
  uni.showModal({
    title: '确认兑现',
    content: '确定要标记为已兑现吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({
            url: `/parent/redemption/confirm/${id}`,
            method: 'POST'
          });
          onRefresh();
          uni.showToast({ title: '已兑现', icon: 'success' });
        } catch (e) {
          console.error(e);
          uni.showToast({ title: '操作失败', icon: 'none' });
        }
      }
    }
  });
};

const cancelRedemption = async (id: number) => {
  uni.showModal({
    title: '取消兑换',
    content: '确定要取消该兑换吗？积分将自动返还给孩子。',
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({
            url: `/parent/redemption/cancel/${id}`,
            method: 'POST'
          });
          onRefresh();
          uni.showToast({ title: '已取消', icon: 'success' });
        } catch (e) {
          console.error(e);
          uni.showToast({ title: '操作失败', icon: 'none' });
        }
      }
    }
  });
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
}

.list-scroll {
  height: 100%;
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

  .item-image {
    width: 140rpx;
    height: 140rpx;
    border-radius: 8rpx;
    margin-right: 20rpx;
    background-color: #f0f0f0;
  }

  .item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;

      .name {
        font-size: 32rpx;
        font-weight: bold;
        color: #333;
      }

      .status {
        font-size: 24rpx;
        padding: 4rpx 12rpx;
        border-radius: 6rpx;

        &.pending { background-color: #fff7e6; color: #fa8c16; }
        &.completed { background-color: #f6ffed; color: #52c41a; }
        &.confirmed { background-color: #f6ffed; color: #52c41a; }
        &.cancelled { background-color: #fff1f0; color: #f5222d; }
      }
    }

    .body {
      .user {
        font-size: 26rpx;
        color: #666;
        display: block;
      }
      .time {
        font-size: 24rpx;
        color: #999;
        margin-top: 4rpx;
      }
    }

    .actions {
      display: flex;
      justify-content: flex-end;
      margin-top: 10rpx;
      .action-btn {
        margin: 0;
        padding: 0 30rpx;
        height: 54rpx;
        line-height: 54rpx;
        font-size: 24rpx;
        border-radius: 27rpx;
        &::after { border: none; }
        &.cancel {
          background-color: #ffffff;
          color: #999;
          border: 1px solid #ddd;
          margin-left: 12rpx;
        }
        &.confirm {
          background-color: #007AFF;
          color: #ffffff;
        }
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

.loading-more, .no-more {
  text-align: center;
  padding: 20rpx;
  font-size: 24rpx;
  color: #999;
}
</style>
