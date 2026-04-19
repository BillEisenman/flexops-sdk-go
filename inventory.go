package flexops

import (
	"context"
	"net/url"
)

type InventoryService struct{ c *Client }

const inventoryProxy = "/api/ApiProxy"

func (s *InventoryService) PostAsnReceipt(ctx context.Context, receipt any) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, inventoryProxy+"/api/v1/Inventory/postNewAsnReceipt", receipt)) }
func (s *InventoryService) GetWarehouseSnapshot(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getWarehouseInventorySnapshot", params)) }
func (s *InventoryService) GetCompleteSnapshot(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getCompleteInventorySnapshot", params)) }
func (s *InventoryService) GetPartNumbers(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getPartNumberList", params)) }
func (s *InventoryService) UpdateInventory(ctx context.Context, data any) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.post(ctx, inventoryProxy+"/api/v2/Inventory/postCustomerInventoryUpdate", data)) }
