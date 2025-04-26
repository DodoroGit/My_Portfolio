let socket;

document.addEventListener("DOMContentLoaded", function() {
    const token = localStorage.getItem("jwt");
    if (!token) {
        alert("請先登入！");
        window.location.href = "/usermanagement";
        return;
    }

    // 換成 wss://yourdomain.com/ws/chat (如果有https)
    const wsUrl = (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/ws/chat";
    socket = new WebSocket(wsUrl);

    socket.onmessage = function(event) {
        const msg = JSON.parse(event.data);
        const timeStr = new Date(msg.timestamp).toLocaleTimeString();
        const chatBox = document.getElementById("chat-box");
        const p = document.createElement("p");
        p.textContent = `[${timeStr}] ${msg.user_name}: ${msg.content}`;
        chatBox.appendChild(p);
        chatBox.scrollTop = chatBox.scrollHeight; // 自動捲到最底
    };

    socket.onclose = function(event) {
        alert("連線中斷，請重新整理頁面！");
    };
});

function sendMessage() {
    const input = document.getElementById("message-input");
    if (input.value.trim() === "") return;
    socket.send(JSON.stringify({ content: input.value.trim() }));
    input.value = "";
}
