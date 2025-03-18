# 🌐 MY_PORTFOLIO

> **一個基於 Golang Gin + HTML/CSS/JavaScript 的全端專案，部署於 AWS。**

**MY_PORTFOLIO** 是一個個人作品展示網站

使用 **Golang Gin** 作為後端框架，前端採用 **HTML、CSS、JavaScript**，並透過 **Nginx** 進行反向代理。

專案使用 **PostgreSQL** 作為資料庫，並透過 **GitHub Actions** 進行自動化部署。

---

## 📂 **專案目錄結構**
```plaintext
MY_PORTFOLIO/
├── .github/workflows/                  # GitHub Actions 自動化部署設定
│   └── deploy.yaml                     # CI/CD 部署設定
│
├── backend/                            # 後端 Golang Gin API
│   ├── database/                       # 資料庫相關檔案
│   │   └── db.go                       # 資料庫連線設定
│   ├── handlers/                       # 處理請求的函式
│   │   ├── auth.go                     # 使用者身份驗證邏輯
│   │   └── user.go                     # 使用者相關 API
│   ├── middlewares/                    # 中介層（Middleware）
│   │   └── auth_middleware.go          # JWT 驗證機制
│   ├── models/                         # 資料庫模型
│   │   └── user.go                     # 使用者模型
│   ├── routes/                         # API 路由設定
│   │   ├── auth_routes.go              # 身份驗證路由
│   │   ├── user_routes.go              # 使用者相關 API 路由
│   │   └── web_routes.go               # 靜態文件與前端路由
│   ├── .env                            # 環境變數設定檔
│   ├── go.mod                          # Golang 依賴管理
│   ├── go.sum                          # 依賴版本管理
│   └── main.go                         # 伺服器進入點
│
├── frontend/                           # 前端 HTML / CSS / JavaScript
│   ├── assets/                         # 前端靜態資源
│   │   ├── admin_styles.css            # 管理員界面樣式
│   │   ├── auth.js                     # 登入/驗證邏輯
│   │   ├── check_login.js              # 檢查使用者是否登入
│   │   ├── hind_page.css               # 隱藏頁面樣式
│   │   ├── login_button.js             # 登入按鈕邏輯
│   │   ├── styles.css                  # 全站 CSS 樣式
│   │   └── user_management.js          # 使用者管理 JS
│   ├── about.html                      # 關於我們頁面
│   ├── contact.html                    # 聯絡我們頁面
│   ├── index.html                      # 首頁
│   ├── projects.html                   # 專案展示頁面
│   ├── skills.html                     # 技能介紹頁面
│   ├── user_management_dashboard.html  # 管理後台儀表板
│   └── user_management.html            # 使用者管理頁面
│
├── docker-compose.yml                  # Docker 多容器部署設定
├── Dockerfile.backend                  # Golang 後端 Docker 設定
├── Dockerfile.frontend                 # 前端 Docker 設定
└── README.md                           # 本文件
```
---
## 🛠 **技術棧**
### 🔹 **後端**
-  **Golang Gin** - 高效能 Web 框架
-  **PostgreSQL** - 關聯式資料庫，存儲使用者與專案數據
-  **JWT (JSON Web Token)** - 使用者身份驗證
-  **Nginx** - 反向代理與靜態檔案伺服器

### 🔹 **前端**
-  **HTML / CSS / JavaScript** - 基礎前端開發
-  **AJAX (Fetch API)** - 用於前後端非同步請求

### 🔹 **DevOps & 部署**
-  **AWS EC2** - 伺服器主機
-  **AWS Route 53** - 網域名稱解析
-  **Nginx** - 反向代理伺服器
-  **GitHub Actions** - CI/CD 自動部署
-  **Docker & Docker Compose** - 容器化應用

---
