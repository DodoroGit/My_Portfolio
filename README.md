# My Portfolio

這是一個個人作品集網站，採用 **Golang Gin** 作為後端，前端使用 **純 HTML、CSS、JavaScript**。

## 專案結構

```bash
My_Portfolio/
│── backend/                 # 後端 Golang Gin
│   │── main.go              # 入口點
│   │── handler/             # 處理 API 邏輯
│   │── route/               # 設定 API 路由
│   └── static/              # 提供前端靜態文件
│
└── frontend/                # 前端 (純 HTML, CSS, JS)
    │── index.html           # 首頁
    │── about.html           # 關於我
    │── projects.html        # 作品展示
    │── skills.html          # 技能介紹
    │── contact.html         # 聯絡方式
    │── assets/              # 靜態資源 (圖片、CSS、JS)
    │   │── styles.css       # 全站 CSS
    │   └── script.js        # 全站 JS
```
