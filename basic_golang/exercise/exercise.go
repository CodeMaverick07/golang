package exercise

import "fmt"

type Player struct {
	Name      string
	Inventory []Item
}

type Item struct {
	Name string
	Type string
}

func (p *Player) PickUpItem(item Item) {
	p.Inventory = append(p.Inventory, item)
}

func (p *Player) DropItem(itemName string) {
	for idx, value := range p.Inventory {
		if value.Name == itemName {
			p.Inventory = append(p.Inventory[0:idx], p.Inventory[idx+1:]...)
		}
	}
}

func (p *Player) UseItem(itemName string) {
	for _, value := range p.Inventory {
		if value.Name == itemName {
			fmt.Println("item is", value)
		}
	}
}
