const API_BASE = window.location.origin + "/api/auth";

// 當網頁載入時，檢查使用者是否已登入
document.addEventListener("DOMContentLoaded", function () {
    checkAuthStatus();
});

// ✅ 檢查 JWT Token，控制按鈕顯示
function checkAuthStatus() {
    const token = localStorage.getItem("jwt");

    if (token) {
        document.getElementById("profile-btn").style.display = "inline-block";
        document.getElementById("logout-btn").style.display = "inline-block";
        document.getElementById("login-btn").style.display = "none";

        // ✅ 如果用戶已登入但在 /login 頁面，應該導向 dashboard
        if (window.location.pathname === "/auth") {
            window.location.href = "/dashboard";
        }
    } else {
        document.getElementById("profile-btn").style.display = "none";
        document.getElementById("logout-btn").style.display = "none";
        document.getElementById("login-btn").style.display = "inline-block";

        // ✅ 如果用戶未登入但在 /dashboard，應該導向 /login
        if (window.location.pathname === "/dashboard") {
            alert("請先登入！");
            window.location.href = "/auth";
        }
    }
}

// ✅ 切換登入 & 註冊表單
function toggleForm() {
    document.querySelectorAll(".form-box").forEach(box => {
        box.style.display = box.style.display === "none" ? "block" : "none";
    });
}

// ✅ 註冊新使用者
async function register() {
    const name = document.getElementById("register-name").value;
    const email = document.getElementById("register-email").value;
    const password = document.getElementById("register-password").value;

    const res = await fetch(`${API_BASE}/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, password }),
    });

    const data = await res.json();
    if (res.ok) {
        alert("註冊成功！請登入");
        toggleForm(); // 切換回登入頁面
    } else {
        alert(data.error);
    }
}

// ✅ 登入並存儲 JWT Token
async function login() {
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    const res = await fetch(`${API_BASE}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    });

    const data = await res.json();
    if (data.token) {
        localStorage.setItem("jwt", data.token);
        alert("登入成功！");
        window.location.href = "/dashboard"; // 導向個人頁面
    } else {
        alert(data.error);
    }
}

// ✅ 取得個人資訊（可選）
async function getProfile() {
    const token = localStorage.getItem("jwt");
    if (!token) return;

    const res = await fetch(`${window.location.origin}/api/user/profile`, {
        method: "GET",
        headers: { "Authorization": `Bearer ${token}` }
    });

    if (res.status === 401) {
        alert("登入已過期，請重新登入！");
        logout();
    }

    const data = await res.json();
    if (data.user) {
        document.getElementById("user-name").innerText = data.user.name;
        document.getElementById("user-email").innerText = data.user.email;
        document.getElementById("user-role").innerText = data.user.role;
    }
}

// ✅ 登出功能
function logout() {
    localStorage.removeItem("jwt");
    alert("登出成功！");
    window.location.href = "/auth";
}

// ✅ 連結按鈕點擊行為
function goToProfile() {
    window.location.href = "/dashboard";
}

function goToLogin() {
    window.location.href = "/auth";
}

document.addEventListener("DOMContentLoaded", function () {
    // 創建漢堡選單按鈕
    const menuToggle = document.createElement("button");
    menuToggle.classList.add("menu-toggle");
    menuToggle.innerHTML = "☰"; // 漢堡圖示
    document.querySelector("header").appendChild(menuToggle);

    // 取得導覽列
    const navMenu = document.querySelector("nav ul");

    // 監聽按鈕點擊，切換選單顯示
    menuToggle.addEventListener("click", function () {
        if (navMenu.style.display === "flex") {
            navMenu.style.display = "none"; // 收起
        } else {
            navMenu.style.display = "flex"; // 展開
            navMenu.style.flexDirection = "column"; // 手機版直列顯示
            navMenu.style.position = "absolute";
            navMenu.style.top = "60px";
            navMenu.style.right = "20px";
            navMenu.style.background = "rgba(0, 0, 0, 0.9)";
            navMenu.style.padding = "10px";
            navMenu.style.borderRadius = "10px";
            navMenu.style.zIndex = "1000";
        }
    });

    // 監聽視窗變化，確保桌機恢復預設顯示
    window.addEventListener("resize", function () {
        if (window.innerWidth > 768) {
            navMenu.style.display = "flex"; // 桌機版顯示
            navMenu.style.flexDirection = "row"; // 恢復水平排列
            navMenu.style.position = ""; // 移除手機版的額外樣式
            navMenu.style.background = "";
            navMenu.style.padding = "";
            navMenu.style.borderRadius = "";
        } else {
            navMenu.style.display = "none"; // 小螢幕時隱藏，需點擊才展開
        }
    });

    // 初始檢查 RWD 狀態
    if (window.innerWidth <= 768) {
        navMenu.style.display = "none"; // 手機版一開始隱藏選單
    }
});
