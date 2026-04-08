package simulation

import "testing"

func TestRegionalFoodPressureDelaysLocallyDepletedSlotLonger(t *testing.T) {
	world := &World{
		foodSlots: []Food{
			{ID: "food-1", X: 100, Y: 100, Radius: DefaultFoodRadius},
			{ID: "food-2", X: 240, Y: 100, Radius: DefaultFoodRadius},
			{ID: "food-3", X: 620, Y: 100, Radius: DefaultFoodRadius},
		},
		foods: []Food{},
		missingFoodSince: map[string]int64{
			"food-1": 0,
			"food-2": 0,
			"food-3": 0,
		},
	}

	world.regenerateFood(DefaultFoodRegenDelay + 4)
	if _, found := world.foodSlotByID("food-3"); !found {
		t.Fatal("expected deterministic food slot lookup to find food-3")
	}

	hasFood1 := false
	hasFood3 := false
	for _, food := range world.foods {
		if food.ID == "food-1" {
			hasFood1 = true
		}
		if food.ID == "food-3" {
			hasFood3 = true
		}
	}

	if !hasFood3 {
		t.Fatal("expected isolated missing slot to regenerate at the shared global pressure boundary")
	}
	if hasFood1 {
		t.Fatal("expected locally pressured slot to remain missing longer than the isolated slot")
	}

	world.regenerateFood(DefaultFoodRegenDelay + 6)

	hasFood1 = false
	for _, food := range world.foods {
		if food.ID == "food-1" {
			hasFood1 = true
			break
		}
	}
	if !hasFood1 {
		t.Fatal("expected locally pressured slot to regenerate after the additional regional delay")
	}
}

func TestRegionalFoodYieldDropsInLocallyDepletedSlot(t *testing.T) {
	world := &World{
		foodSlots: []Food{
			{ID: "food-1", X: 100, Y: 100, Radius: DefaultFoodRadius},
			{ID: "food-2", X: 220, Y: 100, Radius: DefaultFoodRadius},
			{ID: "food-3", X: 280, Y: 100, Radius: DefaultFoodRadius},
			{ID: "food-4", X: 900, Y: 100, Radius: DefaultFoodRadius},
		},
		missingFoodSince: map[string]int64{
			"food-1": 0,
			"food-2": 0,
		},
		foodGain: DefaultFoodEnergy,
	}

	if gain := world.foodGainForSlot("food-4"); gain != DefaultFoodEnergy {
		t.Fatalf("expected isolated slot to keep baseline food gain %v, got %v", DefaultFoodEnergy, gain)
	}

	expectedReducedGain := DefaultFoodEnergy - 2
	if gain := world.foodGainForSlot("food-3"); gain != expectedReducedGain {
		t.Fatalf("expected regionally depleted slot to yield %v, got %v", expectedReducedGain, gain)
	}
}
