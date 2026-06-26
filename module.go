package servowiggle

import (
  generic "go.viam.com/rdk/components/generic"
  "context"
commonpb "go.viam.com/api/common/v1"
genericpb "go.viam.com/api/component/generic/v1"
"go.viam.com/utils/protoutils"
"go.viam.com/utils/rpc"
"go.viam.com/rdk/logging"
rprotoutils "go.viam.com/rdk/protoutils"
"go.viam.com/rdk/referenceframe"
"go.viam.com/rdk/resource"
"go.viam.com/rdk/spatialmath"
)

var (
	Servo = resource.NewModel("jmatthviam", "servo-wiggle", "servo")
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
	/*
	Put config attributes here. There should be public/exported fields
	with a `json` parameter at the end of each attribute.

	Example config struct:
		type Config struct {
			Pin   string `json:"pin"`
			Board string `json:"board"`
			MinDeg *float64 `json:"min_angle_deg,omitempty"`
		}

	If your model does not need a config, replace *Config in the init
	function with resource.NoNativeConfig
	*/
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns three values:
//   1. Required dependencies: other resources that must exist for this resource to work.
//   2. Optional dependencies: other resources that may exist but are not required.
//   3. An error if any Config fields are missing or invalid.
//
// The `path` parameter indicates
// where this resource appears in the machine's JSON configuration
// (for example, "components.0"). You can use it in error messages 
// to indicate which resource has a problem.
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	 return nil, nil, nil
}

type servoWiggleServo struct {
	resource.AlwaysRebuild
	resource.Named

	name   resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
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
	}
	return s, nil
}

func (s *servoWiggleServo) Name() resource.Name {
	return s.name
}

func (s *servoWiggleServo) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
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
