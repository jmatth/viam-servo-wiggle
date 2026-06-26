package servowiggle

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/mo"

	generic "go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

var (
	Servo            = resource.NewModel("jmatthviam", "servo-wiggle", "servo")
	errUnimplemented = errors.New("unimplemented")
)

func init() {
	resource.RegisterComponent(generic.API, Servo,
		resource.Registration[resource.Resource, *Config]{
			Constructor: newServoWiggleServo,
		},
	)
}

type Config struct {
	Servo      string `json:"servo"`
	Repeat     int    `json:"repeat"`
	DelayMS    int    `json:"delay_ms"`
	StartAngle uint32 `json:"start_angle"`
	EndAngle   uint32 `json:"end_angle"`
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
	if cfg.StartAngle < 0 || cfg.StartAngle > 360 {
		return nil, nil, errors.New("Start angle must be between 0 and 360")
	}
	if cfg.EndAngle < 0 || cfg.EndAngle > 360 {
		return nil, nil, errors.New("Start angle must be between 0 and 360")
	}
	if cfg.DelayMS < 1 {
		return nil, nil, errors.New("Delay must be positive")
	}
	if cfg.Repeat < 1 {
		return nil, nil, errors.New("Repeat must be positive")
	}
	return []string{cfg.Servo}, nil, nil
}

type servoWiggleServo struct {
	resource.AlwaysRebuild
	resource.Named

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()

	upstream servo.Servo
}

func newServoWiggleServo(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (resource.Resource, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewServo(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewServo(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (resource.Resource, error) {
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &servoWiggleServo{
		name:       name,
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		upstream:   mo.TupleToResult(servo.FromProvider(deps, conf.Servo)).MustGet(),
	}
	return s, nil
}

func (s *servoWiggleServo) Name() resource.Name {
	return s.name
}

func (s *servoWiggleServo) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	delay := time.Duration(s.cfg.DelayMS) * time.Millisecond
	start := s.cfg.StartAngle
	end := s.cfg.EndAngle

	for range s.cfg.Repeat {
		s.upstream.Move(ctx, start, nil)
		time.Sleep(delay)
		s.upstream.Move(ctx, end, nil)
		time.Sleep(delay)
	}
	s.upstream.Move(ctx, 105, nil)
	return nil, nil
}

func (s *servoWiggleServo) Status(ctx context.Context) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *servoWiggleServo) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *servoWiggleServo) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}
