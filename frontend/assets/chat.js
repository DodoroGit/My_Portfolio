let socket;
let currentUserName = "";
let currentUserId = 0;

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
});

function sendMessage() {
    const input = document.getElementById("message-input");
    if (input.value.trim() === "") return;
    socket.send(JSON.stringify({ content: input.value.trim() }));
    input.value = "";
}

// ⭐ 新增：清除聊天紀錄
async function clearChatHistory() {
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
