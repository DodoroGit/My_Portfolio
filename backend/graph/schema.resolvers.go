package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/DodoroGit/My_Portfolio/backend/database"
	"github.com/DodoroGit/My_Portfolio/backend/graph/generated"
	"github.com/DodoroGit/My_Portfolio/backend/graph/model"
)

func (r *Resolver) MyFoodLogs(ctx context.Context) ([]*model.FoodLog, error) {
	userID := ctx.Value("user_id").(int) // 你需要從 middleware 帶進來
	rows, err := database.DB.Query(`SELECT id, name, calories, protein, fat, carbs, quantity, logged_at FROM food_logs WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.FoodLog
	for rows.Next() {
		var f model.FoodLog
		var t time.Time
		if err := rows.Scan(&f.ID, &f.Name, &f.Calories, &f.Protein, &f.Fat, &f.Carbs, &f.Quantity, &t); err != nil {
			continue
		}
		formatted := t.Format("2006-01-02")
		f.LoggedAt = &formatted

		logs = append(logs, &f)
	}
	return logs, nil
}

func (r *Resolver) AddFoodLog(ctx context.Context, input model.FoodLogInput) (*model.FoodLog, error) {
	userID := ctx.Value("user_id").(int)

	// ✅ 安全處理指標型別（避免 nil pointer crash）
	if input.LoggedAt == nil {
		return nil, fmt.Errorf("請提供 logged_at 日期")
	}

	// ✅ 解參考 *string 並轉為 time.Time
	t, err := time.Parse("2006-01-02", *input.LoggedAt)
	if err != nil {
		return nil, fmt.Errorf("日期格式錯誤：%v", err)
	}

	// ✅ 寫入資料庫
	_, err = database.DB.Exec(`
		INSERT INTO food_logs (user_id, name, calories, protein, fat, carbs, quantity, logged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, input.Name, input.Calories, input.Protein, input.Fat, input.Carbs, input.Quantity, t,
	)
	if err != nil {
		return nil, err
	}

	// ✅ 回傳輸入內容（可用 *input.LoggedAt）
	return &model.FoodLog{
		Name:     input.Name,
		Calories: input.Calories,
		Protein:  input.Protein,
		Fat:      input.Fat,
		Carbs:    input.Carbs,
		Quantity: input.Quantity,
		LoggedAt: input.LoggedAt, // 這邊保留原本傳入的日期字串
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
