package sample

import (
	pcbook "pcbook/proto"
	proto "pcbook/proto"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *proto.Keyboard {
	keyboard := &proto.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

func NewCPU() *pcbook.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	number_of_cores := randomInt(2, 8)
	minGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minGhz, 5.0)

	cpu := &pcbook.CPU{
		Brand:       brand,
		Name:        name,
		NumberCores: cast.ToUint32(number_of_cores),
		MinGhz:      minGhz,
		MaxGhz:      maxGhz,
	}

	return cpu
}

func NewGPU() *pcbook.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)

	memory := &pcbook.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pcbook.Memory_GIGABYTE,
	}

	gpu := &pcbook.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}

	return gpu
}

func NewRAM() *pcbook.Memory {
	ram := &pcbook.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pcbook.Memory_GIGABYTE,
	}
	return ram
}

func NewSSD() *pcbook.Storage {
	ssd := &pcbook.Storage{
		Driver: pcbook.Storage_SSD,
		Memory: &pcbook.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pcbook.Memory_GIGABYTE,
		},
	}
	return ssd
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}

func NewHDD() *pcbook.Storage {
	ssd := &pcbook.Storage{
		Driver: pcbook.Storage_HDD,
		Memory: &pcbook.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pcbook.Memory_TERABYTE,
		},
	}
	return ssd
}

func NewScree() *pcbook.Screen {
	height := randomInt(1080, 4320)
	weight := height * 16 / 9
	screen := &pcbook.Screen{
		SizeInch: randomFloat32(13, 17),
		Resolution: &proto.Screen_Resolution{
			Height: uint32(height),
			Width:  uint32(weight),
		},
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
	return screen
}

func NewLaptop() *pcbook.Laptop {
	brand := randomLaptoBrand()
	name := randomLaptopName(brand)

	laptop := &pcbook.Laptop{
		Id:       randomID(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCPU(),
		Ram:      NewRAM(),
		Gpus:     []*pcbook.GPU{NewGPU()},
		Storages: []*pcbook.Storage{NewHDD(), NewSSD()},
		Screen:   NewScree(),
		Keyboard: NewKeyboard(),
		Weight: &pcbook.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   timestamppb.Now(),
	}
	return laptop
}
