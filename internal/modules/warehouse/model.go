package warehouse

type Warehouse struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Name         string    `json:"name"`
	UpdatedAt    *string   `json:"updated_at,omitempty"`
	CreatedAt    *string   `json:"created_at,omitempty"`
	Longitude    float64   `json:"longitude"`
	Latitude     float64   `json:"latitude"`
	WBID         *string   `json:"wb_id,omitempty"`
	WorkingHours *string   `json:"working_hours,omitempty"`
	Address      string    `json:"address"`
}

type CreateWarehouseRequest struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	WorkingHours string   `json:"working_hours"`
}

type UpdateWarehouseRequest struct {
	Name         *string  `json:"name,omitempty"`
	Address      *string  `json:"address,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	WorkingHours *string  `json:"working_hours,omitempty"`
} 