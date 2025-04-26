# 🎯 贡献指南（Contributing Guide）
首先感谢对 Autoscale Controller 的关注和贡献！

我们欢迎任何形式的贡献，包括但不限于：

 + 提交 Issue 报告错误或提出建议
 + 修复 Bug
 + 补充功能
 + 优化文档
 + 代码规范性改进
 + 测试用例增强

## ⚙️ 如何开始贡献
1. Fork 仓库
点击右上角 `Fork` 按钮，将项目 Fork 到你的 GitHub 账号下。

2. 克隆代码仓库
```bash
git clone https://github.com/infraflows/autoscale-controller.git
cd autoscale-controller
```

3. 创建特性分支
```bash
git checkout -b feature/feature-name
```

4. 提交修改

 + 请确保每次提交具有清晰的 Commit Message
 + 遵循 Conventional Commits 规范（例如：fix: 修复删除 VPA 时的错误处理、feat: 支持 StatefulSet 的 VPA 扩缩容）
 + 尽可能补充单元测试或验证案例

5. 推送到你的远程仓库
```bash
git push origin feature/your-feature-name
```

6. 提交 Pull Request (PR)
 + 提交到 main 分支
 + 请在 PR 中简要描述你的变更内容和动机
 + 保持 PR 小而清晰，便于 Review

## 🧹 代码规范
 + 保持风格一致，遵循 Go 官方 gofmt 格式化
 + 遵循项目现有目录结构和模块划分。
 + 确保本地通过基础构建和测试命令：
```bash
make build
make test
```
 + 使用 `golangci-lint` 进行代码检查，确保代码符合最佳实践。
 + 使用 `go vet` 检查代码中的潜在问题。
 + 控制器相关逻辑尽量保持幂等性（Idempotent），符合 Kubernetes Operator 开发最佳实践。

## ✅ Pull Request 合并标准
 + 所有单元测试通过
 + 通过 Reviewer 审核（至少 1 个 Approve）
 + Commit Message 清晰、合理
 + 文档更新同步（如果功能有影响）

## 📚 参考资料
[Kubernetes Controller 开发指南](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

[Kubernetes API conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)

[Operator SDK 文档](https://sdk.operatorframework.io/docs/)

## 🚀 欢迎加入我们！
如果你有兴趣持续参与 Autoscale Controller 的开发与优化，非常欢迎加入我们的开发讨论！

