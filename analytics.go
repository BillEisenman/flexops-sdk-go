package flexops

import (
	"context"
	"net/url"
)

// AnalyticsService reads analytics through the legacy proxy (/api/ApiProxy/api/v4/Analytics).
//
// With an API key, each method requires the permission for the data it reads, named on
// the method: shipping:read for shipment, carrier and delivery metrics, inventory:read
// for inventory metrics, orders:read for order metrics and returns:read for returns
// metrics. A key without it gets 403 API_KEY_SCOPE_REQUIRED. Keys created with the
// default permissions have all four.
type AnalyticsService struct{ c *Client }

const analyticsBase = "/api/ApiProxy/api/v4/Analytics"

func (s *AnalyticsService) dateQ(startDate, endDate string) url.Values {
	v := url.Values{}
	if startDate != "" { v.Set("startDate", startDate) }
	if endDate != "" { v.Set("endDate", endDate) }
	return v
}

// ShipmentsTrend returns the shipments trend over time. With an API key it requires shipping:read.
func (s *AnalyticsService) ShipmentsTrend(ctx context.Context, startDate, endDate string) (ApiResponse[[]ShipmentsTrend], error) { return decode[ApiResponse[[]ShipmentsTrend]](s.c.http.get(ctx, analyticsBase+"/ShipmentsTrend", s.dateQ(startDate, endDate))) }
// CarrierSummary returns carrier usage. With an API key it requires shipping:read.
func (s *AnalyticsService) CarrierSummary(ctx context.Context, startDate, endDate string) (ApiResponse[[]CarrierSummary], error) { return decode[ApiResponse[[]CarrierSummary]](s.c.http.get(ctx, analyticsBase+"/CarrierSummary", s.dateQ(startDate, endDate))) }
// TopDestinations returns the top shipping destinations. With an API key it requires shipping:read.
func (s *AnalyticsService) TopDestinations(ctx context.Context, startDate, endDate string, limit string) (ApiResponse[any], error) { q := s.dateQ(startDate, endDate); if limit != "" { q.Set("limit", limit) }; return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/TopDestinations", q)) }
// InventoryMetrics returns inventory metrics (stock levels, low-stock alerts). With an API key it requires inventory:read.
func (s *AnalyticsService) InventoryMetrics(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/InventoryMetrics", nil)) }
// StockByWarehouse returns stock levels by warehouse. With an API key it requires inventory:read.
func (s *AnalyticsService) StockByWarehouse(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/StockByWarehouse", nil)) }
// OrderMetrics returns order metrics (volume, revenue). With an API key it requires orders:read.
func (s *AnalyticsService) OrderMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/OrderMetrics", s.dateQ(startDate, endDate))) }
// OrderTrend returns the order trend over time. With an API key it requires orders:read.
func (s *AnalyticsService) OrderTrend(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/OrderTrend", s.dateQ(startDate, endDate))) }
// TopSellingProducts returns the top selling products. With an API key it requires orders:read.
func (s *AnalyticsService) TopSellingProducts(ctx context.Context, startDate, endDate string, limit string) (ApiResponse[any], error) { q := s.dateQ(startDate, endDate); if limit != "" { q.Set("limit", limit) }; return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/TopSellingProducts", q)) }
// ReturnsMetrics returns returns metrics. With an API key it requires returns:read.
func (s *AnalyticsService) ReturnsMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnsMetrics", s.dateQ(startDate, endDate))) }
// ReturnsTrend returns the returns trend over time. With an API key it requires returns:read.
func (s *AnalyticsService) ReturnsTrend(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnsTrend", s.dateQ(startDate, endDate))) }
// ReturnReasons returns the return reasons breakdown. With an API key it requires returns:read.
func (s *AnalyticsService) ReturnReasons(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ReturnReasons", s.dateQ(startDate, endDate))) }
// PerformanceMetrics returns fulfillment performance metrics. With an API key it requires shipping:read.
func (s *AnalyticsService) PerformanceMetrics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/PerformanceMetrics", s.dateQ(startDate, endDate))) }
// CarrierPerformance returns carrier delivery performance. With an API key it requires shipping:read.
func (s *AnalyticsService) CarrierPerformance(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/CarrierPerformance", s.dateQ(startDate, endDate))) }
// ShippingCostAnalytics returns shipping cost analytics. With an API key it requires shipping:read.
func (s *AnalyticsService) ShippingCostAnalytics(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/ShippingCostAnalytics", s.dateQ(startDate, endDate))) }
// DeliveryPerformance returns delivery performance (on-time percentage). With an API key it requires shipping:read.
func (s *AnalyticsService) DeliveryPerformance(ctx context.Context, startDate, endDate string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, analyticsBase+"/DeliveryPerformance", s.dateQ(startDate, endDate))) }
