document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("jwt");

    // 加載股票清單
    fetch("/api/stocks/", {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => renderTable(data.stocks));

    // 顯示總資產概覽
    fetch("/api/stocks/summary", {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(summary => {
        const totalBox = document.getElementById("summary-box");
        if (summary.total_cost !== undefined) {
            totalBox.innerHTML = `
                <p>💰 持股成本總額：${summary.total_cost}</p>
                <p>📈 現值總額：${summary.total_value}</p>
                <p>📊 未實現損益：<strong style="color:${parseFloat(summary.unrealized_pnl) >= 0 ? 'red' : 'green'}">${summary.unrealized_pnl}</strong></p>
            `;
        } else {
            totalBox.innerHTML = `<p>無法取得總資產資料</p>`;
        }
    });

    // 提交新增/更新股票表單
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

    // 提交賣出股票表單
    document.getElementById("sell-form").addEventListener("submit", (e) => {
        e.preventDefault();
        const symbol = document.getElementById("sell-symbol").value.trim();
        const shares = parseInt(document.getElementById("sell-shares").value);
        const sell_price = parseFloat(document.getElementById("sell-price").value);
        const note = document.getElementById("sell-note").value.trim();

        if (!symbol || shares <= 0 || sell_price <= 0) {
            return alert("請填寫完整賣出資料");
        }

        fetch("/api/stocks/sell", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ symbol, shares, sell_price, note })
        })
        .then(res => res.json())
        .then(data => {
            alert(data.message + (data.realized_profit ? `\n本次損益：${data.realized_profit}` : ""));
            location.reload();
        });
    });

    // 匯出 Excel
    document.getElementById("export-btn").addEventListener("click", () => {
        fetch("/api/stocks/export", {
            headers: { "Authorization": `Bearer ${token}` }
        })
        .then(res => res.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = "stocks.xlsx";
            a.click();
            window.URL.revokeObjectURL(url);
        });
    });

    // WebSocket 連線
    connectWebSocket();
});

function renderTable(stocks) {
    const tbody = document.getElementById("stock-table-body");
    tbody.innerHTML = "";
    stocks.forEach(stock => {
        const row = document.createElement("tr");
        row.id = `stock-row-${stock.symbol}`;
        row.innerHTML = `
            <td><button onclick="viewChart('${stock.symbol}')">${stock.symbol}</button></td>
            <td>${stock.shares}</td>
            <td id="avg-${stock.symbol}">${stock.avg_price !== undefined ? stock.avg_price.toFixed(2) : '-'}</td>
            <td id="price-${stock.symbol}">-</td>
            <td id="profit-${stock.symbol}">-</td>
            <td><button onclick="deleteStock('${stock.id}')">刪除</button></td>
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

    const profitClass = data.profit >= 0 ? "profit-positive" : "profit-negative";
    const avgPriceText = data.avg_price !== undefined ? data.avg_price.toFixed(2) : "-";
    const priceText = data.price !== undefined ? data.price.toFixed(2) : "-";
    const profitText = data.profit !== undefined ? data.profit.toFixed(2) : "-";

    if (!row) {
        row = document.createElement("tr");
        row.id = rowId;
        row.innerHTML = `
            <td><button onclick="viewChart('${data.symbol}')">${data.symbol}</button></td>
            <td>${data.shares}</td>
            <td id="avg-${data.symbol}">${avgPriceText}</td>
            <td id="price-${data.symbol}">${priceText}</td>
            <td id="profit-${data.symbol}" class="${profitClass}">${profitText}</td>
        `;
        document.getElementById("stock-table-body").appendChild(row);
    } else {
        document.getElementById(`avg-${data.symbol}`).textContent = avgPriceText;
        document.getElementById(`price-${data.symbol}`).textContent = priceText;
        const profitCell = document.getElementById(`profit-${data.symbol}`);
        profitCell.textContent = profitText;
        profitCell.className = profitClass;
    }
}

function viewChart(symbol) {
    const modal = document.getElementById("chart-modal");
    const canvas = document.getElementById("stock-chart");
    const ctx = canvas.getContext("2d");
    const token = localStorage.getItem("jwt");

    fetch(`/api/stocks/history/${symbol}`, {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => {
        const labels = data.map(p => p.date);
        const prices = data.map(p => p.price);

        if (window.myChart) window.myChart.destroy();
        window.myChart = new Chart(ctx, {
            type: "line",
            data: {
                labels,
                datasets: [{
                    label: `${symbol} 股價趨勢`,
                    data: prices,
                    fill: false,
                    borderWidth: 2
                }]
            },
            options: {
                responsive: true,
                scales: {
                    x: { title: { display: true, text: "日期" } },
                    y: { title: { display: true, text: "價格" } }
                }
            }
        });
        modal.style.display = "block";
    });
}

window.addEventListener("click", (e) => {
    const modal = document.getElementById("chart-modal");
    if (e.target === modal) modal.style.display = "none";
});
