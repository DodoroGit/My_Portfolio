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

    document.getElementById("export-tx-btn").addEventListener("click", () => {
        fetch("/api/stocks/transactions/export", {
            headers: { "Authorization": `Bearer ${token}` }
        })
        .then(res => res.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = "transactions.xlsx";
            a.click();
            window.URL.revokeObjectURL(url);
        });
    });

    
    connectWebSocket();// 連線 WebSocket
    loadTransactions(); // 加這行載入交易紀錄
    loadProfitSummary(); // 取得總損益
});

// 渲染股票表格（初始用）
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
            <td>
                <button onclick="sellStockPrompt('${stock.symbol}', ${stock.shares})">賣出</button>
                <button onclick="receiveDividendPrompt('${stock.symbol}')">💰領股息</button>
                <button onclick="deleteStock(${stock.id})">🗑️刪除</button>
            </td>
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

function sellStockPrompt(symbol, currentShares) {
    const sellShares = prompt(`請輸入要賣出的股數（最多 ${currentShares} 股）：`);
    const sellPrice = prompt(`請輸入每股賣出價格：`);
    const note = prompt("備註（可留空）：");

    if (!sellShares || !sellPrice) return alert("請輸入完整資訊");

    const token = localStorage.getItem("jwt");
    fetch("/api/stocks/sell", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            symbol,
            shares: parseInt(sellShares),
            sell_price: parseFloat(sellPrice),
            note: note || ""
        })
    })
    .then(res => res.json())
    .then(data => {
        if (data.error) return alert(data.error);
        alert(`賣出成功，損益：${data.realized_profit}`);
        location.reload();
    });
}

function loadTransactions() {
    const token = localStorage.getItem("jwt");
    fetch("/api/stocks/transactions", {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => renderTransactions(data.transactions));
}

function renderTransactions(transactions) {
    const container = document.getElementById("tx-records");
    container.innerHTML = "<h2>交易紀錄</h2>";
    const table = document.createElement("table");
    table.border = "1";
    table.width = "100%";
    table.innerHTML = `
        <thead>
            <tr>
                <th>代碼</th><th>股數</th><th>均價</th><th>賣價</th><th>損益</th><th>備註</th><th>時間</th><th>操作</th>
            </tr>
        </thead>
        <tbody>` +
        transactions.map(tx => `
            <tr>
                <td>${tx.symbol}</td>
                <td>${tx.shares}</td>
                <td>${tx.avg_price}</td>
                <td>${tx.sell_price}</td>
                <td class="${tx.profit >= 0 ? 'profit-positive' : 'profit-negative'}">${tx.profit}</td>
                <td class="${tx.note.includes('股息') ? 'profit-dividend' : ''}">${tx.note || ""}</td>
                <td>${new Date(tx.created_at).toLocaleString()}</td>
                <td><button onclick="deleteTransaction(${tx.id})">🗑️</button></td>
            </tr>
        `).join("") +
        "</tbody>";
    container.appendChild(table);
}



// 點擊外部區域或關閉按鈕關掉 modal
window.addEventListener("click", (e) => {
    const modal = document.getElementById("chart-modal");
    if (e.target === modal) modal.style.display = "none";
});

function loadProfitSummary() {
    const token = localStorage.getItem("jwt");
    fetch("/api/stocks/summary", {
        headers: { "Authorization": `Bearer ${token}` }
    })
    .then(res => res.json())
    .then(data => {
        const div = document.getElementById("profit-summary");
        const unrealized = data.unrealized_profit.toFixed(2);
        const realized = data.realized_profit.toFixed(2);
        const total = data.total_profit.toFixed(2);
        const totalClass = data.total_profit >= 0 ? "profit-positive" : "profit-negative";

        div.innerHTML = `
            <div>
                🧾 <strong>總損益：</strong>
                <span class="${totalClass}" style="font-size: 1.4rem;">${total}</span><br>
                <span style="font-size: 14px; color: #555;">
                    （未實現：<strong>${unrealized}</strong>，已實現：<strong>${realized}</strong>）
                </span>
            </div>

            <div style="margin-top: 1rem; text-align: left; font-size: 16px; background: #f8f9fa; padding: 1rem; border-radius: 10px; box-shadow: 0 2px 6px rgba(0,0,0,0.05); line-height: 1.8; font-family: '標楷體', Cambria; font-weight: normal;">
                <h3 style="font-size: 15px; margin: 0 0 0.5rem 0; font-weight: normal;">💡 計算公式說明：</h3>
                <ul style="padding-left: 1rem; margin: 0;">
                    <li>每筆損益計算方式如下：</li>
                    <li>(持股數*即時價格-手續費-證交稅)-(持股數*均價+手續費)</li>
                    <li>手續費依據 0.001425 * 0.35 計算(四捨五入,最低為1元)</li>
                    <li>證交稅為 0.3%，且無條件捨去</li>
                </ul>
            </div>

            <p style="color: red; font-size: 13px; margin-top: 8px; font-family: '標楷體', Cambria; font-weight: normal;">
                ⚠️ 最終數字可能與券商 App 有誤差，僅供參考，請以官方資訊為準。
            </p>
        `;
    });
}



function receiveDividendPrompt(symbol) {
    const amount = prompt(`請輸入「${symbol}」股息金額：`);
    if (!amount || isNaN(amount) || parseFloat(amount) <= 0) return alert("請輸入有效金額");

    const note = prompt("備註（可選）：") || "";
    const token = localStorage.getItem("jwt");

    fetch("/api/stocks/dividend", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            symbol,
            amount: parseFloat(amount),
            note
        })
    })
    .then(res => res.json())
    .then(data => {
        if (data.error) return alert(data.error);
        alert("股息已記錄！");
        location.reload();
    });
}

function deleteTransaction(id) {
    if (!confirm("確定要刪除這筆交易紀錄？")) return;
    fetch(`/api/stocks/transactions/${id}`, {
        method: "DELETE",
        headers: { "Authorization": `Bearer ${localStorage.getItem("jwt")}` }
    })
    .then(res => res.json())
    .then(data => {
        alert(data.message || "已刪除");
        location.reload();
    });
}
