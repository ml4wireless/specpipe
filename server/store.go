package server

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Store struct {
	conn *nats.Conn
	kv   jetstream.KeyValue
}

func NewStore(conn *nats.Conn, kv jetstream.KeyValue) *Store {
	return &Store{conn, kv}
}

func (s *Store) GetDevice(ctx context.Context, sdrType common.SDRType, deviceName string) (bool, common.Device, error) {
	entry, err := s.kv.Get(ctx, common.KVStoreKey(sdrType, deviceName))
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	switch sdrType {
	case common.FM:
		var device common.FMDevice
		if err = json.Unmarshal(entry.Value(), &device); err != nil {
			return false, nil, err
		}
		return true, &device, nil
	}
	return false, nil, nil
}

func (s *Store) GetDevices(ctx context.Context, sdrType common.SDRType) ([]common.Device, error) {
	keys, err := s.kv.Keys(ctx)
	if err != nil {
		return nil, err
	}
	devices := []common.Device{}
	for _, key := range keys {
		exist, device, err := s.GetDevice(ctx, sdrType, common.DeviceNameFromKey(key))
		if err != nil {
			return nil, err
		}
		if exist {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (s *Store) UpdateDevice(ctx context.Context, sdrType common.SDRType, deviceName string, deviceInfo common.Device) error {
	deviceInfoBytes, err := json.Marshal(deviceInfo)
	if err != nil {
		return err
	}
	if _, err = s.kv.Put(ctx, common.KVStoreKey(sdrType, deviceName), deviceInfoBytes); err != nil {
		return err
	}
	_, err = s.conn.Request(common.ClusterSubject(sdrType, deviceName, common.WatchConfigCmd), nil, 3*time.Second)
	return err
}

func (s *Store) DeleteDevice(ctx context.Context, sdrType common.SDRType, deviceName string) error {
	if err := s.kv.Delete(ctx, common.KVStoreKey(sdrType, deviceName)); err != nil {
		return err
	}
	return nil
}
