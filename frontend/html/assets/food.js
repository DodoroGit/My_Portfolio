document.addEventListener("DOMContentLoaded", () => {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("請先登入！");
    window.location.href = "/user_management.html";
    return;
  }

  fetchFoodLogs();

  document.getElementById("food-form").addEventListener("submit", (e) => {
    e.preventDefault();
    addFoodLog();
  });

  document.getElementById("export-food-btn").addEventListener("click", () => {
    alert("匯出功能尚未實作");
  });
});

function fetchFoodLogs() {
  const query = `
    query {
      myFoodLogs {
        id
        name
        calories
        protein
        fat
        carbs
        quantity
        loggedAt
      }
    }
  `;

  fetch("/graphql", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("jwt")}`
    },
    body: JSON.stringify({ query })
  })
    .then(res => res.json())
    .then(data => {
      if (data.errors) {
        alert("查詢錯誤：" + data.errors[0].message);
        return;
      }
      const logs = data.data.myFoodLogs;
      renderTable(logs);
      renderDateOptions(logs);
    })
    .catch(err => {
      console.error("查詢失敗", err);
      alert("載入失敗！");
    });
}

function addFoodLog() {
  const input = {
    name: document.getElementById("name").value.trim(),
    calories: parseFloat(document.getElementById("calories").value || 0),
    protein: parseFloat(document.getElementById("protein").value || 0),
    fat: parseFloat(document.getElementById("fat").value || 0),
    carbs: parseFloat(document.getElementById("carbs").value || 0),
    quantity: document.getElementById("quantity").value.trim(),
    loggedAt: document.getElementById("loggedAt").value
  };

  const mutation = `
    mutation ($input: FoodLogInput!) {
      addFoodLog(input: $input) {
        id
        name
      }
    }
  `;

  fetch("/graphql", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("jwt")}`
    },
    body: JSON.stringify({ query: mutation, variables: { input } })
  })
    .then(res => res.json())
    .then(data => {
      if (data.errors) {
        alert("新增失敗：" + data.errors[0].message);
        return;
      }
      alert("新增成功！");
      document.getElementById("food-form").reset();
      fetchFoodLogs();
    })
    .catch(err => {
      console.error("新增失敗", err);
      alert("新增失敗！");
    });
}

function renderTable(logs) {
  const tbody = document.getElementById("food-table-body");
  tbody.innerHTML = "";

  logs.forEach(log => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${log.loggedAt || ""}</td>
      <td>${log.name}</td>
      <td>${log.calories ?? ""}</td>
      <td>${log.protein ?? ""}</td>
      <td>${log.fat ?? ""}</td>
      <td>${log.carbs ?? ""}</td>
      <td>${log.quantity ?? ""}</td>
      <td>
        <button onclick="editFoodLog(${log.id})">✏️</button>
        <button onclick="deleteFoodLog(${log.id})">🗑️</button>
      </td>
    `;
    tbody.appendChild(tr);
  });
}

function renderDateOptions(logs) {
  const select = document.getElementById("date-filter");
  const dates = [...new Set(logs.map(log => log.loggedAt))].sort().reverse();
  select.innerHTML = `<option value="">請選擇日期</option>` +
    dates.map(date => `<option value="${date}">${date}</option>`).join("");

  select.onchange = () => {
    const selected = select.value;
    const filtered = logs.filter(log => log.loggedAt === selected);
    renderTable(filtered);
    const totalCalories = filtered.reduce((sum, f) => sum + (f.calories || 0), 0);
    const totalProtein = filtered.reduce((sum, f) => sum + (f.protein || 0), 0);
    document.getElementById("summary-text").textContent = 
      `總熱量：${totalCalories.toFixed(1)} kcal，總蛋白質：${totalProtein.toFixed(1)} g`;
  };
}

function deleteFoodLog(id) {
  if (!confirm("確定要刪除這筆紀錄嗎？")) return;

  const mutation = `
    mutation {
      deleteFoodLog(id: ${id})
    }
  `;

  fetch("/graphql", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("jwt")}`
    },
    body: JSON.stringify({ query: mutation })
  })
    .then(res => res.json())
    .then(() => {
      alert("已刪除！");
      fetchFoodLogs();
    })
    .catch(err => console.error("刪除失敗", err));
}

function editFoodLog(id) {
  const row = [...document.querySelectorAll("#food-table-body tr")].find(tr => 
    tr.querySelector("button")?.getAttribute("onclick")?.includes(`editFoodLog(${id})`)
  );

  if (!row) return;

  const cells = row.querySelectorAll("td");
  const originalData = Array.from(cells).map(td => td.textContent);

  // 替換成輸入欄位
  const fields = ["date", "name", "cal", "protein", "fat", "carbs", "qty"];
  const inputs = fields.map((_, i) => {
    return `<input type="${i === 0 ? 'date' : 'text'}" value="${originalData[i] || ''}" style="width:80px;" />`;
  });

  row.innerHTML = `
    <td>${inputs[0]}</td>
    <td>${inputs[1]}</td>
    <td>${inputs[2]}</td>
    <td>${inputs[3]}</td>
    <td>${inputs[4]}</td>
    <td>${inputs[5]}</td>
    <td>${inputs[6]}</td>
    <td>
      <button onclick="saveFoodLog(${id}, this)">💾</button>
      <button onclick="cancelEdit()">❌</button>
    </td>
  `;
}

function saveFoodLog(id, btn) {
  const row = btn.closest("tr");
  const inputs = row.querySelectorAll("input");

  const input = {
    loggedAt: inputs[0].value,
    name: inputs[1].value.trim(),
    calories: parseFloat(inputs[2].value || 0),
    protein: parseFloat(inputs[3].value || 0),
    fat: parseFloat(inputs[4].value || 0),
    carbs: parseFloat(inputs[5].value || 0),
    quantity: inputs[6].value.trim(),
  };

  const mutation = `
    mutation ($id: Int!, $input: FoodLogInput!) {
      updateFoodLog(id: $id, input: $input) {
        id
      }
    }
  `;

  fetch("/graphql", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${localStorage.getItem("jwt")}`
    },
    body: JSON.stringify({ query: mutation, variables: { id, input } })
  })
    .then(res => res.json())
    .then(data => {
      if (data.errors) {
        alert("更新失敗：" + data.errors[0].message);
        return;
      }
      alert("更新成功！");
      fetchFoodLogs();
    })
    .catch(err => {
      console.error("更新失敗", err);
      alert("更新失敗！");
    });
}

function cancelEdit() {
  fetchFoodLogs();
}
