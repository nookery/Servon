---
title: 基本信息
description: 返回当前操作系统的基本信息
---

本文档描述了获取系统基本信息的 API 接口。

## 获取基本系统信息

获取系统的基本运行信息，包括操作系统信息等基础数据。

### 请求

```http
GET /api/system/basic
```

### 请求头

| 名称          | 必填  | 描述                |
|--------------|-------|-------------------|
| Content-Type | 是    | application/json  |

### 响应

#### 成功响应

- 状态码：200 OK
- 内容类型：application/json

响应示例：

```json
{
    // 具体返回字段取决于系统实现
    "os": "linux",
    "arch": "amd64",
    // ... 其他系统基本信息
}
```

#### 错误响应

- 状态码：500 Internal Server Error
- 内容类型：application/json

```json
{
    "error": "错误信息描述"
}
```

### CORS 支持

该 API 支持跨域请求，允许以下配置：

- 允许所有来源 (`Access-Control-Allow-Origin: *`)
- 支持的方法：GET, POST, PUT, DELETE, OPTIONS
- 支持的请求头：Content-Type, Authorization

