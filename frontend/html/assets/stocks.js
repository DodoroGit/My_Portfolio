document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("jwt");

    // 加載股票清單
    fetch("/api/stocks/", {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => renderTable(data.stocks));

    // 提交表單
    document.getElementById("stock-form").addEventListener("submit", (e) => {
        e.preventDefault();
        const symbol = document.getElementById("symbol").value.trim();
        const shares = parseInt(document.getElementById("shares").value);
        const avg_price = parseFloat(document.getElementById("avgPrice").value || "0");
        if (!symbol || !shares) return alert("請填寫完整欄位");

        fetch("/api/stocks/", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ symbol, shares, avg_price })
        })
        .then(res => res.json())
        .then(data => {
            alert(data.message || "操作完成");
            location.reload();
        });
    });

    // 連線 WebSocket
    connectWebSocket();
});

// 渲染股票表格（初始用）
function renderTable(stocks) {
    const tbody = document.getElementById("stock-table-body");
    tbody.innerHTML = "";
    stocks.forEach(stock => {
        const row = document.createElement("tr");
        row.id = `stock-row-${stock.symbol}`;
        row.innerHTML = `
            <td>${stock.symbol}</td>
            <td>${stock.shares}</td>
            <td id="price-${stock.symbol}">-</td>
            <td id="profit-${stock.symbol}">-</td>
        `;
        tbody.appendChild(row);
    });
}

function deleteStock(id) {
    const token = localStorage.getItem("jwt");
    fetch(`/api/stocks/${id}`, {
        method: "DELETE",
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => {
        alert(data.message || "已刪除");
        location.reload();
    });
}

let socket;

function connectWebSocket() {
    const token = localStorage.getItem("jwt");
    socket = new WebSocket(`wss://${window.location.host}/ws/stocks/?token=${token}`);

    socket.onopen = () => {
        console.log("✅ WebSocket 已連線");
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        updateStockRow(data);
    };

    socket.onclose = () => {
        console.log("❌ WebSocket 斷線，5秒後重連...");
        setTimeout(connectWebSocket, 5000);
    };
}

function updateStockRow(data) {
    const rowId = `stock-row-${data.symbol}`;
    let row = document.getElementById(rowId);
    if (!row) {
        row = document.createElement("tr");
        row.id = rowId;
        row.innerHTML = `
            <td>${data.symbol}</td>
            <td>${data.shares}</td>
            <td id="price-${data.symbol}">${data.price}</td>
            <td id="profit-${data.symbol}">${data.profit.toFixed(2)}</td>
        `;
        document.getElementById("stock-table-body").appendChild(row);
    } else {
        document.getElementById(`price-${data.symbol}`).textContent = data.price;
        document.getElementById(`profit-${data.symbol}`).textContent = data.profit.toFixed(2);
    }
}
