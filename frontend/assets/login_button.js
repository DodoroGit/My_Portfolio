// ✅ 檢查 JWT Token，控制按鈕顯示
function checkAuthStatus() {
    const token = localStorage.getItem("jwt");

    if (token) {
        document.getElementById("profile-btn").style.display = "inline-block";
        document.getElementById("logout-btn").style.display = "inline-block";
        document.getElementById("login-btn").style.display = "none";

        // ✅ 如果用戶已登入但在 /login 頁面，應該導向 dashboard
        if (window.location.pathname === "/auth") {
            window.location.href = "/usermanagementdashboard";
        }
    } else {
        document.getElementById("profile-btn").style.display = "none";
        document.getElementById("logout-btn").style.display = "none";
        document.getElementById("login-btn").style.display = "inline-block";

        // ✅ 如果用戶未登入但在 /dashboard，應該導向 /login
        if (window.location.pathname === "/usermanagementdashboard") {
            alert("請先登入！");
            window.location.href = "/usermanagement";
        }
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
    window.location.href = "/usermanagementdashboard";
}

function goToLogin() {
    window.location.href = "/usermanagement";
}