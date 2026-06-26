package inlinegenericservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	genericservice "go.viam.com/rdk/services/generic"
)

// IMPORTANT: Do not change the model name triplet, the struct names, or the public function names below.
// The platform uses these auto-generated values to identify your module.
// Changing them will break your inline module.
var (
	GenericService   = resource.NewModel("jmatthviam", "3f9f4aed-c8ed-4f52-a55b-7fdfa33ef3fc", "generic-service")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterService(genericservice.API, GenericService,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newGenericService,
		},
	)
}

type Config struct {
	Board string `json:"board"`
	Servo string `json:"servo"`
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns three values:
//  1. Required dependencies: other resources that must exist for this resource to work.
//  2. Optional dependencies: other resources that may exist but are not required.
//  3. An error if any Config fields are missing or invalid.
//
// The `path` parameter indicates
// where this resource appears in the machine's JSON configuration
// (for example, "components.0"). You can use it in error messages
// to indicate which resource has a problem.
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	requiredDeps := []string{}
	optionalDeps := []string{}

	// board
	if cfg.Board == "" {
		return nil, nil, fmt.Errorf(
			"%s: attribute 'board' (non-empty string) is required",
			path,
		)
	}
	requiredDeps = append(requiredDeps, cfg.Board)

	// servo
	if cfg.Servo == "" {
		return nil, nil, fmt.Errorf(
			"%s: attribute 'servo' (non-empty string) is required",
			path,
		)
	}
	requiredDeps = append(requiredDeps, cfg.Servo)

	return requiredDeps, optionalDeps, nil
}

type genericService struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	board board.Board
	servo servo.Servo
}

// Status implements [resource.Resource].
func (s *genericService) Status(ctx context.Context) (map[string]any, error) {
	return map[string]any{}, nil
}

func newGenericService(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewGenericService(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewGenericService(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	boardDep, err := board.FromProvider(deps, conf.Board)
	if err != nil {
		return nil, err
	}

	servoDep, err := servo.FromProvider(deps, conf.Servo)
	if err != nil {
		return nil, err
	}

	s := &genericService{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		board:      boardDep,
		servo:      servoDep,
	}

	return s, nil
}

func (s *genericService) Name() resource.Name {
	return s.name
}

func (s *genericService) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	s.logger.Info("Moving servo back and forth between 105 and 180 degrees")

	// Move to 105 degrees
	delay := 250 * time.Millisecond
	for range 3 {
		err := s.servo.Move(ctx, 105, nil)
		if err != nil {
			s.logger.Errorf("Failed to move servo to 105 degrees: %v", err)
			return nil, fmt.Errorf("failed to move servo to 105 degrees: %w", err)
		}
		s.logger.Info("Servo moved to 105 degrees")
		time.Sleep(delay)

		// Move to 180 degrees
		err = s.servo.Move(ctx, 180, nil)
		if err != nil {
			s.logger.Errorf("Failed to move servo to 180 degrees: %v", err)
			return nil, fmt.Errorf("failed to move servo to 180 degrees: %w", err)
		}
		s.logger.Info("Servo moved to 180 degrees")
		time.Sleep(delay)

		// Move back to 105 degrees
		err = s.servo.Move(ctx, 105, nil)
		if err != nil {
			s.logger.Errorf("Failed to move servo back to 105 degrees: %v", err)
			return nil, fmt.Errorf("failed to move servo back to 105 degrees: %w", err)
		}
		s.logger.Info("Servo moved back to 105 degrees")
	}

	return map[string]interface{}{"result": "servo moved between 105 and 180 degrees"}, nil
}

func (s *genericService) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
