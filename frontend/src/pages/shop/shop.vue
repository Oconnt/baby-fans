<template>
  <view class="container">
    <!-- Parent View: Item Management -->
    <view v-if="userRole === 'parent'" class="admin-view">
      <view class="header-section">
        <text class="fredoka title">Shop Admin</text>
        <view class="card add-card" @click="showAddItem = true">
          <text class="plus">+</text>
          <text class="label">上架新奖品</text>
        </view>
      </view>

      <view class="item-list">
        <view v-for="item in items" :key="item.id" class="card item-card">
          <view class="icon-box">🎁</view>
          <view class="details">
            <text class="name">{{ item.name }}</text>
            <text class="price">{{ item.price }} ⭐ | 库存: {{ item.stock }}</text>
          </view>
          <view class="actions">
            <text class="del-btn" @click="deleteItem(item.id)">下架</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Child View: Wish Shop -->
    <view v-else class="shop-view">
      <view class="hero-header">
        <text class="fredoka title">Wish Shop</text>
        <text class="subtitle">挑选你心仪的奖品吧 🎁</text>
      </view>

      <view class="item-grid">
        <view v-for="item in items" :key="item.id" class="card product-card">
          <view class="image-box">🎁</view>
          <text class="product-name">{{ item.name }}</text>
          <view class="price-row">
            <text class="price">{{ item.price }}</text>
            <text class="unit">⭐</text>
          </view>
          <button class="buy-btn" @click="handleExchange(item)">兑换</button>
        </view>
      </view>
    </view>

    <!-- Add Item Modal -->
    <view v-if="showAddItem" class="modal-mask" @click="showAddItem = false">
      <view class="modal-content card" @click.stop>
        <text class="modal-title fredoka">New Prize</text>
        <input class="input-box" v-model="newItem.name" placeholder="奖品名称" />
        <input class="input-box" type="number" v-model.number="newItem.price" placeholder="所需积分" />
        <input class="input-box" type="number" v-model.number="newItem.stock" placeholder="初始库存" />
        <button class="btn-primary" @click="saveItem">确认上架</button>
        <text class="cancel-link" @click="showAddItem = false">返回</text>
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
const newItem = ref({ name: '', price: 0, stock: 10 });

const loadData = async () => {
  const userInfo = JSON.parse(uni.getStorageSync('userInfo') || '{}');
  userRole.value = userInfo.role;

  try {
    const res = await request({ url: '/parent/items', method: 'GET' });
    items.value = res;
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' });
  }
};

const saveItem = async () => {
  if (!newItem.value.name || newItem.value.price <= 0) return;
  try {
    await request({
      url: '/parent/items',
      method: 'POST',
      data: newItem.value
    });
    showAddItem.value = false;
    newItem.value = { name: '', price: 0, stock: 10 };
    loadData();
  } catch (e) {
    uni.showToast({ title: '保存失败', icon: 'none' });
  }
};

const deleteItem = (id: number) => {
  uni.showModal({
    title: '确认下架',
    content: '确定要删除这个奖品吗？',
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
            url: `/child/exchange`,
            method: 'POST',
            data: { item_id: item.id }
          });
          uni.showToast({ title: '兑换成功！', icon: 'success' });
          loadData();
        } catch (e: any) {
          uni.showToast({ title: e.error || '积分不足', icon: 'none' });
        }
      }
    }
  });
};

onMounted(loadData);
</script>

<style lang="scss" scoped>
.container {
  padding-bottom: 40rpx;
  min-height: 100vh;
}

.header-section, .hero-header {
  padding: 40rpx 32rpx;
  .title { font-size: 60rpx; color: var(--text); display: block; }
  .subtitle { font-size: 24rpx; color: var(--text2); font-weight: 700; }
}

.add-card {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx dashed var(--purple);
  color: var(--purple);
  margin-top: 20rpx;
  padding: 30rpx;
  .plus { font-size: 40rpx; margin-right: 10rpx; font-weight: bold; }
  .label { font-size: 28rpx; font-weight: 800; }
}

.item-card {
  display: flex;
  align-items: center;
  padding: 24rpx;
  .icon-box {
    width: 80rpx; height: 80rpx; background: #FFF9C4;
    border-radius: 20rpx; display: flex; align-items: center;
    justify-content: center; font-size: 40rpx; margin-right: 20rpx;
  }
  .details {
    flex: 1;
    .name { font-size: 30rpx; font-weight: 800; display: block; }
    .price { font-size: 24rpx; color: var(--text2); font-weight: 700; }
  }
  .del-btn { color: #FF4D4F; font-size: 24rpx; font-weight: 800; }
}

.item-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24rpx;
  padding: 0 32rpx;
}

.product-card {
  margin: 0;
  padding: 24rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  .image-box {
    width: 140rpx; height: 140rpx; background: #FFF9C4;
    border-radius: 30rpx; display: flex; align-items: center;
    justify-content: center; font-size: 60rpx; margin-bottom: 20rpx;
  }
  .product-name { font-size: 28rpx; font-weight: 800; margin-bottom: 10rpx; }
  .price-row {
    display: flex; align-items: baseline; gap: 4rpx; margin-bottom: 20rpx;
    .price { font-size: 36rpx; font-weight: 900; color: var(--orange); }
    .unit { font-size: 20rpx; color: var(--orange); }
  }
  .buy-btn {
    width: 100%; background: var(--orange); color: white;
    font-size: 24rpx; font-weight: 800; border-radius: 12rpx;
    padding: 10rpx 0; line-height: 1.5;
    &::after { border: none; }
  }
}

.modal-mask {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0,0,0,0.6); z-index: 1000; display: flex;
  align-items: center; justify-content: center; backdrop-filter: blur(10rpx);
}

.modal-content {
  width: 80%; padding: 48rpx; text-align: center;
  .modal-title { font-size: 48rpx; color: var(--purple); margin-bottom: 40rpx; display: block; }
  .input-box {
    border: 4rpx solid #F0F0F0; border-radius: 20rpx; padding: 20rpx;
    margin-bottom: 24rpx; font-size: 28rpx; width: 100%; box-sizing: border-box;
  }
  .cancel-link { font-size: 24rpx; color: #BBB; margin-top: 30rpx; display: block; }
}
</style>
