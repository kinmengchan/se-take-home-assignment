package main

import (
	"fmt"
	"time"

	"github.com/feedme/order-controller/internal/core"
	"github.com/feedme/order-controller/internal/engine"
)

func ts() string {
	return time.Now().Format("15:04:05")
}

func main() {
	fmt.Println("McDonald's Order Management System - Simulation Results")
	fmt.Println()

	e := engine.New(10 * time.Second)

	e.OnEvent(func(ev engine.Event) {
		switch ev.Type {
		case "order_created":
			o := ev.Payload.(*core.Order)
			fmt.Printf("[%s] Created %s Order #%d - Status: %s\n", ts(), o.Type, o.ID, o.Status)
		case "order_processing":
			p := ev.Payload.(map[string]interface{})
			b := p["bot"].(*core.Bot)
			o := p["order"].(*core.Order)
			fmt.Printf("[%s] Bot #%d picked up %s Order #%d - Status: PROCESSING\n", ts(), b.ID, o.Type, o.ID)
		case "order_completed":
			p := ev.Payload.(map[string]interface{})
			b := p["bot"].(*core.Bot)
			o := p["order"].(*core.Order)
			fmt.Printf("[%s] Bot #%d completed %s Order #%d - Status: COMPLETE (Processing time: 10s)\n", ts(), b.ID, o.Type, o.ID)
		case "bot_created":
			b := ev.Payload.(*core.Bot)
			fmt.Printf("[%s] Bot #%d created - Status: ACTIVE\n", ts(), b.ID)
		case "bot_destroyed":
			b := ev.Payload.(*core.Bot)
			fmt.Printf("[%s] Bot #%d destroyed\n", ts(), b.ID)
		case "bot_idle":
			b := ev.Payload.(*core.Bot)
			fmt.Printf("[%s] Bot #%d is now IDLE - No pending orders\n", ts(), b.ID)
		}
	})

	fmt.Printf("[%s] System initialized with 0 bots\n", ts())

	// Scenario: exercises all requirements
	e.AddOrder(core.Normal) // #1
	time.Sleep(1 * time.Second)
	e.AddOrder(core.VIP)    // #2
	e.AddOrder(core.Normal) // #3
	time.Sleep(1 * time.Second)

	e.AddBot() // Bot #1 - picks up VIP #2
	time.Sleep(1 * time.Second)
	e.AddBot() // Bot #2 - picks up Normal #1

	// Wait for both to complete
	time.Sleep(12 * time.Second)

	// Add another VIP
	e.AddOrder(core.VIP) // #4
	time.Sleep(12 * time.Second)

	// Remove bot #2
	e.RemoveBot()
	time.Sleep(1 * time.Second)

	state := e.State()
	fmt.Println()
	fmt.Println("Final Status:")
	vipCount := 0
	normalCount := 0
	for _, o := range state.CompletedOrders {
		if o.Type == core.VIP {
			vipCount++
		} else {
			normalCount++
		}
	}
	total := vipCount + normalCount
	fmt.Printf("- Total Orders Processed: %d (%d VIP, %d Normal)\n", total, vipCount, normalCount)
	fmt.Printf("- Orders Completed: %d\n", total)
	fmt.Printf("- Active Bots: %d\n", len(state.Bots))
	fmt.Printf("- Pending Orders: %d\n", len(state.PendingOrders))
}
