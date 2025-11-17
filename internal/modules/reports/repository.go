package reports

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type SalesSummary struct {
	TotalOrders       int64   `json:"total_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
	TotalItemsSold    int64   `json:"total_items_sold"`
}

func (r *Repository) GetSalesSummary(startDate, endDate time.Time) (*SalesSummary, error) {
	summary := SalesSummary{}

	query := `
		SELECT
			COALESCE(COUNT(DISTINCT c.id_compra), 0) AS total_orders,
			COALESCE(SUM(d.cantidad * d.precio_unitario), 0) AS total_revenue,
			COALESCE(SUM(d.cantidad), 0) AS total_items_sold
		FROM compra c
		LEFT JOIN detalle_compra d ON d.id_compra = c.id_compra
		WHERE c.fecha_compra BETWEEN ? AND ?`

	if err := r.db.Raw(query, startDate, endDate).Scan(&summary).Error; err != nil {
		return nil, err
	}

	if summary.TotalOrders > 0 {
		summary.AverageOrderValue = summary.TotalRevenue / float64(summary.TotalOrders)
	}

	return &summary, nil
}

type TopProduct struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitsSold   int64   `json:"units_sold"`
	Revenue     float64 `json:"revenue"`
}

func (r *Repository) GetTopSellingProducts(startDate, endDate time.Time, limit int) ([]TopProduct, error) {
	if limit <= 0 {
		limit = 5
	}

	var products []TopProduct
	query := `
		SELECT
			p.id_producto AS product_id,
			p.nombre AS product_name,
			COALESCE(SUM(d.cantidad), 0) AS units_sold,
			COALESCE(SUM(d.cantidad * d.precio_unitario), 0) AS revenue
		FROM detalle_compra d
		JOIN producto p ON p.id_producto = d.id_producto
		JOIN compra c ON c.id_compra = d.id_compra
		WHERE c.fecha_compra BETWEEN ? AND ?
		GROUP BY p.id_producto, p.nombre
		ORDER BY revenue DESC
		LIMIT ?`

	if err := r.db.Raw(query, startDate, endDate, limit).Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

type InventoryStatus struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

func (r *Repository) GetLowStockProducts(threshold int) ([]InventoryStatus, error) {
	if threshold <= 0 {
		threshold = 10
	}

	var products []InventoryStatus
	query := `
		SELECT
			id_producto AS product_id,
			nombre AS product_name,
			stock
		FROM producto
		WHERE stock <= ?
		ORDER BY stock ASC`

	if err := r.db.Raw(query, threshold).Scan(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

type DailySalesEntry struct {
	Date        time.Time `json:"date"`
	TotalOrders int64     `json:"total_orders"`
	Revenue     float64   `json:"revenue"`
}

func (r *Repository) GetDailySalesTrend(startDate, endDate time.Time) ([]DailySalesEntry, error) {
	var entries []DailySalesEntry
	query := `
		SELECT
			c.fecha_compra AS date,
			COALESCE(COUNT(DISTINCT c.id_compra), 0) AS total_orders,
			COALESCE(SUM(d.cantidad * d.precio_unitario), 0) AS revenue
		FROM compra c
		LEFT JOIN detalle_compra d ON d.id_compra = c.id_compra
		WHERE c.fecha_compra BETWEEN ? AND ?
		GROUP BY c.fecha_compra
		ORDER BY c.fecha_compra ASC`

	if err := r.db.Raw(query, startDate, endDate).Scan(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

type CustomerRanking struct {
	CustomerID   uint    `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	OrdersCount  int64   `json:"orders_count"`
	TotalSpent   float64 `json:"total_spent"`
}

func (r *Repository) GetTopCustomers(startDate, endDate time.Time, limit int) ([]CustomerRanking, error) {
	if limit <= 0 {
		limit = 5
	}

	var rankings []CustomerRanking
	query := `
		SELECT
			u.id_usuario AS customer_id,
			u.nombre AS customer_name,
			COALESCE(COUNT(DISTINCT c.id_compra), 0) AS orders_count,
			COALESCE(SUM(d.cantidad * d.precio_unitario), 0) AS total_spent
		FROM compra c
		JOIN cliente cl ON cl.id_cliente = c.id_cliente
		JOIN usuario u ON u.id_usuario = cl.id_cliente
		LEFT JOIN detalle_compra d ON d.id_compra = c.id_compra
		WHERE c.fecha_compra BETWEEN ? AND ?
		GROUP BY u.id_usuario, u.nombre
		ORDER BY total_spent DESC
		LIMIT ?`

	if err := r.db.Raw(query, startDate, endDate, limit).Scan(&rankings).Error; err != nil {
		return nil, err
	}

	return rankings, nil
}

type PaymentMethodSummary struct {
	MethodID    *uint   `json:"method_id"`
	MethodName  *string `json:"method_name"`
	TotalAmount float64 `json:"total_amount"`
	Payments    int64   `json:"payments"`
}

func (r *Repository) GetPaymentMethodSummary(startDate, endDate time.Time) ([]PaymentMethodSummary, error) {
	var summaries []PaymentMethodSummary
	query := `
		SELECT
			mp.id_metodo_pago AS method_id,
			mp.nombre AS method_name,
			COALESCE(SUM(p.monto), 0) AS total_amount,
			COALESCE(COUNT(p.id_pago), 0) AS payments
		FROM pago p
		LEFT JOIN metodo_pago mp ON mp.id_metodo_pago = p.id_metodo_pago
		WHERE p.fecha_pago BETWEEN ? AND ?
		GROUP BY mp.id_metodo_pago, mp.nombre
		ORDER BY total_amount DESC`

	if err := r.db.Raw(query, startDate, endDate).Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

type PendingPayment struct {
	OrderID       uint      `json:"order_id"`
	CustomerName  string    `json:"customer_name"`
	OrderTotal    float64   `json:"order_total"`
	TotalPaid     float64   `json:"total_paid"`
	PendingAmount float64   `json:"pending_amount"`
	OrderDate     time.Time `json:"order_date"`
}

func (r *Repository) GetPendingPaymentsReport() ([]PendingPayment, error) {
	var report []PendingPayment
	query := `
		WITH order_totals AS (
			SELECT
				id_compra,
				SUM(cantidad * precio_unitario) AS total
			FROM detalle_compra
			GROUP BY id_compra
		),
		payment_totals AS (
			SELECT
				id_compra,
				SUM(monto) AS total
			FROM pago
			GROUP BY id_compra
		)
		SELECT
			c.id_compra AS order_id,
			u.nombre AS customer_name,
			COALESCE(ot.total, 0) AS order_total,
			COALESCE(pt.total, 0) AS total_paid,
			COALESCE(ot.total, 0) - COALESCE(pt.total, 0) AS pending_amount,
			c.fecha_compra AS order_date
		FROM compra c
		JOIN cliente cl ON cl.id_cliente = c.id_cliente
		JOIN usuario u ON u.id_usuario = cl.id_cliente
		LEFT JOIN order_totals ot ON ot.id_compra = c.id_compra
		LEFT JOIN payment_totals pt ON pt.id_compra = c.id_compra
		WHERE COALESCE(ot.total, 0) > COALESCE(pt.total, 0)
		ORDER BY c.fecha_compra DESC`

	if err := r.db.Raw(query).Scan(&report).Error; err != nil {
		return nil, err
	}

	return report, nil
}
