summary:
  # 注册与登录优化

  1. **新增注册页面**:
     - 前端新增 `pages/register/register.vue`，包含用户名、密码、身份选择（家长/孩子）和昵称。
     - 在登录页新增“还没有账号？点击注册”入口。

  2. **后端接口实现**:
     - 新增 `/register` 路由 (`router.go`)。
     - 新增 `Register` 处理函数 (`handlers.go`)，自动生成 6 位登录码。
     - 更新 `models.go` 增加 `Password` 字段。
     - 注册成功后自动跳转到登录页。

  3. **细节优化**:
     - 调整了 `records.vue` 中按钮间距（确认与取消按钮间隔 12rpx）。
     - 优化了注册成功后的用户引导流程。
