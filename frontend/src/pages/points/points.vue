<template>
  <view class="container">
    <view class="header card">
      <text class="title">标签管理</text>
      <text class="subtitle">管理积分任务标签</text>
    </view>

    <!-- Usage Guide -->
    <view class="usage-guide card">
      <text class="guide-title">💡 使用说明</text>
      <text class="guide-item">1. 在此创建常用任务标签（如：写完作业 +10）。</text>
      <text class="guide-item">2. 创建后在“孩子管理”点击“管理”即可快速使用。</text>
    </view>

    <!-- Template List -->
    <view class="template-list">
      <view v-for="item in templates" :key="item.id" class="template-card card">
        <text class="delete-icon" @click="handleDelete(item.id)">✕</text>
        <text class="template-title">{{ item.title }}</text>
        <text class="template-content">{{ item.content }}</text>
        <text class="template-amount" :class="item.amount >= 0 ? 'plus' : 'minus'">
          {{ item.amount >= 0 ? '+' : '' }}{{ item.amount }}
        </text>
      </view>

      <!-- Add New Card -->
      <view class="template-card card add-card" @click="showAddModal = true">
        <text class="add-icon">+</text>
        <text class="add-text">新增标签</text>
      </view>
    </view>

    <!-- Add Template Modal -->
    <view v-if="showAddModal" class="modal-mask" @click="showAddModal = false">
      <view class="modal-content card" @click.stop>
        <text class="modal-title">新建标签</text>
        <input class="input-box" v-model="newTemplate.title" placeholder="标签名称" />
        <input class="input-box" v-model="newTemplate.content" placeholder="标签描述" />
        <input class="input-box" type="number" v-model.number="newTemplate.amount" placeholder="积分分值 (例如 10 或 -5)" />
        <view class="modal-actions">
          <view class="btn-cancel" @click="showAddModal = false">取消</view>
          <view class="btn-submit" @click="saveTemplate">保存</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { request } from '../../utils/request';

const templates = ref<any[]>([]);
const showAddModal = ref(false);
const newTemplate = ref({ title: '', content: '', amount: 0 });

const fetchTemplates = async () => {
  try {
    const res = await request({ url: '/parent/templates', method: 'GET' });
    templates.value = res || [];
  } catch (e) {
    templates.value = [];
  }
};

const saveTemplate = async () => {
  if (!newTemplate.value.title) {
    return uni.showToast({ title: '请填写标签名称', icon: 'none' });
  }
  try {
    await request({
      url: '/parent/templates',
      method: 'POST',
      data: newTemplate.value
    });
    uni.showToast({ title: '保存成功', icon: 'success' });
    showAddModal.value = false;
    newTemplate.value = { title: '', content: '', amount: 0 };
    fetchTemplates();
  } catch (e) {
    uni.showToast({ title: '保存失败', icon: 'none' });
  }
};

const handleDelete = (id: number) => {
  uni.showModal({
    title: '确认删除',
    content: '确定要删除这个标签吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({ url: `/parent/templates/${id}`, method: 'DELETE' });
          fetchTemplates();
        } catch (e) {
          uni.showToast({ title: '删除失败', icon: 'none' });
        }
      }
    }
  });
};

onMounted(fetchTemplates);
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
  background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
  color: white;
  .title { font-size: 40rpx; font-weight: bold; display: block; }
  .subtitle { font-size: 24rpx; opacity: 0.9; margin-top: 8rpx; display: block; }
}

.usage-guide {
  padding: 30rpx;
  margin-bottom: 30rpx;
  background: #fff9f6;
  border: 1px solid #ffe8de;
  .guide-title { font-size: 28rpx; font-weight: bold; color: #FF6B35; margin-bottom: 12rpx; display: block; }
  .guide-item { font-size: 24rpx; color: #888; display: block; line-height: 1.6; }
}

.template-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
}

.template-card {
  padding: 28rpx;
  position: relative;
  display: flex;
  flex-direction: column;

  .delete-icon {
    position: absolute;
    top: 16rpx;
    right: 16rpx;
    font-size: 28rpx;
    color: #ccc;
    padding: 8rpx;
  }

  .template-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #333;
    margin-bottom: 12rpx;
    padding-right: 40rpx;
  }

  .template-content {
    font-size: 24rpx;
    color: #888;
    line-height: 1.5;
    margin-bottom: 12rpx;
  }

  .template-amount {
    font-size: 32rpx;
    font-weight: bold;
    align-self: flex-end;
    &.plus { color: #FF6B35; }
    &.minus { color: #888; }
  }
}

.add-card {
  align-items: center;
  justify-content: center;
  border: 2px dashed #ddd;
  background: transparent;
  box-shadow: none;
  min-height: 160rpx;

  .add-icon { font-size: 60rpx; color: #ccc; }
  .add-text { font-size: 24rpx; color: #bbb; margin-top: 8rpx; }
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
  width: 80%; padding: 48rpx; position: relative;
  background: #ffffff; border-radius: 24rpx;
  pointer-events: auto;
  .modal-title { font-size: 36rpx; font-weight: bold; color: #333; margin-bottom: 40rpx; display: block; text-align: center; }
  .input-box {
    background: #f8f9fa; border: 1px solid #ddd; border-radius: 16rpx; padding: 24rpx;
    margin-bottom: 24rpx; font-size: 30rpx; width: 100%; box-sizing: border-box;
    display: block; height: 90rpx;
  }
  .modal-actions {
    display: flex; gap: 20rpx; margin-top: 20rpx;
    .btn-cancel, .btn-submit {
      flex: 1; height: 80rpx; line-height: 80rpx; font-size: 28rpx;
      border-radius: 40rpx; text-align: center;
      cursor: pointer;
    }
    .btn-cancel { background: #f0f0f0; color: #666; }
    .btn-submit { background: #FF6B35; color: white; }
  }
}
</style>
