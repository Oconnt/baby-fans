package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleParent Role = "parent"
	RoleChild  Role = "child"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Role      Role           `gorm:"not null" json:"role"`
	Password  string         `json:"password"`
	LoginCode string         `json:"login_code"`
	Points    int            `gorm:"default:0" json:"points"`
	OpenID    string         `gorm:"index" json:"openid"`
	UnionID   string         `gorm:"index" json:"unionid"`
	Nickname  string         `json:"nickname"`
	AvatarURL string         `json:"avatar_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserBinding struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentID  uint      `gorm:"index" json:"parent_id"`
	ChildID   uint      `gorm:"index" json:"child_id"`
	BindCode  string    `gorm:"uniqueIndex" json:"bind_code"`
	Status    string    `gorm:"default:'pending'" json:"status"` // pending, accepted
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Parent User `gorm:"foreignKey:ParentID" json:"parent"`
	Child  User `gorm:"foreignKey:ChildID" json:"child"`
}

type ParentChild struct {
	ParentID uint `gorm:"primaryKey" json:"parent_id"`
	ChildID  uint `gorm:"primaryKey" json:"child_id"`
}

type PointsRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Amount     int       `json:"amount"`
	Reason     string    `json:"reason"`
	OperatorID uint      `json:"operator_id"` // The parent who added/subtracted points
	Operator   User      `gorm:"foreignKey:OperatorID" json:"operator"`
	CreatedAt  time.Time `json:"created_at"`
}

type ShopItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       int            `json:"price"`
	Stock       int            `json:"stock"`
	ImagePath   string         `json:"image_path"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RedemptionStatus string

const (
	RedemptionPending   RedemptionStatus = "pending"
	RedemptionCompleted RedemptionStatus = "completed"
	RedemptionCancelled RedemptionStatus = "cancelled"
)

type Redemption struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"index" json:"user_id"`
	ItemID    uint             `json:"item_id"`
	Status    RedemptionStatus `gorm:"default:'pending'" json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	User User     `gorm:"foreignKey:UserID" json:"user"`
	Item ShopItem `gorm:"foreignKey:ItemID" json:"item"`
}

type FaceLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	PhotoPath string    `json:"photo_path"`
	CreatedAt time.Time `json:"created_at"`
}

type PointsTemplate struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Amount  int    `json:"amount"`
}

// Task status constants
const (
	TaskPending   = 1 // 待办
	TaskCompleted = 2 // 已完成
	TaskExpired   = 3 // 已过期
)

type TaskTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	CreatedAt   time.Time `json:"created_at"`
}

type Task struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Points       int        `json:"points"`
	Status       int        `gorm:"default:1" json:"status"` // 1=待办, 2=已完成, 3=已过期
	PublisherID  uint       `gorm:"index" json:"publisher_id"`
	HandlerID    uint       `gorm:"index" json:"handler_id"`
	PublishTime  time.Time  `json:"publish_time"`
	ExpireTime   time.Time  `json:"expire_time"`
	CompleteTime *time.Time `json:"complete_time,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`

	Publisher User `gorm:"foreignKey:PublisherID" json:"publisher"`
	Handler   User `gorm:"foreignKey:HandlerID" json:"handler"`
}
