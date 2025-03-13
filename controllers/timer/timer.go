package timer

import (
	"fmt"
	"sync"
	"time"

	"github.com/GeminiZA/iot-device-manager/models"
	"github.com/gofiber/fiber/v2/log"
)

type Timer struct {
	dr        *models.DeviceRepository
	running   bool
	mux       sync.Mutex
	wg        sync.WaitGroup
	deviceIds []uint
}

func NewTimer(dr *models.DeviceRepository) *Timer {
	return &Timer{
		dr:      dr,
		running: false,
	}
}

func (timer *Timer) Run() error {
	if timer.running {
		return fmt.Errorf("timer already running")
	}
	timer.mux.Lock()
	devices, err := timer.dr.GetAllDevices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		timer.deviceIds = append(timer.deviceIds, device.ID)
	}
	timer.running = true
	timer.mux.Unlock()
	go timer.run()
	timer.wg.Add(1)
	return nil
}

func (timer *Timer) Stop() {
	timer.mux.Lock()
	timer.running = false
	timer.mux.Unlock()
	timer.wg.Wait()
}

func (timer *Timer) run() {
	for timer.running {
		for _, deviceId := range timer.deviceIds {
			tel, err := timer.dr.GetLatestTelemetryByDeviceID(deviceId)
			if err != nil {
				log.Error(err)
				continue
			}
			if tel == nil {
				continue
			}
			if tel.Timestamp.Add(30 * time.Second).Before(time.Now()) {
				err := timer.dr.UpdateDevice(deviceId, &models.NewDeviceDetails{
					Status: "offline",
				})
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}
		time.Sleep(time.Second)
	}
	timer.wg.Done()
}

func (timer *Timer) RemoveDevice(id uint) {
	timer.mux.Lock()
	defer timer.mux.Unlock()
	for i := range timer.deviceIds {
		if timer.deviceIds[i] == id {
			timer.deviceIds[len(timer.deviceIds)-1] = timer.deviceIds[i]
			timer.deviceIds = timer.deviceIds[:len(timer.deviceIds)-2]
			break
		}
	}
}

func (timer *Timer) AddDevice(id uint) {
	timer.mux.Lock()
	defer timer.mux.Unlock()
	timer.deviceIds = append(timer.deviceIds, id)
}
