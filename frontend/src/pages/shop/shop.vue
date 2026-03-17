<template>
  <view class="container">
    <!-- Parent View: Item Management -->
    <view v-if="userRole === 'parent'" class="admin-view">
      <view class="card add-card" @click="openAddModal">
        <text class="plus">+</text>
        <text class="label">上架新商品</text>
      </view>

      <view class="item-list">
        <view v-for="item in items" :key="item.id" class="card item-card">
          <view class="details">
            <text class="name">{{ item.name }}</text>
            <text class="desc" v-if="item.description">{{ item.description }}</text>
            <text class="price">{{ item.price }} ⭐ | 库存: {{ item.stock }}</text>
          </view>
          <view class="actions">
            <text class="stock-btn" @click="openEditModal(item)">修改</text>
            <text class="del-btn" @click="deleteItem(item.id)">下架</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Child View: Wish Shop -->
    <view v-else class="shop-view">
      <view class="item-grid">
        <view v-for="item in items" :key="item.id" class="card product-card">
          <text class="product-name">{{ item.name }}</text>
          <text class="product-desc" v-if="item.description">{{ item.description }}</text>
          <view class="price-row">
            <text class="price">{{ item.price }}</text>
            <text class="unit">⭐</text>
          </view>
          <view class="stock-row">
            <text :class="['stock', item.stock === 0 ? 'out' : '']">库存: {{ item.stock }}</text>
          </view>
          <button class="buy-btn" :class="{ disabled: item.stock === 0 }" :disabled="item.stock === 0" @click="handleExchange(item)">{{ item.stock === 0 ? '已售罄' : '兑换' }}</button>
        </view>
      </view>
    </view>

    <!-- Add/Edit Item Modal -->
    <view v-if="showAddItem" class="modal-mask" @click="closeModal">
      <view class="modal-content card" @click.stop>
        <text class="modal-title">{{ editingItem?.id ? '修改商品' : '上架新商品' }}</text>

        <view class="form-group">
          <text class="label">商品名称</text>
          <input class="input-box" v-model="editingItem.name" placeholder="请输入商品名称" placeholder-style="color: #999;" />
        </view>

        <view class="form-group">
          <text class="label">商品描述</text>
          <input class="input-box" v-model="editingItem.description" placeholder="请输入商品描述" placeholder-style="color: #999;" />
        </view>

        <view class="form-group">
          <text class="label">所需积分</text>
          <input class="input-box" type="number" v-model.number="editingItem.price" placeholder="请输入所需积分" placeholder-style="color: #999;" />
        </view>

        <view class="form-group">
          <text class="label">库存数量</text>
          <input class="input-box" type="number" v-model.number="editingItem.stock" placeholder="请输入库存数量" placeholder-style="color: #999;" />
        </view>

        <view class="modal-actions">
          <button class="btn-cancel" @click="closeModal">取消</button>
          <button class="btn-submit" @click="saveItem">{{ editingItem?.id ? '保存修改' : '确认上架' }}</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { request } from '../../utils/request';

const userRole = ref('');
const items = ref<any[]>([]);
const showAddItem = ref(false);

const loadData = async () => {
  const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
  userRole.value = userInfo.role;
  try {
    const res = await request({ url: '/parent/items', method: 'GET' });
    items.value = res || [];
  } catch (e) {
    items.value = [];
  }
};

const openAddModal = () => {
  editingItem.value = { name: '', description: '', price: 0, stock: 10 };
  showAddItem.value = true;
};

const openEditModal = (item: any) => {
  editingItem.value = { ...item }; // Copy original values
  showAddItem.value = true;
};

const closeModal = () => {
  showAddItem.value = false;
  editingItem.value = { name: '', description: '', price: 0, stock: 10 };
};

const saveItem = async () => {
  if (!editingItem.value.name) {
    return uni.showToast({ title: '请填写商品名称', icon: 'none' });
  }
  if (editingItem.value.price <= 0) {
    return uni.showToast({ title: '所需积分必须大于0', icon: 'none' });
  }
  if (editingItem.value.stock < 0) {
    return uni.showToast({ title: '库存不能为负数', icon: 'none' });
  }
  try {
    await request({
      url: '/parent/items',
      method: 'POST',
      data: editingItem.value
    });
    uni.showToast({ title: '保存成功', icon: 'success' });
    closeModal();
    loadData();
  } catch (e) {
    uni.showToast({ title: '保存失败', icon: 'none' });
  }
};

const editingItem = ref<any>(null);

const deleteItem = (id: number) => {
  uni.showModal({
    title: '确认下架',
    content: '确定要下架这个商品吗？',
    success: async (res) => {
      if (res.confirm) {
        await request({ url: `/parent/items/${id}`, method: 'DELETE' });
        loadData();
      }
    }
  });
};

const handleExchange = (item: any) => {
  uni.showModal({
    title: '确认兑换',
    content: `确定使用 ${item.price} 积分兑换 ${item.name} 吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await request({
            url: '/child/exchange',
            method: 'POST',
            data: { item_id: item.id }
          });
          uni.showToast({ title: '兑换成功！', icon: 'success' });
          loadData();
        } catch (e: any) {
          uni.showToast({ title: e?.error || '积分不足', icon: 'none' });
        }
      }
    }
  });
};

onMounted(() => {
  const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
  if (!userInfo.role) {
    uni.reLaunch({ url: '/pages/login/login' });
    return;
  }
  loadData();
});
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

.add-card {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed #FF6B35;
  color: #FF6B35;
  padding: 30rpx;
  margin-bottom: 30rpx;
  background: transparent;
  box-shadow: none;
  .plus { font-size: 40rpx; margin-right: 10rpx; font-weight: bold; }
  .label { font-size: 28rpx; font-weight: bold; }
}

.item-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.item-card {
  display: flex;
  align-items: center;
  padding: 28rpx;
  .details {
    flex: 1;
    .name { font-size: 30rpx; font-weight: bold; color: #333; display: block; }
    .desc { font-size: 24rpx; color: #666; margin-top: 4rpx; display: block; }
    .price { font-size: 24rpx; color: #888; margin-top: 8rpx; display: block; }
  }
  .actions {
    display: flex;
    gap: 20rpx;
    .stock-btn { color: #FF6B35; font-size: 24rpx; font-weight: bold; }
    .del-btn { color: #ccc; font-size: 24rpx; font-weight: bold; }
  }
}

.item-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20rpx;
}

.product-card {
  padding: 28rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  .product-name { font-size: 28rpx; font-weight: bold; margin-bottom: 8rpx; }
  .product-desc { font-size: 22rpx; color: #999; margin-bottom: 12rpx; text-align: center; height: 60rpx; overflow: hidden; }
  .price-row {
    display: flex; align-items: baseline; gap: 4rpx; margin-bottom: 10rpx;
    .price { font-size: 36rpx; font-weight: 900; color: #FF6B35; }
    .unit { font-size: 20rpx; color: #FF6B35; }
  }
  .stock-row {
    margin-bottom: 20rpx;
    .stock { font-size: 22rpx; color: #999; }
    .stock.out { color: #ff4d4f; }
  }
  .buy-btn {
    width: 100%; background: #FF6B35; color: white;
    font-size: 24rpx; font-weight: bold; border-radius: 12rpx;
    padding: 10rpx 0; line-height: 1.5;
    &::after { border: none; }
  }
}

.card {
  background: white;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 20rpx rgba(0,0,0,0.05);
}

.modal-mask {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.5); z-index: 100; display: flex;
  align-items: center; justify-content: center;
}

.modal-content {
  width: 80%; padding: 48rpx; z-index: 101;
  .modal-title { font-size: 36rpx; font-weight: bold; color: #333; margin-bottom: 40rpx; display: block; text-align: center; }
  .form-group { margin-bottom: 20rpx; }
  .label { font-size: 26rpx; color: #666; margin-bottom: 10rpx; display: block; font-weight: bold; }
  .input-box {
    display: block;
    background: #f8f9fa; border: 1px solid #FF6B35; border-radius: 16rpx;
    padding: 20rpx;
    margin-bottom: 24rpx; font-size: 28rpx; width: 100%; box-sizing: border-box;
    color: #333333;
    height: 80rpx;
    caret-color: #333333;
    font-weight: bold;
  }
  .modal-actions {
    display: flex; gap: 20rpx; margin-top: 20rpx;
    button { flex: 1; height: 80rpx; line-height: 80rpx; font-size: 28rpx; border-radius: 40rpx; &::after { border: none; } }
    .btn-cancel { background: #f0f0f0; color: #666; }
    .btn-submit { background: #FF6B35; color: white; }
  }
}
</style>
