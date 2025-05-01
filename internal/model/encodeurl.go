package model

type EncodeURL struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Original  string `json:"original_url" gorm:"not null"`
	Strategy  string `json:"strategy" gorm:"not null"`
	TenantID  string `json:"tenant_id" gorm:"not null;index:idx_tenant_url"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
}

type CreateEncodeURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	Strategy    string `json:"strategy" binding:"required,oneof=random sequential"`
	TenantID    string `json:"tenant_id"`
}

type EncodeURLResponse struct {
	EncodeURL   string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	TenantID    string `json:"tenant_id" binding:"required"`
}
