<script setup lang="ts">
import { onLaunch } from "@dcloudio/uni-app";
import { request } from './utils/request';

// App version - synced with manifest.json
const APP_VERSION = '100';

onLaunch(() => {
  console.log("App Launch");
  // Check and apply pending update first
  applyPendingUpdate();
  checkForUpdate();
});

interface VersionInfo {
  version: string;
  build: string;
  update_url: string;
  force_update: boolean;
}

const checkForUpdate = async () => {
  try {
    const res = await request<VersionInfo>({
      url: '/version',
      method: 'GET'
    });

    const latestBuild = parseInt(res.build || '0');
    const currentBuild = parseInt(APP_VERSION);

    if (latestBuild > currentBuild && res.update_url) {
      // Silent download in background
      downloadUpdate(res.update_url);
    }
  } catch (e) {
    console.error('Version check failed:', e);
  }
};

const downloadUpdate = (url: string) => {
  uni.downloadFile({
    url: url,
    success: (downloadRes) => {
      if (downloadRes.statusCode === 200) {
        // Store for next launch
        uni.setStorageSync('pendingUpdate', downloadRes.tempFilePath);
      }
    },
    fail: (err) => {
      console.error('Download failed:', err);
    }
  });
};

const applyPendingUpdate = () => {
  const pendingPath = uni.getStorageSync('pendingUpdate') as string;
  if (!pendingPath) return;

  uni.removeStorageSync('pendingUpdate');

  if (uni.canIUse('uni.installBundle')) {
    uni.installBundle({
      url: pendingPath,
      fail: () => {
        // Fallback
        plus.runtime.install(pendingPath, { force: true }, () => {
          uni.reLaunch({ url: '/' });
        }, () => {});
      }
    });
  } else {
    plus.runtime.install(pendingPath, { force: true }, () => {
      uni.reLaunch({ url: '/' });
    }, () => {});
  }
};
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
