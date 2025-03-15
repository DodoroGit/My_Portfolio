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
    window.location.href = "/user_management";
}