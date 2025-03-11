package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Device struct {
	ID        uint           `gorm:"primaryKey;not null;index"`
	Name      string         `gorm:"unique;not null"`
	Status    string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	Telemetry datatypes.JSON `gorm:"type:jsonb"`
}

type TelemetryData map[string]interface{}

// Device methods

func NewDevice(name string, id uint) *Device {
	device := Device{
		ID:        id,
		Name:      name,
		Status:    "offline",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Telemetry: datatypes.JSON{},
	}
	return &device
}

// Database Methods

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	db.AutoMigrate(Device{})
	dr := DeviceRepository{
		db: db,
	}
	return &dr
}

func (r *DeviceRepository) CreateDevice(asset *Device) error {
	res := r.db.Create(asset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *DeviceRepository) GetDevice(id uint) (*Device, error) {
	var asset Device
	err := r.db.First(&asset, id).Error
	return &asset, err
}

// NewDeviceDetails struct for passing to UpdateDevice
type NewDeviceDetails struct {
	Telemetry datatypes.JSON
	Status    string
}

func (r *DeviceRepository) UpdateDevice(id uint, details *NewDeviceDetails) error {
	res := r.db.Model(&Device{}).Where("id = ?", id).Updates(Device{
		Status:    details.Status,
		Telemetry: details.Telemetry,
	})
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *DeviceRepository) DeleteDevice(id uint) error {
	res := r.db.Delete(&Device{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
