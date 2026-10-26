package flexops

import (
	"context"
	"net/url"
)

// InventoryService reads inventory through the legacy proxy. The proxy refuses API-key
// writes, so it has no ASN receipt or inventory update.
//
// With an API key, every method requires the inventory:read permission; a key without it
// gets 403 API_KEY_SCOPE_REQUIRED. Keys created with the default permissions include
// inventory:read.
type InventoryService struct{ c *Client }

const inventoryProxy = "/api/ApiProxy"

func (s *InventoryService) GetWarehouseSnapshot(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getWarehouseInventorySnapshot", params)) }
func (s *InventoryService) GetCompleteSnapshot(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getCompleteInventorySnapshot", params)) }
func (s *InventoryService) GetPartNumbers(ctx context.Context, params url.Values) (ApiResponse[any], error) { return decode[ApiResponse[any]](s.c.http.get(ctx, inventoryProxy+"/api/v1/Inventory/getPartNumberList", params)) }
