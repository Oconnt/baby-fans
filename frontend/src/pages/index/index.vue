<template>
  <view class="container">

    <!-- Parent View: Bind Child -->
    <view v-if="userRole === 'parent'" class="bind-section">
      <view class="bind-row">
        <input class="bind-input" v-model="loginCode" placeholder="输入孩子登录码" />
        <button class="bind-btn" @click="bindChild">绑定</button>
      </view>
    </view>

    <!-- Parent View: Children List -->
    <view v-if="userRole === 'parent'" class="children-list">
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
      <view v-if="childrenList.length === 0" class="empty-state card">
        <text>暂无绑定的孩子，请输入登录码绑定</text>
      </view>
    </view>

    <!-- Child View: Points Overview -->
    <view v-if="userRole === 'child'" class="child-overview">
      <view class="points-circle" @click="toggleTasks">
        <text class="points-label">当前积分</text>
        <text class="points-value">{{ childOverview.points }}</text>
        <text class="toggle-hint">{{ showTasks ? '点击隐藏任务' : '点击查看今日任务' }}</text>
      </view>

      <view class="parent-info card">
        <text class="parent-label">绑定家长</text>
        <view class="parent-names">
          <text v-if="childOverview.parent_names && childOverview.parent_names.length">
            {{ childOverview.parent_names.join('、') }}
          </text>
          <text v-else>未绑定</text>
        </view>
      </view>

      <!-- Today's Tasks Section -->
      <view v-if="showTasks" class="task-section">
        <text class="section-title">今日任务</text>
        <view v-if="todayTasks.length > 0" class="task-list">
          <view v-for="task in todayTasks" :key="task.id" class="task-item card" @click="showTaskDetail(task)">
            <view class="task-info">
              <text class="task-name">{{ task.name }}</text>
              <text class="task-meta">
                <text class="task-status" :class="getStatusClass(task.status)">{{ getStatusText(task.status) }}</text>
                <text class="task-expire"> | 过期: {{ formatTime(task.expire_time) }}</text>
              </text>
            </view>
            <view class="task-points">+{{ task.points }}</view>
          </view>
        </view>
        <view v-else class="empty-tip">今日暂无任务</view>
      </view>

      <view v-if="showRecords" class="record-section">
        <text class="section-title">最近记录</text>
        <view class="record-list">
          <view v-for="item in childOverview.records" :key="item.id" class="record-item card">
            <view class="record-left">
              <text class="record-reason">{{ item.reason }}</text>
              <text class="record-operator">{{ item.operator ? (item.operator.nickname || item.operator.name) : '' }}</text>
            </view>
            <text class="record-amount" :class="item.amount >= 0 ? 'plus' : 'minus'">
              {{ item.amount >= 0 ? '+' : '' }}{{ item.amount }}
            </text>
          </view>
          <view v-if="childOverview.records.length === 0" class="empty-tip">暂无记录</view>
        </view>
      </view>
    </view>

    <!-- Task Detail Modal -->
    <view v-if="showTaskModal" class="modal-mask" @click="showTaskModal = false">
      <view class="modal-content card" @click.stop>
        <view class="modal-header">
          <text class="modal-title">任务详情</text>
          <view class="modal-close" @click="showTaskModal = false">✕</view>
        </view>
        <view v-if="selectedTask" class="task-detail">
          <view class="detail-row">
            <text class="detail-label">任务名称</text>
            <text class="detail-value">{{ selectedTask.name }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">任务描述</text>
            <text class="detail-value">{{ selectedTask.description || '无' }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">奖励积分</text>
            <text class="detail-value points">+{{ selectedTask.points }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">状态</text>
            <text class="detail-value" :class="getStatusClass(selectedTask.status)">{{ getStatusText(selectedTask.status) }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">发布时间</text>
            <text class="detail-value">{{ formatTime(selectedTask.publish_time) }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">过期时间</text>
            <text class="detail-value">{{ formatTime(selectedTask.expire_time) }}</text>
          </view>
          <view v-if="selectedTask.status === 1" class="btn-complete" @click="completeTask(selectedTask)">完成任务</view>
        </view>
      </view>
    </view>

    <!-- Points Adjust Modal (Parent Only) -->
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
import { onShow } from '@dcloudio/uni-app';
import { request } from '../../utils/request';

    const userRole = ref('');
    const childrenList = ref<any[]>([]);
    const loginCode = ref('');

    // Child Overview Data
    const childOverview = ref<any>({ points: 0, records: [], parent_names: [] });
    const showRecords = ref(false);

    // Task related
    const showTasks = ref(false);
    const todayTasks = ref<any[]>([]);
    const showTaskModal = ref(false);
    const selectedTask = ref<any>(null);

    // Parent Modal Data
    const showPointsModal = ref(false);
    const selectedChild = ref<any>(null);
    const pointTemplates = ref<any[]>([]);
    const manualAmount = ref('');

    const updateRoleAndData = () => {
      const stored = uni.getStorageSync('userInfo');
      if (!stored) {
        uni.reLaunch({ url: '/pages/login/login' });
        return;
      }
      const userInfo = JSON.parse(stored);
      userRole.value = userInfo.role;

      if (userRole.value === 'parent') {
        uni.setNavigationBarTitle({ title: '孩子管理' });
        fetchChildren();
      } else {
        uni.setNavigationBarTitle({ title: '积分详情' });
        fetchChildOverview();
      }
    };

    onMounted(() => {
      const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
      if (userInfo.role === 'parent') {
        uni.switchTab({ url: '/pages/home-parent/home-parent' });
        return;
      }
      updateRoleAndData();
    });

    onShow(() => {
      updateRoleAndData();
    });

    const onTabItemTap = () => {
      updateRoleAndData();
    };

    const fetchChildOverview = async () => {
      try {
        const res = await request({ url: '/child/overview', method: 'GET' });
        console.log('child overview response:', res);
        childOverview.value = {
          points: res?.points || 0,
          records: res?.records || [],
          parent_names: res?.parent_names || []
        };
        console.log('childOverview.value:', childOverview.value);
      } catch (e) {
        console.error(e);
      }
    };

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

    // Task functions
    const toggleTasks = async () => {
      console.log('toggleTasks called, showTasks before:', showTasks.value);
      showTasks.value = !showTasks.value;
      showRecords.value = false; // Hide records when showing tasks
      console.log('showTasks after:', showTasks.value);
      if (showTasks.value && todayTasks.value.length === 0) {
        await fetchTodayTasks();
      }
    };

    const fetchTodayTasks = async () => {
      try {
        const res = await request({ url: '/child/tasks/today', method: 'GET' });
        todayTasks.value = res || [];
      } catch (e) {
        console.error(e);
        todayTasks.value = [];
      }
    };

    const showTaskDetail = async (task: any) => {
      selectedTask.value = task;
      showTaskModal.value = true;
    };

    const completeTask = async (task: any) => {
      try {
        await request({
          url: `/child/tasks/${task.id}/complete`,
          method: 'PUT'
        });
        uni.showToast({ title: '任务完成，获得 ' + task.points + ' 积分', icon: 'success' });
        showTaskModal.value = false;
        // Refresh data
        fetchTodayTasks();
        fetchChildOverview();
      } catch (e) {
        uni.showToast({ title: '操作失败', icon: 'none' });
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

.child-overview {
  .points-circle {
    padding: 80rpx;
    text-align: center;
    margin-bottom: 40rpx;
    background: linear-gradient(135deg, #FF6B35 0%, #FFB347 100%);
    color: white;
    border-radius: 50%; /* Make it circular */
    width: 400rpx; /* Fixed width */
    height: 400rpx; /* Fixed height */
    display: flex; /* Use flex to center content */
    flex-direction: column;
    justify-content: center;
    align-items: center;
    margin: 40rpx auto; /* Center horizontally */
    box-shadow: 0 8rpx 20rpx rgba(255, 107, 53, 0.3);

    .points-label { font-size: 28rpx; opacity: 0.9; margin-bottom: 10rpx; }
    .points-value { font-size: 80rpx; font-weight: bold; display: block; margin-bottom: 10rpx; }
    .toggle-hint { font-size: 22rpx; opacity: 0.8; }
  }

  .record-section {
    margin-top: 30rpx;
  }

  .parent-info {
    padding: 30rpx;
    background: white;
    border-radius: 16rpx;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30rpx;
    .parent-label { font-size: 28rpx; color: #999; }
    .parent-name { font-size: 32rpx; font-weight: bold; color: #333; }
  }

  .section-title { font-size: 30rpx; font-weight: bold; color: #333; margin-bottom: 20rpx; display: block; }
  .record-list {
    display: flex; flex-direction: column; gap: 20rpx;
    .record-item {
      padding: 24rpx; display: flex; justify-content: space-between; align-items: center;
      .record-left {
        flex: 1;
        .record-reason { font-size: 28rpx; color: #333; display: block; }
        .record-operator { font-size: 24rpx; color: #999; display: block; margin-top: 4rpx; }
      }
      .record-amount { font-size: 32rpx; font-weight: bold; }
      .plus { color: #FF6B35; }
      .minus { color: #666; }
    }
  }
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

.task-section {
  margin-top: 30rpx;
  .task-list {
    display: flex; flex-direction: column; gap: 20rpx;
    .task-item {
      padding: 24rpx; display: flex; justify-content: space-between; align-items: center;
      .task-info {
        flex: 1;
        .task-name { font-size: 28rpx; color: #333; display: block; font-weight: bold; }
        .task-meta { font-size: 22rpx; color: #999; display: block; margin-top: 4rpx; }
        .task-status {
          &.pending { color: #007AFF; }
          &.completed { color: #4CD964; }
          &.expired { color: #FF3B30; }
        }
      }
      .task-points { font-size: 32rpx; font-weight: bold; color: #FF6B35; }
    }
  }
}

.task-detail {
  .detail-row {
    display: flex; justify-content: space-between; align-items: center;
    padding: 20rpx 0; border-bottom: 1px solid #f0f0f0;
    .detail-label { font-size: 26rpx; color: #999; }
    .detail-value { font-size: 28rpx; color: #333; font-weight: 500; }
    .detail-value.points { color: #FF6B35; font-weight: bold; }
    .detail-value.pending { color: #007AFF; }
    .detail-value.completed { color: #4CD964; }
    .detail-value.expired { color: #FF3B30; }
  }
  .btn-complete {
    margin-top: 30rpx; background: #FF6B35; color: white; border-radius: 40rpx;
    height: 80rpx; line-height: 80rpx; text-align: center;
    font-size: 28rpx; font-weight: bold;
  }
}

.modal-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 30rpx;
  .modal-title { font-size: 36rpx; font-weight: bold; color: #333; text-align: center; flex: 1; }
  .modal-close { font-size: 36rpx; color: #999; padding: 10rpx; }
}

.empty-tip {
  text-align: center; color: #999; font-size: 26rpx; padding: 40rpx;
}
</style>
