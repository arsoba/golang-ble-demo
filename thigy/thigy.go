package thigy

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ble/ble"
	goble "github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
)

type ThigyState struct {
	Temp [2]int
}

// Thigy struct
type Thigy struct {
	gattDevice      ble.Device
	client          goble.Client
	characteristics map[string]map[string]*goble.Characteristic
	done            chan bool
	Msg             chan BleMsg
	ThigyState      ThigyState
	Connected       bool
}

type BleMsg struct {
	ServiceUUID        string
	CharacteristicUUID string
	B                  []byte
}

func NewThigy() *Thigy {
	log.SetOutput(os.Stdout)
	thigy := &Thigy{
		nil,
		nil,
		make(map[string]map[string]*goble.Characteristic),
		make(chan bool),
		make(chan BleMsg),
		ThigyState{
			[2]int{0, 0},
		},
		false,
	}

	return thigy
}

func (t *Thigy) Init() (err error) {
	if t.gattDevice, err = linux.NewDeviceWithName("hci0"); err != nil {
		fmt.Println(err)
		os.Exit(1)
		return err
	}
	ble.SetDefaultDevice(t.gattDevice)
	t.startDiscover()
	return nil
}

func (t *Thigy) startDiscover() {
	var err error

	id := "CF:AA:13:A1:5C:A5"
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Addr().String()) == id || id == "*"
	}

	fmt.Println("Scanning...")
	t.client, err = ble.Connect(context.Background(), filter)
	if err != nil {
		fmt.Printf("can't connect : %s", err)
		os.Exit(1)
	}
	t.Connected = true
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", t.client.Addr(), t.client.Name())

	go func() {
		<-t.client.Disconnected()
		fmt.Printf("[ %s ] is disconnected \n", t.client.Addr())
		close(t.done)
	}()

	fmt.Printf("Discovering profile...\n")
	p, err := t.client.DiscoverProfile(true)
	if err != nil {
		fmt.Printf("can't discover profile: %s", err)
		os.Exit(1)
	}

	// Start the exploration.
	t.explore(p)

	<-t.done

	fmt.Println("Done")
}

func (t *Thigy) explore(p *ble.Profile) error {
	go t.writeLoop(p)
	t.readTemp(p)

	return nil
}

func (t *Thigy) readTemp(p *ble.Profile) error {
	characteristicUUID := parseUUID(TempUUID)
	h := func(req []byte) {
		t.ThigyState.Temp = [2]int{int(req[0]), int(req[1])}
	}
	if u := p.Find(ble.NewCharacteristic(characteristicUUID)); u != nil {
		err := t.client.Subscribe(u.(*ble.Characteristic), false, h)
		return errors.Wrap(err, "can't subscribe to characteristic")
	}
	return nil
}

func (t *Thigy) writeLoop(p *ble.Profile) {
	for {
		select {
		case msg := <-t.Msg:
			t.writeToCharacteristic(p, msg)
		}
	}
}

func (t *Thigy) writeToCharacteristic(p *ble.Profile, m BleMsg) (err error) {
	characteristicUUID := parseUUID(m.CharacteristicUUID)
	if u := p.Find(ble.NewCharacteristic(characteristicUUID)); u != nil {
		err := t.client.WriteCharacteristic(u.(*ble.Characteristic), m.B, false)
		return errors.Wrap(err, "can't write characteristic")
	}
	return err
}

func parseUUID(uuidString string) goble.UUID {
	uuid := goble.MustParse(uuidString)
	return uuid
}
