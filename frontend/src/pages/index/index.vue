<template>
  <view class="container">
    <!-- Header Section -->
    <view class="header card">
      <view class="user-info">
        <text class="greeting">你好, {{ userInfo.nickname }}! 👋</text>
        <text class="role-tag" :class="userInfo.role">{{ userInfo.role === 'parent' ? '家长模式' : '宝宝模式' }}</text>
      </view>
      <view class="stats" v-if="userInfo.role === 'child'">
        <text class="points-label">我的积分</text>
        <text class="points-value">{{ userInfo.points }}</text>
      </view>
    </view>

    <!-- Parent: Children Management Section -->
    <view v-if="userInfo.role === 'parent'" class="section">
      <view class="section-header">
        <text class="section-title">我的孩子</text>
        <button class="bind-btn" @click="generateBindCode">绑定新孩子 +</button>
      </view>

      <view v-if="childrenList.length > 0" class="children-list">
        <view v-for="child in childrenList" :key="child.id" class="child-card card">
          <view class="child-info">
            <text class="child-name">{{ child.nickname || child.name }}</text>
            <text class="child-points">当前积分: {{ child.points }} ⭐</text>
          </view>
          <view class="child-actions">
            <button class="action-btn primary" @click="adjustPoints(child)">管理积分</button>
            <button class="action-btn" @click="viewRecords(child)">变动记录</button>
          </view>
        </view>
      </view>
      <view v-else class="empty-state card">
        <text>还没绑定孩子，点击上方按钮开始吧</text>
      </view>
    </view>

    <!-- Navigation Grids -->
    <view class="grid-container">
      <view v-for="(tag, index) in activeTags" :key="index" class="grid-item card" @click="handleTagClick(tag)">
        <text class="icon">{{ tag.icon }}</text>
        <text class="label">{{ tag.label }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { request } from '../../utils/request';

const userInfo = ref({
  id: 0,
  nickname: '访客',
  role: 'child',
  points: 0
});

const childrenList = ref<any[]>([]);

onMounted(() => {
  const stored = uni.getStorageSync('userInfo');
  if (stored) {
    userInfo.value = JSON.parse(stored);
    if (userInfo.value.role === 'parent') {
      fetchChildren();
    }
  } else {
    uni.reLaunch({ url: '/pages/login/login' });
  }
});

const fetchChildren = async () => {
  try {
    const res = await request({ url: '/parent/children', method: 'GET' });
    childrenList.value = res;
  } catch (e) {
    console.error('Failed to fetch children', e);
  }
};

const generateBindCode = async () => {
  try {
    const res = await request({ url: '/parent/binding/code', method: 'POST' });
    uni.showModal({
      title: '专属绑定码',
      content: res.bind_code,
      showCancel: false,
      confirmText: '复制并发送'
    });
  } catch (e) {
    uni.showToast({ title: '获取失败', icon: 'none' });
  }
};

const adjustPoints = (child: any) => {
  uni.showModal({
    title: `为 ${child.nickname || child.name} 调整积分`,
    editable: true,
    placeholderText: '输入变动数值 (如: 10 或 -5)',
    success: async (res) => {
      if (res.confirm && res.content) {
        const amount = parseInt(res.content);
        if (isNaN(amount)) return uni.showToast({ title: '请输入有效数字', icon: 'none' });

        try {
          await request({
            url: '/parent/points/manage',
            method: 'POST',
            data: {
              user_id: child.id,
              amount: amount,
              reason: '手动调整'
            }
          });
          uni.showToast({ title: '操作成功', icon: 'success' });
          fetchChildren();
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' });
        }
      }
    }
  });
};

const viewRecords = (child: any) => {
  uni.navigateTo({ url: `/pages/records/records?userId=${child.id}` });
};

const parentTags = [
  { label: '刷新列表', icon: '🔄', action: () => fetchChildren() },
  { label: '积分配置', icon: '⚙️', action: () => uni.navigateTo({ url: '/pages/points/points' }) },
  { label: '前往商城', icon: '🛍️', action: () => uni.switchTab({ url: '/pages/shop/shop' }) },
  { label: '个人中心', icon: '👤', action: () => uni.switchTab({ url: '/pages/mine/mine' }) }
];

const childTags = [
  { label: '变动记录', icon: '📜', action: () => uni.navigateTo({ url: '/pages/records/records' }) },
  { label: '心愿商城', icon: '🎁', action: () => uni.switchTab({ url: '/pages/shop/shop' }) },
  { label: '我的资料', icon: '⭐', action: () => uni.switchTab({ url: '/pages/mine/mine' }) },
  { label: '退出登录', icon: '🚪', action: () => {
    uni.clearStorageSync();
    uni.reLaunch({ url: '/pages/login/login' });
  }}
];

const activeTags = computed(() => userInfo.value.role === 'parent' ? parentTags : childTags);

const handleTagClick = (tag: any) => {
  if (tag.action) tag.action();
};
</script>

<style lang="scss" scoped>
.container {
  padding: 30rpx;
  background-color: #f8f9fa;
  min-height: 100vh;
}

.header {
  padding: 40rpx;
  margin-bottom: 40rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FF8FAB 100%);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .greeting {
    font-size: 36rpx;
    font-weight: bold;
    display: block;
  }

  .role-tag {
    font-size: 20rpx;
    padding: 4rpx 16rpx;
    border-radius: 20rpx;
    background: rgba(255,255,255,0.2);
    margin-top: 10rpx;
    display: inline-block;
    &.parent { border: 1px solid #C77DFF; color: #fdfdfd; }
  }

  .points-value {
    font-size: 48rpx;
    font-weight: 800;
    display: block;
    text-align: right;
  }
  .points-label {
    font-size: 20rpx;
    opacity: 0.8;
  }
}

.section {
  margin-bottom: 40rpx;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20rpx;

    .section-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #333;
    }

    .bind-btn {
      font-size: 24rpx;
      background-color: #FF6B35;
      color: white;
      border-radius: 30rpx;
      padding: 0 30rpx;
      height: 60rpx;
      line-height: 60rpx;
      margin: 0;
    }
  }
}

.children-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.child-card {
  padding: 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .child-info {
    .child-name {
      font-size: 30rpx;
      font-weight: bold;
      color: #333;
      display: block;
    }
    .child-points {
      font-size: 24rpx;
      color: #FF6B35;
      margin-top: 10rpx;
      display: block;
    }
  }

  .child-actions {
    display: flex;
    gap: 10rpx;

    .action-btn {
      font-size: 22rpx;
      padding: 0 20rpx;
      height: 50rpx;
      line-height: 50rpx;
      border-radius: 10rpx;
      background-color: #f0f0f0;
      color: #666;
      margin: 0;
      &::after { border: none; }
      &.primary {
        background-color: #FF6B35;
        color: white;
      }
    }
  }
}

.empty-state {
  padding: 60rpx;
  text-align: center;
  color: #999;
  font-size: 26rpx;
}

.grid-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30rpx;
}

.grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60rpx 0;
  transition: transform 0.2s;

  &:active {
    transform: scale(0.95);
  }

  .icon {
    font-size: 80rpx;
    margin-bottom: 20rpx;
  }

  .label {
    font-size: 28rpx;
    font-weight: 600;
    color: #333;
  }
}

.card {
  background: white;
  border-radius: 32rpx;
  box-shadow: 0 10rpx 30rpx rgba(0,0,0,0.05);
}
</style>
