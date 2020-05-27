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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
			},
			models.Group {
				Sets: getAddressValueInt(4),
				Repetitions: getAddressValueInt(8),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(2),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
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
				Sets: getAddressValueInt(1),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(50),
				Sequence: getAddressValueInt64(1),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
			},
			models.Group {
				Sets: getAddressValueInt(4),
				Repetitions: getAddressValueInt(8),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(2),
			},
			models.Group {
				Sets: getAddressValueInt(5),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(30),
				Sequence: getAddressValueInt64(3),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(4),
				Repetitions: getAddressValueInt(8),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(2),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(5),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(30),
				Sequence: getAddressValueInt64(3),
				Variation: getAddressValueString("variation"),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(1),
				Repetitions: getAddressValueInt(8),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(2),
				Variation: getAddressValueString("variation"),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(10),
				Weight: getAddressValueFloat32(20),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(1),
				Repetitions: getAddressValueInt(8),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(2),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(5),
				Repetitions: getAddressValueInt(7),
				Weight: getAddressValueFloat32(40),
				Sequence: getAddressValueInt64(3),
				Variation: getAddressValueString("variation"),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(5),
				Sequence: getAddressValueInt64(1),
			},
			models.Group {
				Sets: getAddressValueInt(8),
				Repetitions: getAddressValueInt(40),
				Sequence: getAddressValueInt64(2),
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
				Sets: getAddressValueInt(3),
				Repetitions: getAddressValueInt(5),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
			},
			models.Group {
				Sets: getAddressValueInt(8),
				Repetitions: getAddressValueInt(40),
				Sequence: getAddressValueInt64(2),
				Variation: getAddressValueString("variation"),
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
				Sets: getAddressValueInt(9),
				Repetitions: getAddressValueInt(5),
				Sequence: getAddressValueInt64(1),
				Variation: getAddressValueString("variation"),
			},
		}
		assert.ElementsMatch(expectedGroups, groups)
	})
}

func getAddressValueInt(val int) *int {
	return &val
}

func getAddressValueFloat32(val float32) *float32 {
	return &val
}

func getAddressValueInt64(val int64) *int64 {
	return &val
}
func getAddressValueString(val string) *string{
	return &val
}