package game

// ItemsTurn handles item updates each turn (cooldowns, etc.)
func (gs *State) ItemsTurn() {
	for i := range gs.Items {
		item := &gs.Items[i]
		// Decrement cooldown for bomb refill generators
		if item.Cooldown > 0 {
			item.Cooldown--
		}
	}
}
