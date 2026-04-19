# Baby-Fans 积分管理系统后端接口文档

这是一个基于 Go (Gin) 和 SQLite 的后端系统，专为家长和孩子设计的积分与虚拟商店管理平台。

## 🚀 快速启动

### 运行环境
- Go 1.18+
- 系统已配置纯 Go SQLite 驱动，无需安装 GCC。

### 启动服务器
```bash
go run cmd/server/main.go
```
服务器将运行在：`http://localhost:18080`

---

## 🔑 身份认证

系统使用 JWT 进行身份认证。
1. 通过登录接口获取 `token`。
2. 在后续请求的 Header 中携带：`Authorization: Bearer <your_token>`

---

## 📡 接口清单

### 1. 公开接口 (Public)

#### 人脸登录模拟 (自动注册)
- **URL**: `/login/face`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`
- **Body**:
  - `name`: 孩子姓名 (string)
  - `photo`: 照片文件 (file)
- **说明**: 首次登录会自动注册为孩子角色。系统会自动执行**照片轮转**逻辑（仅保留最近 7 张）。返回 `token` 和 6 位离线登录码。

#### 验证码登录
- **URL**: `/login/code`
- **Method**: `GET`
- **Query**: `code=123456`
- **返回**: `token`

---

### 2. 家长端接口 (Parent Role)
*请求头需包含 Parent 角色的 Token*

#### 积分管理 (加分/减分)
- **URL**: `/parent/points/manage`
- **Method**: `POST`
- **Body**:
```json
{
  "user_id": 1,
  "amount": 10,
  "reason": "按时完成作业"
}
```

#### 确认兑换
- **URL**: `/parent/redemption/confirm/:id`
- **Method**: `POST`
- **说明**: 家长确认孩子兑换的商品已发放。

---

### 3. 孩子端接口 (Child Role)
*请求头需包含 Child 角色的 Token*

#### 兑换商品
- **URL**: `/child/exchange/:id`
- **Method**: `POST`
- **说明**: 孩子使用积分兑换虚拟商店中的商品。系统会自动检查积分余额和库存。

---

## 🛠️ 核心功能说明

1. **照片轮转逻辑**:
   - 存储在 `uploads/` 目录下。
   - 每个孩子最多保留 7 张照片，第 8 次登录时会自动删除最旧的一张（包括文件和数据库记录）。
2. **积分事务**:
   - 积分变动与历史记录 (`PointsRecord`) 强绑定，确保数据一致性。
3. **后台自动任务**:
   - 系统每小时运行一次清理任务，自动删除库存为 0 且超过 24 小时的商品。

## 🧪 开发测试
运行单元测试验证核心逻辑：
```bash
go test -v ./internal/service/...
```
