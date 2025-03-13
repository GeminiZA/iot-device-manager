package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Device struct {
	ID        uint      `gorm:"primaryKey;not null;index"`
	Name      string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Telemetry struct {
	ID        uint           `gorm:"primaryKey;not null;index"`
	DeviceID  uint           `gorm:"not null;index"`
	Timestamp time.Time      `gorm:"not null"`
	Data      datatypes.JSON `gorm:"type:jsonb;notnull"`
}

// Device methods

func NewDevice(name string, id uint) *Device {
	device := Device{
		ID:        id,
		Name:      name,
		Status:    "offline",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &device
}

// Database Methods

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	db.AutoMigrate(Device{})
	db.AutoMigrate(Telemetry{})
	dr := DeviceRepository{
		db: db,
	}
	return &dr
}

func (r *DeviceRepository) CreateDevice(device *Device) error {
	res := r.db.Create(device)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *DeviceRepository) GetDevice(id uint) (*Device, error) {
	var device Device
	err := r.db.First(&device, id).Error
	return &device, err
}

// NewDeviceDetails struct for passing to UpdateDevice
type NewDeviceDetails struct {
	Status string
	Name   string
}

func (r *DeviceRepository) UpdateDevice(id uint, details *NewDeviceDetails) error {
	updateDevice := Device{}
	if details.Status != "" {
		updateDevice.Status = details.Status
	}
	if details.Name != "" {
		updateDevice.Name = details.Name
	}
	res := r.db.Model(&Device{}).Where("id = ?", id).Updates(updateDevice)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *DeviceRepository) AddTelemetry(id uint, data *datatypes.JSON) error {
	newTelemetryData := Telemetry{
		DeviceID:  id,
		Timestamp: time.Now(),
		Data:      *data,
	}
	res := r.db.Create(&newTelemetryData)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *DeviceRepository) GetAllTelemetryByDeviceID(deviceId uint) ([]Telemetry, error) {
	var telemetry []Telemetry
	res := r.db.Where("device_id = ?", deviceId).Find(&telemetry)
	if res.Error != nil {
		return nil, res.Error
	}
	return telemetry, nil
}

func (r *DeviceRepository) GetLatestTelemetryByDeviceID(deviceId uint) (*Telemetry, error) {
	var telemetry Telemetry
	res := r.db.Where("device_id = ?", deviceId).Order("timestamp DESC").First(&telemetry)
	if res.Error != nil {
		return nil, res.Error
	}
	return &telemetry, nil
}

func (r *DeviceRepository) DeleteDevice(id uint) error {
	res := r.db.Delete(&Device{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	res = r.db.Delete(&Telemetry{}, "device_id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
