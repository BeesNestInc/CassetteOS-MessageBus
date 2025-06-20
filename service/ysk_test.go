package service_test

import (
	"context"
	"testing"

	"github.com/BeesNestInc/CassetteOS-MessageBus/pkg/ysk"
	"github.com/BeesNestInc/CassetteOS-MessageBus/repository"
	"github.com/BeesNestInc/CassetteOS-MessageBus/service"
	"github.com/BeesNestInc/CassetteOS-MessageBus/utils"
	"gotest.tools/assert"
)

func setup(t *testing.T) (*service.YSKService, func()) {
	repository, err := repository.NewDatabaseRepositoryInMemory()
	assert.NilError(t, err)

	yskService := service.NewYSKService(&repository, nil, nil)
	return yskService, func() {
		repository.Close()
	}
}

func TestInsertAndGetCardList(t *testing.T) {
	yskService, cleanup := setup(t)
	defer cleanup()

	cardList, err := yskService.YskCardList(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, len(cardList), 0)

	cardInsertQueue := []ysk.YSKCard{
		utils.ApplicationInstallProgress, utils.ZimaOSDataStationNotice,
	}

	for _, card := range cardInsertQueue {
		err = yskService.UpsertYSKCard(context.Background(), card)
		assert.NilError(t, err)
	}

	cardList, err = yskService.YskCardList(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, len(cardList), 2)

	for _, card := range cardList {
		if card.Id == utils.ApplicationInstallProgress.Id {
			utils.ApplicationInstallProgress.Created = card.Created
			utils.ApplicationInstallProgress.Updated = card.Updated
			assert.DeepEqual(t, card, utils.ApplicationInstallProgress)
		} else if card.Id == utils.ZimaOSDataStationNotice.Id {
			utils.ZimaOSDataStationNotice.Created = card.Created
			utils.ZimaOSDataStationNotice.Updated = card.Updated
			assert.DeepEqual(t, card, utils.ZimaOSDataStationNotice)
		} else {
			t.Errorf("unexpected card: %v", card)
		}
	}
}

func TestInsertAllTypeCardList(t *testing.T) {
	yskService, cleanup := setup(t)
	defer cleanup()

	cardList, err := yskService.YskCardList(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, len(cardList), 0)

	cardInsertQueue := []ysk.YSKCard{
		utils.ApplicationInstallProgress, utils.ZimaOSDataStationNotice,
		// the notice is short. it didn't be stored
		utils.ApplicationUpdateNotice,
		utils.ApplicationInstallProgress.WithProgress("Installing LinuxServer/Jellyfin", 50), utils.ApplicationInstallProgress.WithProgress("Installing LinuxServer/Jellyfin", 55),
		utils.ApplicationInstallProgress.WithProgress("Installing LinuxServer/Jellyfin", 80), utils.ApplicationInstallProgress.WithProgress("Installing LinuxServer/Jellyfin", 99),
		utils.ApplicationUpdateNotice,
	}

	for _, card := range cardInsertQueue {
		err = yskService.UpsertYSKCard(context.Background(), card)
		assert.NilError(t, err)
	}

	cardList, err = yskService.YskCardList(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, len(cardList), 2)
}
