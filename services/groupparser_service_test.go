package services_test

import (
	"testing"
	"fmt"
	"github.com/cts3njitedu/healthful-heart/services"
	"github.com/cts3njitedu/healthful-heart/models"
	"github.com/stretchr/testify/assert"

)

func TestGetGroupsOneGroup(t *testing.T) {
	t.Run("Testing Get Groups", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsTwoGroups(t *testing.T) {
	t.Run("Testing Get Groups", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20, 4x8x40", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
			},
			models.Group {
				Sets: 4,
				Repetitions: 8,
				Weight: 40,
				Sequence: 2,
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsOneGroupParenthesisAfter(t *testing.T) {
	t.Run("Testing Get Groups With Parenthesis After", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20 (variation)", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsOneGroupRepetitionsAndWeight(t *testing.T) {
	t.Run("Testing Get Groups With Only Repetitions and Weight", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "10x50", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 1,
				Repetitions: 10,
				Weight: 50,
				Sequence: 1,
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsThreeGroups(t *testing.T) {
	t.Run("Testing Get Groups With Three Groups", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20, 4x8x40, 5x10x30", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
			},
			models.Group {
				Sets: 4,
				Repetitions: 8,
				Weight: 40,
				Sequence: 2,
			},
			models.Group {
				Sets: 5,
				Repetitions: 10,
				Weight: 30,
				Sequence: 3,
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsThreeGroupsWithVariation(t *testing.T) {
	t.Run("Testing Get Groups With Three Groups And Variation", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20, 4x8x40, 5x10x30 (variation)", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
				Variation: "variation",
			},
			models.Group {
				Sets: 4,
				Repetitions: 8,
				Weight: 40,
				Sequence: 2,
				Variation: "variation",
			},
			models.Group {
				Sets: 5,
				Repetitions: 10,
				Weight: 30,
				Sequence: 3,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsTwoGroupsWithVariationMixed(t *testing.T) {
	t.Run("Testing Get Groups With Two Groups And Variation Mixed", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20, 8x40 (variation)", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
				Variation: "variation",
			},
			models.Group {
				Sets: 1,
				Repetitions: 8,
				Weight: 40,
				Sequence: 2,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsThreeGroupsWithVariationMixed(t *testing.T) {
	t.Run("Testing Get Groups With Three Groups And Variation Mixed", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("TEST", "3x10x20, 8x40,5x7x40(variation)", "BACK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 10,
				Weight: 20,
				Sequence: 1,
				Variation: "variation",
			},
			models.Group {
				Sets: 1,
				Repetitions: 8,
				Weight: 40,
				Sequence: 2,
				Variation: "variation",
			},
			models.Group {
				Sets: 5,
				Repetitions: 7,
				Weight: 40,
				Sequence: 3,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
	
}

func TestGetGroupsTwoGroupsPullUps(t *testing.T) {
	t.Run("Testing Get Groups With Two Groups Pull Ups", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("BK16", "3x5, 8x40", "BK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 5,
				Sequence: 1,
			},
			models.Group {
				Sets: 8,
				Repetitions: 40,
				Sequence: 2,
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
}


func TestGetGroupsTwoGroupsPullUpsVariation(t *testing.T) {
	t.Run("Testing Get Groups With Two Groups Pull Ups Variation", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("BK16", "3x5, 8x40  (variation)", "BK")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 3,
				Repetitions: 5,
				Sequence: 1,
				Variation: "variation",
			},
			models.Group {
				Sets: 8,
				Repetitions: 40,
				Sequence: 2,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
}

func TestGetGroupsTwoGroupsDipsVariation(t *testing.T) {
	t.Run("Testing Get Groups With Two Groups Dips Variation", func(t *testing.T) {
		assert := assert.New(t)
		service :=  services.NewGroupParserService()
		fmt.Printf("Service type: %T", service)
		groups := service.GetGroups("CH5", "9x5 (variation)", "CH")
		expectedGroups := []models.Group {
			models.Group {
				Sets: 9,
				Repetitions: 5,
				Sequence: 1,
				Variation: "variation",
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
}