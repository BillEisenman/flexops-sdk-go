package flexops

import (
	"context"
	"net/url"
)

type AnalyticsService struct{ c *Client }

const analyticsBase = "/api/ApiProxy/api/v4/Analytics"

func (s *AnalyticsService) dateQ(startDate, endDate string) url.Values {
	v := url.Values{}
	if startDate != "" { v.Set("startDate", startDate) }
	if endDate != "" { v.Set("endDate", endDate) }
	return v
}

func (s *AnalyticsService) ShipmentsTrend(ctx context.Context, startDate, endDate string) (ApiResponse[[]ShipmentsTrend], error) { return decode[ApiResponse[[]ShipmentsTrend]](s.c.http.get(ctx, analyticsBase+"/ShipmentsTrend", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) CarrierSummary(ctx context.Context, startDate, endDate string) (ApiResponse[[]CarrierSummary], error) { return decode[ApiResponse[[]CarrierSummary]](s.c.http.get(ctx, analyticsBase+"/CarrierSummary", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) TopDestinations(ctx context.Context, startDate, endDate string, limit string) (ApiResponse[any], error) { q := s.dateQ(startDate, endDate); if limit != "" { q.Set("limit", limit) }; return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/TopDestinations", q)) }
func (s *AnalyticsService) InventoryMetrics(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/InventoryMetrics", nil)) }
func (s *AnalyticsService) StockByWarehouse(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/StockByWarehouse", nil)) }
func (s *AnalyticsService) OrderMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/OrderMetrics", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) OrderTrend(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/OrderTrend", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) TopSellingProducts(ctx context.Context, startDate, endDate string, limit string) (ApiResponse[any], error) { q := s.dateQ(startDate, endDate); if limit != "" { q.Set("limit", limit) }; return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/TopSellingProducts", q)) }
func (s *AnalyticsService) ReturnsMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnsMetrics", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) ReturnsTrend(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnsTrend", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) ReturnReasons(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnReasons", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) PerformanceMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/PerformanceMetrics", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) CarrierPerformance(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/CarrierPerformance", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) ShippingCostAnalytics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ShippingCostAnalytics", s.dateQ(startDate, endDate))) }
func (s *AnalyticsService) DeliveryPerformance(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/DeliveryPerformance", s.dateQ(startDate, endDate))) }
