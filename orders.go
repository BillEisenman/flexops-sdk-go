package flexops

import (
	"context"
	"net/url"
)

// OrdersService reads orders through the legacy proxy. The proxy refuses API-key
// writes, so it has no create or cancel.
type OrdersService struct{ c *Client }

const ordersBase = "/api/ApiProxy/api/v1/Order"

func (s *OrdersService) GetNewOrders(ctx context.Context, params url.Values) (ApiResponse[[]Order], error) { return decode[ApiResponse[[]Order]](s.c.http.get(ctx, ordersBase+"/getNewOrderList", params)) }
func (s *OrdersService) GetByStatus(ctx context.Context, params url.Values) (ApiResponse[[]Order], error) { return decode[ApiResponse[[]Order]](s.c.http.get(ctx, ordersBase+"/getAllOrderListByStatus", params)) }
func (s *OrdersService) GetDetails(ctx context.Context, orderNumber string) (ApiResponse[Order], error) { return decode[ApiResponse[Order]](s.c.http.get(ctx, ordersBase+"/getCompleteOrderDetailsByOrderNumber", url.Values{"orderNumber": {orderNumber}})) }
func (s *OrdersService) GetExtendedDetails(ctx context.Context, orderNumber string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getExtendedOrderDetailsByOrderNumber", url.Values{"orderNumber": {orderNumber}})) }
func (s *OrdersService) GetStatus(ctx context.Context, orderNumber string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getIndividualOrderStatusByOrderNumber", url.Values{"orderNumber": {orderNumber}})) }
func (s *OrdersService) GetItems(ctx context.Context, orderNumber string) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getAllOrderItemsByOrderNumber", url.Values{"orderNumber": {orderNumber}})) }
func (s *OrdersService) GetShipMethods(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getAvailableShipMethodsList", nil)) }
func (s *OrdersService) GetWarehouses(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getActiveWarehouseList", nil)) }
func (s *OrdersService) GetCountryCodes(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getCountryNameCodeList", nil)) }
func (s *OrdersService) GetStatusTypes(ctx context.Context) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, ordersBase+"/getOrderStatusTypesList", nil)) }
