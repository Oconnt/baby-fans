<template>
  <view class="container">
    <view class="header">
      <text class="title">任务管理</text>
    </view>

    <!-- Task List -->
    <view class="task-list">
      <view v-for="task in tasks" :key="task.id" class="task-card card">
        <view class="task-header">
          <text class="task-name">{{ task.name }}</text>
          <text class="task-status" :class="getStatusClass(task.status)">{{ getStatusText(task.status) }}</text>
        </view>
        <text class="task-desc">{{ task.description || '无描述' }}</text>
        <view class="task-footer">
          <text class="task-handler">{{ task.handler ? (task.handler.nickname || task.handler.name) : '' }}</text>
          <text class="task-points">+{{ task.points }}分</text>
        </view>
        <view class="task-meta">
          <text>发布: {{ formatTime(task.publish_time) }}</text>
          <text>过期: {{ formatTime(task.expire_time) }}</text>
        </view>
        <view v-if="task.status === 1" class="task-actions">
          <view class="btn-action" @click="markExpired(task)">标记过期</view>
        </view>
      </view>
      <view v-if="tasks.length === 0" class="empty-state">
        <text>暂无任务</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import { request } from '../../utils/request';

const tasks = ref<any[]>([]);

onMounted(() => {
  const token = uni.getStorageSync('token');
  if (!token) {
    uni.reLaunch({ url: '/pages/login/login' });
    return;
  }
  fetchTasks();
});

onShow(() => {
  fetchTasks();
});

const fetchTasks = async () => {
  try {
    const res = await request({ url: '/parent/tasks', method: 'GET' });
    tasks.value = res || [];
  } catch (e) {
    console.error(e);
    tasks.value = [];
  }
};

const getStatusText = (status: number) => {
  switch (status) {
    case 1: return '待办';
    case 2: return '已完成';
    case 3: return '已过期';
    default: return '未知';
  }
};

const getStatusClass = (status: number) => {
  switch (status) {
    case 1: return 'pending';
    case 2: return 'completed';
    case 3: return 'expired';
    default: return '';
  }
};

const formatTime = (timeStr: string) => {
  if (!timeStr) return '';
  const date = new Date(timeStr);
  return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`;
};

const markExpired = async (task: any) => {
  uni.showModal({
    title: '确认',
    content: '确定将此任务标记为过期?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({
            url: `/parent/tasks/${task.id}/status`,
            method: 'PUT',
            data: { status: 3 }
          });
          uni.showToast({ title: '已标记过期', icon: 'success' });
          fetchTasks();
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' });
        }
      }
    }
  });
};
</script>

<style lang="scss" scoped>
.container {
  padding: 30rpx;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.header {
  padding: 30rpx;
  margin-bottom: 30rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
  border-radius: 24rpx;
  .title { font-size: 36rpx; font-weight: bold; color: white; }
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.task-card {
  padding: 30rpx;
  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12rpx;
    .task-name { font-size: 32rpx; font-weight: bold; color: #333; }
    .task-status {
      font-size: 24rpx;
      padding: 4rpx 16rpx;
      border-radius: 20rpx;
      &.pending { background: #E6F2FF; color: #007AFF; }
      &.completed { background: #E8F8E8; color: #4CD964; }
      &.expired { background: #FFE8E8; color: #FF3B30; }
    }
  }
  .task-desc { font-size: 26rpx; color: #666; margin-bottom: 16rpx; display: block; }
  .task-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8rpx;
    .task-handler { font-size: 26rpx; color: #666; }
    .task-points { font-size: 28rpx; font-weight: bold; color: #FF6B35; }
  }
  .task-meta {
    display: flex;
    justify-content: space-between;
    font-size: 22rpx;
    color: #999;
  }
  .task-actions {
    margin-top: 16rpx;
    padding-top: 16rpx;
    border-top: 1px solid #f0f0f0;
    .btn-action {
      background: #f0f0f0; color: #666; border-radius: 20rpx;
      height: 56rpx; line-height: 56rpx; text-align: center;
      font-size: 24rpx;
    }
  }
}

.empty-state {
  text-align: center;
  color: #999;
  font-size: 28rpx;
  padding: 80rpx;
}

.card {
  background: white;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}
</style>
