package emulator

import "encoding/json"

type EmulatorRegistration struct {
	Size int
}

type registration struct {
	Size int    `json:size`
	Data string `json:data`
}

func (self *EmulatorRegistration) RegistrationData() (*[]byte, error) {
	r := &registration{
		Size: self.Size,
		Data: self.randomString(self.Size),
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (self *EmulatorRegistration) randomString(strlen int) string {
	chars := "0123456789"
	result := ""

	for i := 0; i < strlen; i++ {
		start := i % 10
		result += chars[start : start+1]
	}

	return result
}
