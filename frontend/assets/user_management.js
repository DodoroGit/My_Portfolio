const API_USER = window.location.origin + "/api/user/profile";

// 取得 JWT Token
const token = localStorage.getItem("jwt");
if (!token) {
    alert("請先登入！");
    window.location.href = "/usermanagement";
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
        window.location.href = "/usermanagement";
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
    window.location.href = "/usermanagement";
}

// 檢查是否為管理員
function checkAdminRole() {
    const userRole = document.getElementById('user-role').textContent;
    const adminSection = document.getElementById('admin-section');
    if (userRole === 'admin') {
        adminSection.style.display = 'block';
        loadPendingUsers();
        loadAllUsers();
    }
}

// 載入待審核用戶
async function loadPendingUsers() {
    try {
        const response = await fetch('/api/admin/pending-users', {
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('jwt')
            }
        });
        const data = await response.json();
        
        const pendingUsersList = document.getElementById('pending-users-list');
        pendingUsersList.innerHTML = '';
        
        data.pending_users.forEach(user => {
            const userDiv = document.createElement('div');
            userDiv.className = 'user-item';
            userDiv.innerHTML = `
                <p><strong>名稱：</strong>${user.name}</p>
                <p><strong>Email：</strong>${user.email}</p>
                <p><strong>註冊時間：</strong>${new Date(user.created_at).toLocaleString()}</p>
                <button onclick="approveUser(${user.id}, 'approve')">批准</button>
                <button onclick="approveUser(${user.id}, 'reject')">拒絕</button>
            `;
            pendingUsersList.appendChild(userDiv);
        });
    } catch (error) {
        console.error('載入待審核用戶失敗：', error);
    }
}

// 審核用戶
async function approveUser(userId, action) {
    try {
        const response = await fetch('/api/admin/approve-user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            },
            body: JSON.stringify({
                user_id: userId,
                action: action
            })
        });
        
        const data = await response.json();
        alert(data.message);
        loadPendingUsers(); // 重新載入待審核用戶列表
        loadAllUsers(); // 重新載入所有用戶列表
    } catch (error) {
        console.error('審核用戶失敗：', error);
        alert('審核失敗，請稍後再試');
    }
}

// 載入所有用戶
async function loadAllUsers() {
    try {
        const response = await fetch('/api/admin/users', {
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('jwt')
            }
        });
        const data = await response.json();
        
        const allUsersList = document.getElementById('all-users-list');
        allUsersList.innerHTML = '';
        
        data.users.forEach(user => {
            const userDiv = document.createElement('div');
            userDiv.className = 'user-item';
            userDiv.innerHTML = `
                <p><strong>名稱：</strong>${user.name}</p>
                <p><strong>Email：</strong>${user.email}</p>
                <p><strong>角色：</strong>${user.role}</p>
                <p><strong>狀態：</strong>${user.status}</p>
                <p><strong>註冊時間：</strong>${new Date(user.created_at).toLocaleString()}</p>
            `;
            allUsersList.appendChild(userDiv);
        });
    } catch (error) {
        console.error('載入所有用戶失敗：', error);
    }
}

// 頁面載入時執行
document.addEventListener('DOMContentLoaded', () => {
    checkAdminRole();
});

