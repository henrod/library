package shelves

import (
	"context"
	"fmt"
	"time"

	"github.com/Henrod/library/domain/entities"
	"github.com/Henrod/library/domain/errors"
	"go.uber.org/zap"
)

type CreateShelfDomain struct {
	gateway CreateShelfGateway
	log     *zap.SugaredLogger

	// This is not scalable, but works for studying purposes.
	pendingShelves map[string]*shelfCreationStatus
}

func NewCreateShelfDomain(
	log *zap.SugaredLogger,
	gateway CreateShelfGateway,
) *CreateShelfDomain {
	domain := &CreateShelfDomain{
		log:            log,
		gateway:        gateway,
		pendingShelves: make(map[string]*shelfCreationStatus),
	}

	go domain.cleanUp()

	return domain
}

type CreateShelfGateway interface {
	CreateShelf(ctx context.Context, shelf *entities.Shelf) (*entities.Shelf, error)
}

type shelfCreationStatus struct {
	stage      int
	err        error
	finished   bool
	finishTime time.Time
}

var shelfCreationStages = []string{
	"PLANTING_TREE",
	"CUTTING_TREE",
	"BUILDING_SHELF",
	"INSTALLING_SHELF",
	"FINISHED_SHELF",
}

const (
	stageTime                 = 1 * time.Second
	failedStageExpirationTime = 10 * time.Second
)

// StartCreateShelfOperation starts a long-running operation to create a shelf.
// After starting it, retrieve the creation status in the Operation API: `GET /operations/shelves/{shelf_name}`
// If the operation fails, its reason is retrievable until expiration time.
func (c *CreateShelfDomain) StartCreateShelfOperation(
	_ context.Context,
	inputShelf *entities.Shelf,
) (*entities.Operation, error) {
	if _, ok := c.pendingShelves[inputShelf.Name]; ok {
		return nil, errors.AlreadyExistsError{
			Details: fmt.Sprintf("create shelf %s operation already exists", inputShelf.Name),
		}
	}

	c.pendingShelves[inputShelf.Name] = new(shelfCreationStatus)

	go c.createShelf(context.Background(), inputShelf)

	operation, _ := c.GetOperation(inputShelf.Name)

	return operation, nil
}

func (c *CreateShelfDomain) cleanUp() {
	for range time.NewTicker(time.Minute).C {
		now := time.Now()
		for shelfName, status := range c.pendingShelves {
			isExpired := now.After(status.finishTime.Add(failedStageExpirationTime))
			if status.finished && isExpired {
				delete(c.pendingShelves, shelfName)
			}
		}
	}
}

func (c *CreateShelfDomain) createShelf(ctx context.Context, inputShelf *entities.Shelf) {
	ticker := time.NewTicker(stageTime)
	log := c.log.With(zap.String("shelf", inputShelf.Name))
	log.Info("started creating shelf")

	for range ticker.C {
		finished, err := c.executeStage(ctx, inputShelf)
		status := c.pendingShelves[inputShelf.Name]

		if err != nil {
			log.With(zap.Error(err)).Error("failure creating shelf")
			status.err = err
			status.finishTime = time.Now()
			status.finished = true

			return
		} else if finished {
			log.Info("finished creating shelf")
			status.finishTime = time.Now()
			status.finished = true

			return
		}
	}
}

func (c *CreateShelfDomain) executeStage(ctx context.Context, inputShelf *entities.Shelf) (finished bool, err error) {
	status := c.pendingShelves[inputShelf.Name]
	status.stage++
	register := status.stage >= len(shelfCreationStages)-1
	if !register {
		return false, nil
	}

	shelf, err := c.gateway.CreateShelf(ctx, inputShelf)
	if err != nil {
		return false, fmt.Errorf("failed to create book in gateway: %w", err)
	}

	if shelf == nil {
		return false, errors.AlreadyExistsError{
			Details: fmt.Sprintf("shelf %s already exists", inputShelf.Name),
		}
	}

	return true, nil
}

func (c *CreateShelfDomain) GetOperation(shelfName string) (*entities.Operation, error) {
	status, ok := c.pendingShelves[shelfName]
	if !ok {
		return nil, errors.NotFoundError{
			Details: fmt.Sprintf("operation for shelf not found: %s", shelfName),
		}
	}

	return &entities.Operation{
		Name:       fmt.Sprintf("operations/shelves/%s", shelfName),
		Stage:      shelfCreationStages[status.stage],
		Percentage: status.stage * 100 / (len(shelfCreationStages) - 1),
		Error:      status.err,
	}, nil
}
