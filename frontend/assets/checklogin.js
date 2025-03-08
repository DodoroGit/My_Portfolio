// 取得 JWT Token
const token = localStorage.getItem("jwt");
if (!token) {
    alert("請先登入！");
    window.location.href = "/auth";
}

document.addEventListener("DOMContentLoaded", function () {
    const token = localStorage.getItem("jwt");
    if (!token) {
        window.location.href = "/auth";
    } else {
        document.documentElement.style.display = ""; // 顯示頁面
    }
});