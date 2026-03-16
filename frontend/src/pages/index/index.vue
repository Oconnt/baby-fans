<template>
  <view class="container">
    <!-- Header -->
    <view class="header card">
      <text class="title">孩子管理</text>
      <text class="subtitle">管理已绑定的孩子</text>
    </view>

    <!-- Usage Tip -->
    <view class="usage-tip card">
      <text class="tip-text">💡 使用方法：点击【管理】按钮可手动调整积分，结合【标签管理】中定义的标签可实现快速发放。</text>
    </view>

    <!-- Bind by Login Code -->
    <view class="bind-section">
      <view class="bind-row">
        <input class="bind-input" v-model="loginCode" placeholder="输入孩子登录码" />
        <button class="bind-btn" @click="bindChild">绑定</button>
      </view>
    </view>

    <!-- Children List -->
    <view v-if="childrenList.length > 0" class="children-list">
      <view v-for="child in childrenList" :key="child.id" class="child-card card">
        <view class="child-info">
          <text class="child-name">{{ child.nickname || child.name }}</text>
          <text class="child-points">积分: {{ child.points }} ⭐</text>
        </view>
        <view class="child-actions">
          <button class="action-btn manage" @click="openAdjustModal(child)">管理</button>
          <button class="action-btn unbind" @click="unbindChild(child)">解绑</button>
        </view>
      </view>
    </view>
    <view v-else class="empty-state card">
      <text>暂无绑定的孩子，请输入登录码绑定</text>
    </view>

    <!-- Points Adjust Modal -->
    <view v-if="showPointsModal" class="modal-mask" @click="showPointsModal = false">
      <view class="modal-content card" @click.stop>
        <text class="modal-title">调整积分 - {{ selectedChild?.nickname || selectedChild?.name }}</text>

        <view v-if="pointTemplates.length > 0" class="template-section">
          <text class="section-label">快捷标签</text>
          <view class="template-grid">
            <view
              v-for="item in pointTemplates"
              :key="item.id"
              class="template-chip"
              :class="item.amount >= 0 ? 'plus' : 'minus'"
              @click="handleTemplateClick(item)"
            >
              <text class="chip-title">{{ item.title }}</text>
              <text class="chip-amount">{{ item.amount >= 0 ? '+' : '' }}{{ item.amount }}</text>
            </view>
          </view>
        </view>
        <view v-else class="empty-templates">
          <text>暂无标签，请前往标签管理添加</text>
        </view>

        <view class="divider"></view>

        <view class="manual-section">
          <text class="section-label">手动输入</text>
          <input class="input-box" v-model="manualAmount" type="number" placeholder="输入积分值 (如: 10 或 -5)" />
          <view class="btn-submit" @click="handleManualSubmit">确定发放</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { request } from '../../utils/request';

    const childrenList = ref<any[]>([]);
    const loginCode = ref('');
    const showPointsModal = ref(false);
    const selectedChild = ref<any>(null);
    const pointTemplates = ref<any[]>([]);
    const manualAmount = ref('');

    onMounted(() => {
      const stored = uni.getStorageSync('userInfo');
      if (!stored) {
        uni.reLaunch({ url: '/pages/login/login' });
        return;
      }
      fetchChildren();
    });

    const fetchChildren = async () => {
      try {
        const res = await request({ url: '/parent/children', method: 'GET' });
        childrenList.value = res || [];
      } catch (e) {
        childrenList.value = [];
      }
    };

    const fetchTemplates = async () => {
      try {
        const res = await request({ url: '/parent/templates', method: 'GET' });
        pointTemplates.value = res || [];
      } catch (e) {
        pointTemplates.value = [];
      }
    };

    const openAdjustModal = async (child: any) => {
      selectedChild.value = child;
      manualAmount.value = '';
      await fetchTemplates();
      showPointsModal.value = true;
    };

    const submitPoints = async (amount: number, reason: string) => {
      if (!selectedChild.value) return;
      try {
        await request({
          url: '/parent/points/manage',
          method: 'POST',
          data: { user_id: selectedChild.value.id, amount, reason }
        });
        uni.showToast({ title: '操作成功', icon: 'success' });
        showPointsModal.value = false;
        fetchChildren();
      } catch (e) {
        uni.showToast({ title: '操作失败', icon: 'none' });
      }
    };

    const handleTemplateClick = (template: any) => {
      submitPoints(template.amount, template.title);
    };

    const handleManualSubmit = () => {
      const amount = parseInt(manualAmount.value);
      if (isNaN(amount)) return uni.showToast({ title: '请输入有效数字', icon: 'none' });
      submitPoints(amount, '手动调整');
    };

const bindChild = async () => {
  const code = loginCode.value.trim();
  if (!code) {
    return uni.showToast({ title: '请输入登录码', icon: 'none' });
  }
  try {
    await request({
      url: '/parent/children/bind',
      method: 'POST',
      data: { login_code: code }
    });
    uni.showToast({ title: '绑定成功', icon: 'success' });
    loginCode.value = '';
    fetchChildren();
  } catch (e: any) {
    uni.showToast({ title: e?.error || '绑定失败', icon: 'none' });
  }
};

const adjustPoints = (child: any) => {
  uni.showModal({
    title: `为 ${child.nickname || child.name} 调整积分`,
    editable: true,
    placeholderText: '输入数值 (如: 10 或 -5)',
    success: async (res) => {
      if (res.confirm && res.content) {
        const amount = parseInt(res.content);
        if (isNaN(amount)) return uni.showToast({ title: '请输入有效数字', icon: 'none' });
        try {
          await request({
            url: '/parent/points/manage',
            method: 'POST',
            data: { user_id: child.id, amount, reason: '手动调整' }
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

const unbindChild = (child: any) => {
  uni.showModal({
    title: '确认解绑',
    content: `确定要解绑 ${child.nickname || child.name} 吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({ url: `/parent/children/${child.id}`, method: 'DELETE' });
          uni.showToast({ title: '已解绑', icon: 'success' });
          fetchChildren();
        } catch (e) {
          uni.showToast({ title: '解绑失败', icon: 'none' });
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
  padding: 40rpx;
  margin-bottom: 30rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FF8FAB 100%);
  color: white;
  .title { font-size: 40rpx; font-weight: bold; display: block; }
  .subtitle { font-size: 24rpx; opacity: 0.9; margin-top: 8rpx; display: block; }
}

.usage-tip {
  padding: 24rpx;
  margin-bottom: 30rpx;
  background: #f0f7ff;
  border: 1px solid #e0eeff;
  .tip-text { font-size: 24rpx; color: #007AFF; line-height: 1.6; }
}

.bind-section {
  margin-bottom: 30rpx;
  .bind-row {
    display: flex;
    gap: 16rpx;
  }
  .bind-input {
    flex: 1;
    background: white;
    border: 1px solid #eee;
    border-radius: 20rpx;
    padding: 0 24rpx;
    height: 80rpx;
    font-size: 28rpx;
  }
  .bind-btn {
    background: #FF6B35;
    color: white;
    border-radius: 20rpx;
    font-size: 28rpx;
    font-weight: bold;
    padding: 0 40rpx;
    height: 80rpx;
    line-height: 80rpx;
    margin: 0;
    &::after { border: none; }
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
    .child-name { font-size: 32rpx; font-weight: bold; color: #333; display: block; }
    .child-points { font-size: 24rpx; color: #FF6B35; margin-top: 8rpx; display: block; }
  }
  .child-actions {
    display: flex;
    gap: 12rpx;
    .action-btn {
      font-size: 24rpx;
      padding: 0 24rpx;
      height: 56rpx;
      line-height: 56rpx;
      border-radius: 12rpx;
      margin: 0;
      &::after { border: none; }
      &.manage { background: #FF6B35; color: white; }
      &.unbind { background: #f0f0f0; color: #999; }
    }
  }
}

.empty-state {
  padding: 80rpx 40rpx;
  text-align: center;
  color: #999;
  font-size: 28rpx;
}

.card {
  background: white;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}

.modal-mask {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  z-index: 9999; display: flex;
  align-items: center; justify-content: center;
  background: rgba(0,0,0,0.6);
}

.modal-content {
  width: 85%; max-height: 80vh; padding: 40rpx;
  background: #ffffff; border-radius: 24rpx;
  pointer-events: auto;
  display: flex; flex-direction: column;
  overflow-y: auto;

  .modal-title {
    font-size: 36rpx; font-weight: bold; color: #333;
    margin-bottom: 30rpx; display: block; text-align: center;
  }

  .section-label {
    font-size: 24rpx; color: #999; margin-bottom: 16rpx; display: block;
  }

  .template-grid {
    display: flex; flex-wrap: wrap; gap: 16rpx;
    margin-bottom: 30rpx;

    .template-chip {
      padding: 12rpx 24rpx;
      border-radius: 30rpx;
      background: #f5f5f5;
      display: flex; flex-direction: column; align-items: center;
      min-width: 120rpx;

      .chip-title { font-size: 24rpx; color: #333; }
      .chip-amount { font-size: 28rpx; font-weight: bold; margin-top: 4rpx; }

      &.plus {
        background: #FFF0E6;
        .chip-amount { color: #FF6B35; }
      }
      &.minus {
        background: #F5F5F5;
        .chip-amount { color: #666; }
      }
    }
  }

  .empty-templates {
    padding: 20rpx; text-align: center; color: #999; font-size: 24rpx;
    background: #f9f9f9; border-radius: 12rpx; margin-bottom: 30rpx;
  }

  .divider {
    height: 1px; background: #eee; margin: 20rpx 0;
  }

  .manual-section {
    .input-box {
      background: #f8f9fa; border: 1px solid #ddd; border-radius: 16rpx;
      padding: 20rpx; margin-bottom: 20rpx; font-size: 30rpx; width: 100%; box-sizing: border-box;
      display: block;
    }

    .btn-submit {
      background: #FF6B35; color: white; border-radius: 40rpx;
      height: 80rpx; line-height: 80rpx; text-align: center;
      font-size: 28rpx; font-weight: bold;
    }
  }
}
</style>
