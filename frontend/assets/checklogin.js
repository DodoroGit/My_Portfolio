document.addEventListener("DOMContentLoaded", function () {
    const token = localStorage.getItem("jwt");
    console.log("JWT Token:", token); // 🛠 確保有取得 Token

    if (!token) {
        alert("請先登入！\n您將被導向至登入頁面。");
        window.location.href = "/auth";
    } else {
        alert("登入成功，正在載入頁面...");
        document.documentElement.style.display = "block"; // 🛠 顯示頁面
    }
});
