package message_converters

import (
	"testing"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/gen/go/flightpath"
)

func TestGpsRawIntToProtobuf(t *testing.T) {
	t.Run("full GPS message", func(t *testing.T) {
		msg := &common.MessageGpsRawInt{
			TimeUsec:          1234567890,
			FixType:           common.GPS_FIX_TYPE_3D_FIX,
			Lat:               473977420, // 47.3977420° (Zurich)
			Lon:               85512630,  // 8.5512630°
			Alt:               408000,    // 408m MSL
			Eph:               150,       // HDOP 1.50
			Epv:               200,       // VDOP 2.00
			Vel:               500,       // 5 m/s
			Cog:               9000,      // 90 degrees
			SatellitesVisible: 12,
			AltEllipsoid:      410000, // 410m ellipsoid
			HAcc:              1500,   // 1.5m
			VAcc:              2500,   // 2.5m
			VelAcc:            500,    // 0.5 m/s
			HdgAcc:            10000,  // 0.1 degrees
			Yaw:               18000,  // 180 degrees
		}

		result := GpsRawIntToProtobuf(msg)

		if result.TimeUsec != 1234567890 {
			t.Errorf("TimeUsec: got %v, want 1234567890", result.TimeUsec)
		}
		if result.FixType != flightpath.GpsFixType_GPS_FIX_TYPE_3D_FIX {
			t.Errorf("FixType: got %v, want GPS_FIX_TYPE_3D_FIX", result.FixType)
		}
		if result.Lat != 473977420 {
			t.Errorf("Lat: got %v, want 473977420", result.Lat)
		}
		if result.Lon != 85512630 {
			t.Errorf("Lon: got %v, want 85512630", result.Lon)
		}
		if result.Alt != 408000 {
			t.Errorf("Alt: got %v, want 408000", result.Alt)
		}
		if result.Eph != 150 {
			t.Errorf("Eph: got %v, want 150", result.Eph)
		}
		if result.Epv != 200 {
			t.Errorf("Epv: got %v, want 200", result.Epv)
		}
		if result.Vel != 500 {
			t.Errorf("Vel: got %v, want 500", result.Vel)
		}
		if result.Cog != 9000 {
			t.Errorf("Cog: got %v, want 9000", result.Cog)
		}
		if result.SatellitesVisible != 12 {
			t.Errorf("SatellitesVisible: got %v, want 12", result.SatellitesVisible)
		}
		if result.AltEllipsoid != 410000 {
			t.Errorf("AltEllipsoid: got %v, want 410000", result.AltEllipsoid)
		}
		if result.HAcc != 1500 {
			t.Errorf("HAcc: got %v, want 1500", result.HAcc)
		}
		if result.VAcc != 2500 {
			t.Errorf("VAcc: got %v, want 2500", result.VAcc)
		}
		if result.VelAcc != 500 {
			t.Errorf("VelAcc: got %v, want 500", result.VelAcc)
		}
		if result.HdgAcc != 10000 {
			t.Errorf("HdgAcc: got %v, want 10000", result.HdgAcc)
		}
		if result.Yaw != 18000 {
			t.Errorf("Yaw: got %v, want 18000", result.Yaw)
		}
	})

	t.Run("no GPS fix", func(t *testing.T) {
		msg := &common.MessageGpsRawInt{
			TimeUsec:          0,
			FixType:           common.GPS_FIX_TYPE_NO_GPS,
			Lat:               0,
			Lon:               0,
			Alt:               0,
			Eph:               65535, // unknown
			Epv:               65535, // unknown
			Vel:               65535, // unknown
			Cog:               65535, // unknown
			SatellitesVisible: 255,   // unknown
		}

		result := GpsRawIntToProtobuf(msg)

		if result.FixType != flightpath.GpsFixType_GPS_FIX_TYPE_NO_GPS {
			t.Errorf("FixType: got %v, want GPS_FIX_TYPE_NO_GPS", result.FixType)
		}
		if result.SatellitesVisible != 255 {
			t.Errorf("SatellitesVisible: got %v, want 255", result.SatellitesVisible)
		}
	})

	t.Run("RTK fixed", func(t *testing.T) {
		msg := &common.MessageGpsRawInt{
			TimeUsec:          9999999999,
			FixType:           common.GPS_FIX_TYPE_RTK_FIXED,
			Lat:               377749300, // San Francisco
			Lon:               -1224194200,
			Alt:               50000,
			Eph:               10, // HDOP 0.10 (RTK precision)
			Epv:               15,
			Vel:               0,
			Cog:               0,
			SatellitesVisible: 18,
			HAcc:              10, // 10mm precision
			VAcc:              15,
		}

		result := GpsRawIntToProtobuf(msg)

		if result.FixType != flightpath.GpsFixType_GPS_FIX_TYPE_RTK_FIXED {
			t.Errorf("FixType: got %v, want GPS_FIX_TYPE_RTK_FIXED", result.FixType)
		}
		if result.HAcc != 10 {
			t.Errorf("HAcc: got %v, want 10", result.HAcc)
		}
	})
}
