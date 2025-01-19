package main

import (
	"context"
	// "encoding/binary"
	// "encoding/json"
	"fmt"
	// "log"
	// "math"
	// "os/exec"
	"sync"
	// "time"

	"go.einride.tech/can"
	"wt-race-dash/pkg/hondata"

	"go.einride.tech/can/pkg/socketcan"
	// "wt-race-dash/pkg/datalogging"
)

type CANInterface interface {
	ProcessCANFrame(frameId uint32, data can.Data, wg sync.WaitGroup, isDatalogging bool) []byte
}

var (
	isDatalogging = false
	canInterface CANInterface
)

func (wsConn *MySocket) HandleCanBusData() {
	// ---------- CANBus data ----------
	canConn, err := socketcan.DialContext(context.Background(), "can", appSettings.CanChannel)
	if err != nil {
		// "Failed to connect to CAN channel"
		fmt.Println("SocketCAN Connection Error: ", err)
		fmt.Println("==========================================")
		fmt.Printf("SocketCAN Connection Error: %v", err)
		fmt.Println("==========================================")
	}
	defer canConn.Close()
	canRecv := socketcan.NewReceiver(canConn)

	switch (appSettings.CarOrEcu) {
	case "hondata":
		canInterface = &hondata.CANFrameHandler{
			FrameMisc: hondata.CANFrameMisc{ Type: 5 },
			Frame660: hondata.CANFrame660{ Type: 1 },
			Frame661: hondata.CANFrame661{ Type: 1 },
			Frame662: hondata.CANFrame662{ Type: 1 },
			Frame664: hondata.CANFrame664{ Type: 1 },
			Frame667: hondata.CANFrame667{ Type: 1 },	
		}
	}

	for canRecv.Receive() {
		frame := canRecv.Frame()
	
		jsonData := canInterface.ProcessCANFrame(frame.ID, frame.Data, wg, isDatalogging)
		
		if ((frame.ID == 67 || frame.ID == 103) && canInterface.(*hondata.CANFrameHandler).FrameMisc.ChangePage) {
			canInterface.(*hondata.CANFrameHandler).FrameMisc.ChangePage = false
		}

		if jsonData != nil {
			wsConn.writeToClient(int8(frame.ID), jsonData)
		}
  }
}