<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
onLaunch(() => {
  console.log("App Launch");
});
onShow(() => {
  console.log("App Show");
});
onHide(() => {
  console.log("App Hide");
});
</script>

<style lang="scss">
/* 字体 - 小程序不支持 Google Fonts，使用系统字体 */
:root {
  --yellow: #FFD93D; --orange: #FF6B35; --pink: #FF8FAB; --purple: #C77DFF;
  --blue: #4CC9F0; --green: #52B788; --bg: #FFF9F0; --card: #FFFFFF;
  --text: #2D2D2D; --text2: #888; --radius: 24px; --shadow: 0 8px 32px rgba(0,0,0,0.10);
}

/* 全局基础样式 - 恢复之前的氛围 */
page {
  font-family: 'PingFang SC', 'Helvetica Neue', sans-serif;
  background: linear-gradient(135deg, #FFF9F0 0%, #FFF0F5 50%, #F0F5FF 100%);
  min-height: 100vh;
  color: #2D2D2D;
}

/* 原有的通用卡片样式 */
.card {
  background: white;
  border-radius: 24px;
  padding: 15px;
  margin: 0 16px 12px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.06);
  border: 1px solid rgba(0,0,0,0.02);
}

.btn-primary {
  width: 100%;
  background: linear-gradient(135deg, #FF6B35, #FF8FAB);
  color: white;
  border: none;
  border-radius: 14px;
  padding: 12px;
  font-weight: 800;
  font-size: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(255, 107, 53, 0.3);
  &::after { border: none; }
}

.btn-primary:active {
  transform: scale(0.96);
  opacity: 0.9;
}

/* 标题字体 */
.fredoka {
  font-family: 'PingFang SC', 'Helvetica Neue', cursive;
  font-weight: 900;
}
</style>
