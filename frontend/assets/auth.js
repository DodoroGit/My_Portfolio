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
        if (window.location.pathname === "/login") {
            window.location.href = "/dashboard";
        }
    } else {
        document.getElementById("profile-btn").style.display = "none";
        document.getElementById("logout-btn").style.display = "none";
        document.getElementById("login-btn").style.display = "inline-block";

        // ✅ 如果用戶未登入但在 /dashboard，應該導向 /login
        if (window.location.pathname === "/dashboard") {
            alert("請先登入！");
            window.location.href = "/login";
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
    window.location.href = "/login";
}

// ✅ 連結按鈕點擊行為
function goToProfile() {
    window.location.href = "/dashboard";
}

function goToLogin() {
    window.location.href = "/login";
}
