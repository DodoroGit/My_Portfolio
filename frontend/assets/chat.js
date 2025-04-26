let socket;
let currentUserName = "";
let currentUserId = 0;
let currentUserRole = ""; // ⭐ 新增：角色

document.addEventListener("DOMContentLoaded", async function() {
    const token = localStorage.getItem("jwt");
    if (!token) {
        alert("請先登入！");
        window.location.href = "/usermanagement";
        return;
    }

    const res = await fetch(`${window.location.origin}/api/user/profile`, {
        headers: { "Authorization": `Bearer ${token}` }
    });
    const data = await res.json();
    if (data.user) {
        currentUserName = data.user.name;
        currentUserId = data.user.id;
        currentUserRole = data.user.role; // ⭐ 把role記下來
    } else {
        alert("取得使用者資訊失敗！");
        return;
    }

    const protocol = location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//${location.host}/ws/chat?token=${token}`;
    socket = new WebSocket(wsUrl);

    socket.onmessage = function(event) {
        const msg = JSON.parse(event.data);
        const chatBox = document.getElementById("chat-box");

        if (msg.type === "system") {
            const systemDiv = document.createElement("div");
            systemDiv.classList.add("system-message");
            systemDiv.textContent = `[系統訊息] ${msg.content}`;
            chatBox.appendChild(systemDiv);
        } else if (msg.type === "message") {
            const messageDiv = document.createElement("div");
            messageDiv.classList.add("message");

            if (msg.user_id === currentUserId) {
                messageDiv.classList.add("right");
            } else {
                messageDiv.classList.add("left");
            }

            const timeStr = new Date(msg.timestamp).toLocaleTimeString();
            messageDiv.innerHTML = `
                <div class="message-author">${msg.user_name || '未知使用者'}</div>
                <div class="message-content">${msg.content}</div>
                <div class="message-time">${timeStr}</div>
            `;
            chatBox.appendChild(messageDiv);
        }

        chatBox.scrollTop = chatBox.scrollHeight;
    };

    socket.onclose = function(event) {
        alert("連線中斷，請重新整理頁面！");
    };

    const input = document.getElementById("message-input");
    input.addEventListener("keydown", function(event) {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            sendMessage();
        }
    });

    // ⭐⭐ 新增：如果是admin，動態產生「清除聊天紀錄」按鈕
    if (currentUserRole === "admin") {
        const clearButton = document.createElement("button");
        clearButton.textContent = "清除聊天紀錄";
        clearButton.classList.add("chat-clear-btn");
        clearButton.onclick = clearChatHistory;
        document.querySelector(".chat-input-container").appendChild(clearButton);
    }
});

function sendMessage() {
    const input = document.getElementById("message-input");
    if (input.value.trim() === "") return;
    socket.send(JSON.stringify({ content: input.value.trim() }));
    input.value = "";
}

async function clearChatHistory() {
    if (!confirm("⚠️ 確定要清除所有聊天紀錄嗎？此操作無法復原！")) {
        return; // 使用者按了取消
    }

    const token = localStorage.getItem("jwt");
    const res = await fetch("/api/chat/clear", {
        method: "POST",
        headers: { "Authorization": `Bearer ${token}` }
    });
    if (res.ok) {
        document.getElementById("chat-box").innerHTML = "";
        alert("聊天紀錄已清空！");
    } else {
        alert("清除失敗！");
    }
}
