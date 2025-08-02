# 💼 個人作品集系統 (My Portfolio Web System)

這是一個部署於 AWS EC2 的個人作品集與應用展示系統，後端以 **Golang Gin** 為核心，整合 **RESTful API**、**GraphQL**、**WebSocket**，並搭配 **HTML/CSS/JavaScript** 的前端介面，呈現完整的全端技術整合能力。

---

## 🔧 專案技術架構

### 📌 前端技術

* HTML / CSS / JavaScript（純原生，無框架）
* 使用多頁式架構（MPA）：

  * `index.html`: 首頁
  * `about.html`: 關於我
  * `projects.html`: 專案展示
  * `skills.html`: 技能介紹
  * `contact.html`: 聯絡方式
  * `chat.html`: 即時聊天室（WebSocket）
  * `expense.html`: 個人記帳功能（REST）
  * `food.html`: 飲食紀錄系統（GraphQL）
  * `stocks.html`: 台股追蹤分析（REST + WebSocket）
  * `user_management.html`: 使用者登入註冊頁面
  * `user_management_dashboard.html`: 使用者後台 & 管理者審核介面

### 🧠 後端技術

* Golang + Gin 框架
* JWT 驗證系統（登入、註冊、身分驗證、RBAC）
* RESTful API：支援記帳系統、股票紀錄、使用者管理
* GraphQL：用於飲食紀錄功能
* WebSocket：聊天室 / 台股即時報價串流
* PostgreSQL 資料庫
* Excel 檔案匯入匯出功能（記帳、交易）

### 🚀 部署與基礎架構

* 使用 **Docker** 與 **docker-compose** 進行容器化部署
* 部署於 **AWS EC2**
* 使用 **AWS Route 53** 購買獨立網域並指向主機
* 設定 **Nginx** 作為反向代理伺服器

  * 靜態資源與 HTML 頁面由 Nginx 提供
  * `/api/`、`/graphql`、`/ws/` 由 Nginx Proxy 到後端
* HTTPS 憑證可透過 Let's Encrypt 整合

---

## 🧩 系統功能模組說明

### 🗂 使用者管理

* 註冊 / 登入 / 登出（JWT）
* 權限控管：使用者與管理員角色
* 管理者可審核新註冊使用者（通過 / 拒絕）

### 📒 記帳功能（expense.html + RESTful）

* 類別選擇（早餐、午餐、交通等）
* 可輸入金額、備註、日期
* 整合計算機介面（虛擬按鍵）
* 可匯出目前資料為 Excel
* 支援 Excel 上傳並覆蓋所有資料

### 📊 台股追蹤分析（stocks.html + REST/WebSocket）

* 可新增追蹤股票代碼與張數與均價
* 透過 WebSocket 即時更新價格
* 計算報酬率、平均價、自動估算手續費與稅金
* 顯示報酬總結（未實現/已實現/總損益）
* 可記錄賣出紀錄與股息收入，並匯出所有交易為 Excel
* 可開啟 Modal 視覺化每支股票近 30 日股價趨勢圖

### 🍱 飲食紀錄系統（food.html + GraphQL）

* 新增每日攝取內容（名稱、熱量、蛋白質、脂肪、碳水、份量）
* 支援日期篩選
* 動態計算單日總熱量與蛋白質
* 可編輯與刪除任意資料

### 💬 即時聊天室（chat.html + WebSocket）

* 使用 WebSocket 建立聊天室連線
* 可顯示發言者、時間、訊息內容
* 系統訊息提示使用者進出聊天室
* 管理員角色可清除所有訊息

---

## 📂 專案結構（重點檔案）

```
.
├── main.go                  # Gin 入口
├── routes.go               # REST 與 GraphQL 路由註冊
├── auth.go / middleware    # JWT 認證處理
├── expense.go / stocks.go  # REST API 功能模組
├── chat.go                 # WebSocket 聊天功能
├── food.js / GraphQL       # 飲食紀錄（GraphQL）
├── model.go                # 資料模型
├── graph/                  # GraphQL resolvers/generated
├── static/*.html           # 所有前端頁面
├── static/assets/*.js      # 對應前端邏輯
├── docker-compose.yml      # 一鍵啟動前後端容器
├── nginx.conf              # Nginx 反向代理設定
└── README.md               # 說明文件
```

---

## 🧪 測試方式

### ✅ 登入驗證測試

* 可使用 `/user_management.html` 註冊與登入
* 登入後即可瀏覽 `/expense.html`、`/stocks.html`、`/chat.html` 等

### ✅ 模擬操作流程

* 建立帳號 → 管理員審核 → 登入
* 新增記帳、股票、飲食紀錄
* 測試 Excel 上傳與下載
* 進入聊天室、測試訊息發送與清除

---

## 🏁 部署方式

```bash
# 在 EC2 上啟動專案
$ git clone <your-repo>
$ cd your-repo
$ docker-compose up --build -d
```

搭配 Route 53 的 DNS 與 Nginx，即可完成公開部署。

---

## 🙋‍♂️ 作者資訊

* 👨‍💻 作者：王偉任 (Ren Wei Wang)
* 📫 Email：[dokebi871218@gmail.com](mailto:dokebi871218@gmail.com)
* 🌐 GitHub：[github.com/DodoroGit](https://github.com/DodoroGit)

---

這個專案展示了我在後端、前端、資料庫、DevOps 以及整體系統架構上的綜合能力，特別適合作為履歷或面試展示使用。如需協助部署或想了解更多，歡迎與我聯絡。
