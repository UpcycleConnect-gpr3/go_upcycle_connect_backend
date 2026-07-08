package stats_models

import (
	"go-upcycle_connect-backend/utils/log"
	"go-upcycle_connect-backend/var/database"
)

type CategoryCount struct {
	Category string `db:"category" json:"category"`
	Count    int    `db:"count" json:"count"`
}

type UserStats struct {
	ObjectsCount  int             `json:"objects_count"`
	ProjectsCount int             `json:"projects_count"`
	Co2Total      float64         `json:"co2_total"`
	ByCategory    []CategoryCount `json:"by_category"`
}

func GetUserStats(userId string) UserStats {
	stats := UserStats{ByCategory: []CategoryCount{}}

	if err := database.UpcycleConnect.Get(&stats.ObjectsCount,
		"SELECT COUNT(*) FROM OBJECTS WHERE user_id = ?", userId); err != nil {
		log.Database("STATS objects count", err)
	}

	if err := database.UpcycleConnect.Get(&stats.ProjectsCount,
		"SELECT COUNT(*) FROM PROJECTS WHERE user_id = ?", userId); err != nil {
		log.Database("STATS projects count", err)
	}

	if err := database.UpcycleConnect.Get(&stats.Co2Total,
		"SELECT COALESCE(SUM(score), 0) FROM OBJECTS WHERE user_id = ?", userId); err != nil {
		log.Database("STATS co2 total", err)
	}

	if err := database.UpcycleConnect.Select(&stats.ByCategory,
		"SELECT category, COUNT(*) AS count FROM OBJECTS WHERE user_id = ? GROUP BY category ORDER BY count DESC", userId); err != nil {
		log.Database("STATS by category", err)
	}

	return stats
}

type FinanceStats struct {
	RevenueCents        int `json:"revenue_cents"`
	CommissionCents     int `json:"commission_cents"`
	PaidTransactions    int `json:"paid_transactions"`
	ActiveSubscriptions int `json:"active_subscriptions"`
	ObjectsCount        int `json:"objects_count"`
	ProjectsCount       int `json:"projects_count"`
}

func GetFinanceStats() FinanceStats {
	var s FinanceStats

	_ = database.UpcycleConnect.QueryRow(
		"SELECT COALESCE(SUM(amount_cents),0), COALESCE(SUM(commission_cents),0), COUNT(*) FROM PAYMENTS WHERE status = 'paid'",
	).Scan(&s.RevenueCents, &s.CommissionCents, &s.PaidTransactions)

	_ = database.UpcycleConnect.QueryRow(
		"SELECT COUNT(*) FROM SUBSCRIPTIONS WHERE status = 'active'",
	).Scan(&s.ActiveSubscriptions)

	_ = database.UpcycleConnect.QueryRow("SELECT COUNT(*) FROM OBJECTS").Scan(&s.ObjectsCount)
	_ = database.UpcycleConnect.QueryRow("SELECT COUNT(*) FROM PROJECTS").Scan(&s.ProjectsCount)

	return s
}
