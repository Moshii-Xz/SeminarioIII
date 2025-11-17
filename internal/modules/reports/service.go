package reports

import (
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) parseDates(start, end time.Time) (time.Time, time.Time, error) {
	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("fecha_fin debe ser posterior a fecha_inicio")
	}
	return start, end, nil
}

func (s *Service) GetSalesSummary(req SalesSummaryRequest) (*SalesSummaryResponse, error) {
	start, end, err := s.parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	summary, err := s.repo.GetSalesSummary(start, end)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de ventas: %w", err)
	}

	return &SalesSummaryResponse{
		TotalOrders:       summary.TotalOrders,
		TotalRevenue:      summary.TotalRevenue,
		AverageOrderValue: summary.AverageOrderValue,
		TotalItemsSold:    summary.TotalItemsSold,
	}, nil
}

func (s *Service) GetTopProducts(req ReportFilterRequest) ([]TopProductResponse, error) {
	start, end, err := s.parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	products, err := s.repo.GetTopSellingProducts(start, end, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo productos top: %w", err)
	}

	resp := make([]TopProductResponse, len(products))
	for i, p := range products {
		resp[i] = TopProductResponse(p)
	}
	return resp, nil
}

func (s *Service) GetLowStock(threshold int) ([]InventoryStatusResponse, error) {
	products, err := s.repo.GetLowStockProducts(threshold)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo productos con poco stock: %w", err)
	}

	resp := make([]InventoryStatusResponse, len(products))
	for i, p := range products {
		resp[i] = InventoryStatusResponse(p)
	}
	return resp, nil
}

func (s *Service) GetDailySales(req SalesSummaryRequest) ([]DailySalesResponse, error) {
	start, end, err := s.parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	entries, err := s.repo.GetDailySalesTrend(start, end)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ventas diarias: %w", err)
	}

	resp := make([]DailySalesResponse, len(entries))
	for i, entry := range entries {
		resp[i] = DailySalesResponse(entry)
	}
	return resp, nil
}

func (s *Service) GetTopCustomers(req ReportFilterRequest) ([]CustomerRankingResponse, error) {
	start, end, err := s.parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	customers, err := s.repo.GetTopCustomers(start, end, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo clientes top: %w", err)
	}

	resp := make([]CustomerRankingResponse, len(customers))
	for i, c := range customers {
		resp[i] = CustomerRankingResponse(c)
	}
	return resp, nil
}

func (s *Service) GetPaymentMethodSummary(req SalesSummaryRequest) ([]PaymentMethodSummaryResponse, error) {
	start, end, err := s.parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	summaries, err := s.repo.GetPaymentMethodSummary(start, end)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo resumen de métodos de pago: %w", err)
	}

	resp := make([]PaymentMethodSummaryResponse, len(summaries))
	for i, summary := range summaries {
		resp[i] = PaymentMethodSummaryResponse(summary)
	}
	return resp, nil
}

func (s *Service) GetPendingPayments() ([]PendingPaymentResponse, error) {
	pendings, err := s.repo.GetPendingPaymentsReport()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo pagos pendientes: %w", err)
	}

	resp := make([]PendingPaymentResponse, len(pendings))
	for i, p := range pendings {
		resp[i] = PendingPaymentResponse(p)
	}
	return resp, nil
}

