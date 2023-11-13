package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ml4wireless/specpipe/common"
	"github.com/nats-io/nats.go/jetstream"
)

type Store struct {
	kv jetstream.KeyValue
}

func NewStore(kv jetstream.KeyValue) *Store {
	return &Store{kv}
}

func (s *Store) GetFmDevice(ctx context.Context, deviceName string) (bool, *common.FMDevice, error) {
	entry, err := s.kv.Get(ctx, common.KVStoreKey(common.FM, deviceName))
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	var device common.FMDevice
	if err = json.Unmarshal(entry.Value(), &device); err != nil {
		return false, nil, err
	}
	return true, &device, nil
}

func (s *Store) GetFmDevices(ctx context.Context) ([]*common.FMDevice, error) {
	keys, err := s.kv.Keys(ctx)
	if err != nil {
		return nil, err
	}
	devices := []*common.FMDevice{}
	for _, key := range keys {
		exist, device, err := s.GetFmDevice(ctx, common.DeviceNameFromKey(key))
		if err != nil {
			return nil, err
		}
		if exist {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (s *Store) UpdateFmDevice(ctx context.Context, fmDevice *common.FMDevice) error {
	deviceInfoBytes, err := json.Marshal(fmDevice)
	if err != nil {
		return err
	}
	_, err = s.kv.Put(ctx, common.KVStoreKey(common.FM, fmDevice.Name), deviceInfoBytes)
	return err
}

func (s *Store) DeleteFmDevice(ctx context.Context, deviceName string) error {
	if err := s.kv.Delete(ctx, common.KVStoreKey(common.FM, deviceName)); err != nil {
		return err
	}
	return nil
}
