document.addEventListener("DOMContentLoaded", () => {
  const token = localStorage.getItem("jwt");
  if (!token) {
    alert("請先登入！");
    window.location.href = "/user_management.html";
    return;
  }

  // 初始載入
  fetchFoodLogs();

  // 表單提交事件
  document.getElementById("food-form").addEventListener("submit", (e) => {
    e.preventDefault();
    addFoodLog();
  });

  // 匯出功能（待實作）
  document.getElementById("export-food-btn").addEventListener("click", () => {
    alert("匯出功能尚未實作");
  });
});

// ⬇ 查詢飲食紀錄
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
        console.error("GraphQL 查詢錯誤", data.errors);
        alert("資料查詢失敗：" + data.errors[0].message);
        return;
      }
      renderTable(data.data.myFoodLogs);
    })
    .catch(err => {
      console.error("查詢失敗", err);
      alert("載入失敗！");
    });
}

// ⬇ 新增飲食紀錄
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
    body: JSON.stringify({
      query: mutation,
      variables: { input }
    })
  })
    .then(res => res.json())
    .then(data => {
      if (data.errors) {
        console.error("GraphQL 新增錯誤", data.errors);
        alert("新增失敗：" + data.errors[0].message);
        return;
      }
      alert("新增成功！");
      document.getElementById("food-form").reset();
      fetchFoodLogs(); // 重新載入
    })
    .catch(err => {
      console.error("新增失敗", err);
      alert("新增失敗！");
    });
}

// ⬇ 渲染表格
function renderTable(logs) {
  const tbody = document.getElementById("food-table-body");
  tbody.innerHTML = "";

  logs.forEach(log => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${log.logged_at || ""}</td>
      <td>${log.name}</td>
      <td>${log.calories ?? ""}</td>
      <td>${log.protein ?? ""}</td>
      <td>${log.fat ?? ""}</td>
      <td>${log.carbs ?? ""}</td>
      <td>${log.quantity}</td>
      <td><!-- 可放刪除按鈕 --></td>
    `;
    tbody.appendChild(tr);
  });
}
