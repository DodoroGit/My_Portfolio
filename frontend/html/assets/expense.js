document.addEventListener("DOMContentLoaded", () => {
    const token = localStorage.getItem("jwt");
    if (!token) {
        alert("請先登入！");
        window.location.href = "/user_management.html";
        return;
    }

    loadExpenses(); // 原本的支出載入

    document.getElementById("expense-form").addEventListener("submit", async (e) => {
        e.preventDefault();
        const category = document.getElementById("category").value;
        const amount = parseFloat(document.getElementById("amount").value);
        const note = document.getElementById("note").value;
        const spentAt = document.getElementById("spentAt").value;

        const res = await fetch("/api/expense/", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            },
            body: JSON.stringify({ category, amount, note, spent_at: spentAt })
        });

        const data = await res.json();
        if (res.ok) {
            alert("新增成功！");
            loadExpenses();
            document.getElementById("expense-form").reset();
        } else {
            alert(data.error);
        }
    });

    // ✅ 把 export 功能也包進 DOMContentLoaded 裡
    document.getElementById("export-btn").addEventListener("click", async () => {
        try {
            const res = await fetch("/api/expense/export", {
                headers: {
                    "Authorization": `Bearer ${token}`
                }
            });

            if (!res.ok) {
                const error = await res.json();
                alert(`匯出失敗：${error.error || res.status}`);
                return;
            }

            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = "expenses.xlsx";
            document.body.appendChild(a);
            a.click();
            a.remove();
        } catch (err) {
            alert("下載錯誤：" + err.message);
        }
    });

    // 🔢 原有計算機邏輯也可考慮包進來以避免綁定失敗
    const amountInput = document.getElementById("amount");
    document.getElementById("calculator").addEventListener("click", (e) => {
        if (e.target.tagName !== "BUTTON") return;

        const value = e.target.textContent;

        if (value === "清除") {
            amountInput.value = "";
        } else if (value === "刪除") {
            amountInput.value = amountInput.value.slice(0, -1);
        } else {
            amountInput.value += value;
        }
    });

    async function loadExpenses() {
        const res = await fetch("/api/expense/", {
            headers: { "Authorization": `Bearer ${token}` }
        });
        const data = await res.json();

        const tbody = document.querySelector("#expense-table tbody");
        tbody.innerHTML = "";

        if (data.expenses) {
            for (const item of data.expenses) {
                const tr = document.createElement("tr");
                tr.innerHTML = `
                    <td>${item.spent_at}</td>
                    <td>${item.category}</td>
                    <td>${item.amount}</td>
                    <td>${item.note || ""}</td>
                `;
                tbody.appendChild(tr);
            }
        }
    }

    document.getElementById("upload-btn").addEventListener("click", async () => {
        const token = localStorage.getItem("jwt");
        const fileInput = document.getElementById("upload-file");
        const file = fileInput.files[0];

        if (!file) {
            alert("請先選擇檔案");
            return;
        }

        const formData = new FormData();
        formData.append("file", file);

        const res = await fetch("/api/expense/upload", {
            method: "POST",
            headers: {
                "Authorization": `Bearer ${token}`
            },
            body: formData
        });

        const data = await res.json();
        if (res.ok) {
            alert("上傳成功！");
            location.reload();
        } else {
            alert(data.error || "上傳失敗");
        }
    });

});
