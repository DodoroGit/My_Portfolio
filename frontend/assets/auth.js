const API_BASE = "http://localhost:8080/api/auth";

// 切換登入 & 註冊表單
function toggleForm() {
    document.querySelectorAll(".form-box").forEach(box => box.style.display = 
        box.style.display === "none" ? "block" : "none");
}

// 註冊新使用者
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
    alert(data.message || data.error);
}

// 登入並存儲 JWT Token
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
