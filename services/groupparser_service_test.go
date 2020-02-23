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