package mqtt

import "testing"

func TestParseMessage(t *testing.T) {
	topic := "/fleet/vehicle/B1234ABC/location"

	tests := []struct {
		name    string
		topic   string
		payload string
		wantErr bool
	}{
		{"valid", topic, `{"vehicle_id":"B1234ABC","latitude":-6.2088,"longitude":106.8456,"timestamp":1715003456}`, false},
		{"tanpa timestamp", topic, `{"vehicle_id":"B1234ABC","latitude":-6.2088,"longitude":106.8456}`, true},
		{"vehicle beda dengan topic", "/fleet/vehicle/B9999ZZZ/location", `{"vehicle_id":"B1234ABC","latitude":-6.2088,"longitude":106.8456,"timestamp":1715003456}`, true},
		{"bukan json", topic, `hello`, true},
	}

	for _, tt := range tests {
		_, err := parseMessage(tt.topic, []byte(tt.payload))
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: wantErr=%v, got %v", tt.name, tt.wantErr, err)
		}
	}
}
