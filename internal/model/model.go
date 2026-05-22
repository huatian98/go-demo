package model

import "time"

// User 用户
type User struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	Openid         string    `gorm:"size:64;uniqueIndex;not null" json:"openid"`
	Unionid        string    `gorm:"size:64" json:"unionid,omitempty"`
	Nickname       string    `gorm:"size:64" json:"nickname"`
	Avatar         string    `gorm:"size:255" json:"avatar"`
	Phone          string    `gorm:"size:20" json:"phone,omitempty"`
	DefaultClaimID uint64    `gorm:"column:default_claim_id" json:"default_claim_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// WineSeries 酒品系列
type WineSeries struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CoverURL    string    `gorm:"size:255" json:"cover_url"`
	BasePrice   float64   `gorm:"type:decimal(10,2);not null" json:"base_price"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (WineSeries) TableName() string { return "wine_series" }

// Cellar 酒窖
type Cellar struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Address   string    `gorm:"size:255" json:"address"`
	Province  string    `gorm:"size:32" json:"province"`
	City      string    `gorm:"size:32" json:"city"`
	Capacity  int       `gorm:"default:0" json:"capacity"`
	Available int       `gorm:"default:0" json:"available"`
	CreatedAt time.Time `json:"created_at"`
}

func (Cellar) TableName() string { return "cellars" }

// WineJar 单坛酒
type WineJar struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	Code            string     `gorm:"size:32;uniqueIndex;not null" json:"code"`
	SeriesID        uint64     `gorm:"index;not null" json:"series_id"`
	CellarID        uint64     `gorm:"index;not null" json:"cellar_id"`
	Year            int        `json:"year"`
	CoverURL        string     `gorm:"size:255" json:"cover_url"`
	CurrentOwnerID  uint64     `gorm:"index" json:"current_owner_id"`
	Status          string     `gorm:"size:16;index;default:'idle'" json:"status"`
	ClaimedAt       *time.Time `json:"claimed_at"`
	ExpectedReadyAt *time.Time `json:"expected_ready_at"`
	Version         int        `gorm:"default:0" json:"-"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (WineJar) TableName() string { return "wine_jars" }

// Claim 认领记录
type Claim struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	ClaimNo        string     `gorm:"size:32;uniqueIndex;not null" json:"claim_no"`
	UserID         uint64     `gorm:"index;not null" json:"user_id"`
	JarID          uint64     `gorm:"index;not null" json:"jar_id"`
	CellarID       uint64     `gorm:"index;not null" json:"cellar_id"`
	ApplicantName  string     `gorm:"size:32;not null" json:"applicant_name"`
	ContactPhone   string     `gorm:"size:32;not null" json:"contact_phone"`
	Price          float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status         string     `gorm:"size:16;index;default:'pending'" json:"status"`
	PaidAt         *time.Time `json:"paid_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (Claim) TableName() string { return "claims" }

// Payment 支付订单
type Payment struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	ClaimID       uint64     `gorm:"index;not null" json:"claim_id"`
	OutTradeNo    string     `gorm:"size:64;uniqueIndex;not null" json:"out_trade_no"`
	Channel       string     `gorm:"size:16;not null" json:"channel"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string     `gorm:"size:16;default:'pending'" json:"status"`
	TransactionID string     `gorm:"size:64" json:"transaction_id"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (Payment) TableName() string { return "payments" }

// JarMetrics 酒坛实时指标(时序数据)
type JarMetrics struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	JarID              uint64    `gorm:"index:idx_jar_time,priority:1;not null" json:"jar_id"`
	PhLevel            float64   `gorm:"type:decimal(4,2);not null" json:"ph_level"`
	PhStatus           string    `gorm:"size:16" json:"ph_status"`
	CellarTemperature  float64   `gorm:"type:decimal(4,1);not null" json:"cellar_temperature"`
	CellarHumidity     float64   `gorm:"type:decimal(4,1);not null" json:"cellar_humidity"`
	OutdoorTemperature float64   `gorm:"type:decimal(4,1)" json:"outdoor_temperature"`
	OutdoorLux         int       `json:"outdoor_lux"`
	BreathingState     string    `gorm:"size:32" json:"breathing_state"`
	AINarrative        string    `gorm:"type:text" json:"ai_narrative"`
	RecordedAt         time.Time `gorm:"index:idx_jar_time,priority:2;not null" json:"recorded_at"`
	CreatedAt          time.Time `json:"created_at"`
}

func (JarMetrics) TableName() string { return "jar_metrics" }

// JarTimeline 成长故事时间线
type JarTimeline struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	JarID       uint64    `gorm:"index;not null" json:"jar_id"`
	EventType   string    `gorm:"size:32" json:"event_type"`
	Title       string    `gorm:"size:64;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	HappenedAt  time.Time `gorm:"not null" json:"happened_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (JarTimeline) TableName() string { return "jar_timeline" }

// WineComponent 黄酒成分科普
type WineComponent struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:32;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	IconURL     string `gorm:"size:255" json:"icon_url"`
	Sort        int    `gorm:"default:0" json:"sort"`
}

func (WineComponent) TableName() string { return "wine_components" }

// CraftStep 古法工艺
type CraftStep struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	StepNo      int    `gorm:"not null" json:"step_no"`
	Name        string `gorm:"size:32;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"size:255" json:"image_url"`
}

func (CraftStep) TableName() string { return "craft_steps" }
