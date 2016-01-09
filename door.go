package gogogarage

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/convertors/mcp3008"
	_ "github.com/kidoman/embd/host/all"
)

//DoorSettings holds the settings on which ADC to use as well as max / min values
type DoorSettings struct {
	ADCChannel        int  //ADC Channel, 0-8
	ADCClosed         int  //ADC when door is closed
	ADCOpened         int  //ADC when door is opened
	Mode, Channel     byte //SPIBus settings
	Speed, Bpw, Delay int  //more SPIBus settings
}

//PercentOpen returns the percentage the door is opened
func (s DoorSettings) PercentOpen(adc int) int {
	if s.ADCOpened == s.ADCClosed {
		panic("ADCOpened cannot be equal to ADCClosed")
	}
	//count_to_value = cnt* 100.0 / ( 3630 - 1083) + 100.0 * 1083 / (1083 - 3630)
	f := float64(adc)*100.0/float64(s.ADCOpened-s.ADCClosed) + 100.0*float64(s.ADCClosed)/float64(s.ADCClosed-s.ADCOpened)
	return int(f)
}

//Door is handler for locating the garage door
type Door struct {
	SPIDev  string           //something like /dev/spi0.0
	Port    int              //ADC port
	Cfg     DoorSettings     //Door Settings
	Spibus  *embd.SPIBus     //
	MPC3008 *mcp3008.MCP3008 // s
}

//PercentOpen from 0 (Closed) to 100 (Open)
func (d *Door) PercentOpen() (int, error) {
	val, err := d.MPC3008.AnalogValueAt(d.Cfg.ADCChannel)
	if err != nil {
		return -1, err
	}
	return d.Cfg.PercentOpen(val), nil
}

//Init door
func (d *Door) Init(cfg DoorSettings) {
	d.Cfg = cfg
	if err := embd.InitSPI(); err != nil {
		panic(err)
	}

	spiBus := embd.NewSPIBus(d.Cfg.Mode, d.Cfg.Channel, d.Cfg.Speed, d.Cfg.Bpw, d.Cfg.Delay)
	// defer embd.CloseSPI()
	// defer spiBus.Close()

	d.MPC3008 = mcp3008.New(mcp3008.SingleMode, spiBus)
}
