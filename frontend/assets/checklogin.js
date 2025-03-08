document.addEventListener("DOMContentLoaded", function () {
    const token = localStorage.getItem("jwt");
    console.log("JWT Token:", token); // 🛠 除錯輸出

    if (!token) {
        alert("請先登入！\n您將被導向至登入頁面。");
        console.log("未登入，重導向至登入頁面");
        window.location.href = "/auth";
    } else {
        alert("登入成功，正在載入頁面...");
        console.log("登入成功，嘗試顯示頁面");

        // 🛠 方法 1：直接修改 `display`
        document.documentElement.style.display = "block"; 

        // 🛠 方法 2：確保 CSS 的 `class` 變更
        document.documentElement.classList.add("show");
        
        console.log("已設置 display:block;");
    }
});
