const API_USER = "http://localhost:8080/api/user/profile";

// 取得 JWT Token
const token = localStorage.getItem("jwt");
if (!token) {
    alert("請先登入！");
    window.location.href = "/auth";
}

// 取得個人資訊
async function getProfile() {
    const res = await fetch(API_USER, {
        method: "GET",
        headers: { "Authorization": `Bearer ${token}` }
    });

    if (res.status === 401) {
        alert("登入已過期，請重新登入！");
        localStorage.removeItem("jwt");
        window.location.href = "/auth.html";
    }

    const data = await res.json();
    if (data.user) {
        document.getElementById("user-name").innerText = data.user.name;
        document.getElementById("user-email").innerText = data.user.email;
        document.getElementById("user-role").innerText = data.user.role;
    } else {
        alert("無法取得個人資料");
    }
}

// 登出
function logout() {
    localStorage.removeItem("jwt");
    alert("登出成功！");
    window.location.href = "/auth";
}

getProfile();
