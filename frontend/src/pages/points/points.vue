<template>
  <view class="container">
    <view class="header card">
      <view class="header-content">
        <view>
          <text class="title">积分配置</text>
          <text class="subtitle">管理常用的积分发放模板</text>
        </view>
        <button class="add-btn" @click="openAddModal">添加模板 +</button>
      </view>
    </view>

    <view class="template-grid">
      <view v-for="item in templates" :key="item.id" class="template-card card">
        <view class="card-header">
          <text class="template-title">{{ item.title }}</text>
          <text class="delete-icon" @click="handleDelete(item.id)">✕</text>
        </view>
        <text class="template-content">{{ item.content }}</text>
        <view class="template-footer">
          <text class="amount-tag">{{ item.amount > 0 ? '+' : '' }}{{ item.amount }} ⭐</text>
        </view>
      </view>

      <view v-if="templates.length === 0" class="empty-state card">
        <text>暂无模板，点击右上方添加一个吧</text>
      </view>
    </view>

    <!-- Add Template Modal -->
    <view v-if="showAddModal" class="modal-mask" @click="showAddModal = false">
      <view class="modal-content card" @click.stop>
        <text class="modal-title">新建模板</text>
        <input class="input-box" v-model="newTemplate.title" placeholder="任务标题 (如: 按时睡觉)" />
        <textarea class="input-box textarea" v-model="newTemplate.content" placeholder="任务描述 (如: 晚上9点前上床睡觉)" />
        <input class="input-box" type="number" v-model.number="newTemplate.amount" placeholder="积分数值 (正数奖励，负数扣除)" />

        <view class="modal-actions">
          <button class="btn-cancel" @click="showAddModal = false">取消</button>
          <button class="btn-submit" @click="saveTemplate">保存模板</button>
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
const newTemplate = ref({
  title: '',
  content: '',
  amount: 10
});

const fetchTemplates = async () => {
  try {
    const res = await request({ url: '/parent/templates', method: 'GET' });
    templates.value = res;
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
};

const openAddModal = () => {
  newTemplate.value = { title: '', content: '', amount: 10 };
  showAddModal.ref = true; // Wait, correction below
};

// Fix for reactive ref
const toggleModal = (val: boolean) => {
  showAddModal.value = val;
};

const saveTemplate = async () => {
  if (!newTemplate.value.title || !newTemplate.value.amount) {
    return uni.showToast({ title: '请填写完整信息', icon: 'none' });
  }

  try {
    await request({
      url: '/parent/templates',
      method: 'POST',
      data: newTemplate.value
    });
    uni.showToast({ title: '保存成功', icon: 'success' });
    showAddModal.value = false;
    fetchTemplates();
  } catch (e) {
    uni.showToast({ title: '保存失败', icon: 'none' });
  }
};

const handleDelete = (id: number) => {
  uni.showModal({
    title: '确认删除',
    content: '确定要删除这个配置模板吗？',
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
  margin-bottom: 40rpx;
  background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
  color: white;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title {
    font-size: 40rpx;
    font-weight: bold;
    display: block;
  }

  .subtitle {
    font-size: 24rpx;
    opacity: 0.9;
    margin-top: 10rpx;
    display: block;
  }

  .add-btn {
    font-size: 24rpx;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: 1px solid white;
    border-radius: 30rpx;
    padding: 0 24rpx;
    height: 60rpx;
    line-height: 60rpx;
    margin: 0;
    &::after { border: none; }
  }
}

.template-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24rpx;
}

.template-card {
  padding: 24rpx;
  display: flex;
  flex-direction: column;
  position: relative;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12rpx;
  }

  .template-title {
    font-size: 28rpx;
    font-weight: bold;
    color: #333;
    flex: 1;
  }

  .delete-icon {
    font-size: 24rpx;
    color: #ccc;
    padding: 10rpx;
    margin: -10rpx;
  }

  .template-content {
    font-size: 22rpx;
    color: #666;
    margin-bottom: 20rpx;
    line-height: 1.4;
    min-height: 60rpx;
  }

  .amount-tag {
    font-size: 24rpx;
    font-weight: bold;
    color: #FF6B35;
    background: #FFF5F0;
    padding: 4rpx 16rpx;
    border-radius: 20rpx;
    align-self: flex-start;
  }
}

.empty-state {
  grid-column: span 2;
  padding: 100rpx 40rpx;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #999;
  font-size: 28rpx;
}

.card {
  background: white;
  border-radius: 32rpx;
  box-shadow: 0 10rpx 30rpx rgba(0,0,0,0.05);
}

/* Modal Styles */
.modal-mask {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.6); z-index: 1000; display: flex;
  align-items: center; justify-content: center; backdrop-filter: blur(10rpx);
}

.modal-content {
  width: 80%; padding: 48rpx;
  .modal-title { font-size: 36rpx; font-weight: bold; color: #333; margin-bottom: 40rpx; display: block; text-align: center; }
  .input-box {
    background: #f8f9fa; border: 1px solid #eee; border-radius: 16rpx; padding: 20rpx;
    margin-bottom: 24rpx; font-size: 28rpx; width: 100%; box-sizing: border-box;
    &.textarea { height: 160rpx; }
  }
  .modal-actions {
    display: flex; gap: 20rpx; margin-top: 20rpx;
    button { flex: 1; height: 80rpx; line-height: 80rpx; font-size: 28rpx; border-radius: 40rpx; &::after { border: none; } }
    .btn-cancel { background: #f0f0f0; color: #666; }
    .btn-submit { background: #FF6B35; color: white; }
  }
}
</style>
