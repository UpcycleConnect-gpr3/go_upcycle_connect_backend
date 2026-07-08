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

// GetUserStats agrege l'activite de l'utilisateur : nombre d'annonces,
// de projets, CO2 total economise et repartition par categorie.
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

// FinanceStats : synthese financiere globale pour le back office admin.
type FinanceStats struct {
	RevenueCents        int `json:"revenue_cents"`         // total encaisse (achats payes)
	CommissionCents     int `json:"commission_cents"`      // commission UpcycleConnect
	PaidTransactions    int `json:"paid_transactions"`     // nombre d'achats payes
	ActiveSubscriptions int `json:"active_subscriptions"`  // abonnements actifs
	ObjectsCount        int `json:"objects_count"`
	ProjectsCount       int `json:"projects_count"`
}

// GetFinanceStats agrege les revenus (paiements payes) et abonnements.
func GetFinanceStats() FinanceStats {
	var s FinanceStats
	// Paiements payes : revenu total + commission + nombre.
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
