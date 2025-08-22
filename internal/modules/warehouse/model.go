package warehouse

type Warehouse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	Name         string  `json:"name"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	WBID         *string `json:"wb_id,omitempty"`
	WorkingHours *string `json:"working_hours,omitempty"`
	Address      string  `json:"address"`
}

type CreateWarehouseRequest struct {
	Name         string  `json:"name" binding:"required,min=1,max=255"`
	Address      string  `json:"address" binding:"required,min=1,max=500"`
	Latitude     float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude    float64 `json:"longitude" binding:"required,min=-180,max=180"`
	WorkingHours string  `json:"working_hours" binding:"required,min=1,max=100"`
}

type UpdateWarehouseRequest struct {
	Name         *string  `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Address      *string  `json:"address,omitempty" binding:"omitempty,min=1,max=500"`
	Latitude     *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`
	WorkingHours *string  `json:"working_hours,omitempty" binding:"omitempty,min=1,max=100"`
}

type CreateBatchWarehouseRequest struct {
	Warehouses []CreateWarehouseRequest `json:"warehouses" binding:"required,min=1,max=100"`
}
