package gamelogic

import (
	"log"
)

/**
 * Interface Class for all logic Sensors. Implements
 * pulsemode,pulsefrequency
 * Use of SG_DList element: link sensors to their respective event manager
 *                          Head: SCA_EventManager::m_sensors
 * Use of SG_QList element: not used
 */

type sensor struct {
	LogicBrick
	eventmgr EventManager

	// Pulse positive  pulses?
	pos_pulsemode bool
	// Pulse negative pulses?
	neg_pulsemode bool
	// Repeat frequency in pulse mode.
	pulse_frequency int
	// Number of ticks since the last positive pulse.
	pos_ticks int
	// Number of ticks since the last negative pulse.
	neg_ticks int
	// invert the output signal
	invert bool
	// detect level instead of edge
	level bool
	// tap mode
	tap bool
	// sensor has been reset
	reset bool
	// Sensor must ignore updates?
	suspended bool
	// number of connections to controller
	links int
	// current sensor state
	state bool
	// previous state (for tap option)
	prev_state        bool
	linkedcontrollers ControllerList
}

const (
	ST_NONE = iota
	ST_TOUCH
	ST_NEAR
	ST_RADAR
)

func NewSensor(gameobj Object, eventmgr EventManager) Sensor {
	// TODO
	return &sensor{
		LogicBrick: NewLogicBrick(gameobj),

		links:           0,
		suspended:       false,
		invert:          false,
		level:           false,
		tap:             false,
		reset:           false,
		pos_ticks:       0,
		neg_ticks:       0,
		pos_pulsemode:   false,
		neg_pulsemode:   false,
		pulse_frequency: 0,
		state:           false,
		prev_state:      false,

		eventmgr: eventmgr,
	}
}

type Sensor interface {
	LogicBrick

	//ReParent(parent Object)

	// Because we want sensors to share some behavior, the Activate has
	// an implementation on this level. It requires an evaluate on the lower
	// level of individual sensors. Mapping the old activate()s is easy.
	// The IsPosTrig() also has to change, to keep things consistent.
	/*
		void Activate(class SCA_LogicManager* logicmgr);
		virtual bool Evaluate() = 0;
		virtual bool IsPositiveTrigger();
		virtual void Init();

		virtual CValue* GetReplica()=0;
	*/

	// Set parameters for the pulsing behavior.
	// \param posmode Trigger positive pulses?
	// \param negmode Trigger negative pulses?
	// \param freq    Frequency to use when doing pulsing.
	SetPulseMode(posmode, negmode bool, freq int)

	// Set inversion of pulses on or off
	SetInvert(inv bool)
	// set the level detection on or off
	SetLevel(lvl bool)
	SetTap(tap bool)

	RegisterToManager()
	UnregisterToManager()
	Replace_EventManager(logicmgr *LogicManager)

	LinkToController(controller Controller)
	UnlinkController(controller Controller)
	UnlinkAllControllers()
	ActivateControllers(logicmgr *LogicManager)

	ProcessReplica()

	SensorType() uint

	// Stop sensing for a while.
	Suspend()
	// Is this sensor switched off?
	IsSuspended() bool

	// get the state of the sensor: positive or negative
	State() bool
	// get the previous state of the sensor: positive or negative
	PrevState() bool
	// get the number of ticks since the last positive pulse
	PosTicks() int
	// get the number of ticks since the last negative pulse
	NegTicks() int

	// Resume sensing.
	Resume()

	ClrLink()
	IncLink()
	DecLink()
	IsNoLink() bool
}

func (s *sensor) SensorType() uint { return ST_NONE }
func (s *sensor) State() bool      { return s.state }
func (s *sensor) PrevState() bool  { return s.prev_state }
func (s *sensor) PosTicks() int    { return s.pos_ticks }
func (s *sensor) NegTicks() int    { return s.neg_ticks }

func (s *sensor) ClrLink() { s.links = 0 }
func (s *sensor) IncLink() {
	if s.links == 0 {
		s.RegisterToManager()
	}
	s.links++
}
func (s *sensor) IsNoLink() bool { return s.links == 0 }

///////////////////////////
// Native functions
func (s *sensor) ReParent(parent Object) {
	s.LogicBrick.ReParent(parent)
	// will be done when the sensor is activated
	//m_eventmgr->RegisterSensor(this);
	//this->SetActive(false);
}

func (s *sensor) ProcessReplica() {
	s.LogicBrick.ProcessReplica()
	s.linkedcontrollers = nil
}

func (s *sensor) IsPositiveTrigger() (result bool) {
	if s.eventval {
		result = (s.eventval.Number() != 0.0)
	}
	if s.invert {
		result = !result
	}

	return
}

func (s *sensor) SetPulseMode(posmode, negmode bool, freq int) {
	s.pos_pulsemode = posmode
	s.neg_pulsemode = negmode
	s.pulse_frequency = freq
}

func (s *sensor) SetInvert(inv bool) { s.invert = inv }
func (s *sensor) SetLevel(lvl bool)  { s.level = lvl }
func (s *sensor) SetTap(tap bool)    { s.tap = tap }

func (s *sensor) Number() float64 { return s.State() }

func (s *sensor) Suspend()          { s.suspended = true }
func (s *sensor) IsSuspended() bool { return s.suspended }
func (s *sensor) Resume()           { s.suspended = false }

func (s *sensor) Init() {
	log.Printf("Sensor %s has no init function, please report this bug to Blender.org\n", s.Name())
}

func (s *sensor) DecLink() {
	s.links--
	if s.links < 0 {
		log.Printf("Warning: sensor %s has negative m_links: %d\n", s.Name(), s.links)
		s.links = 0
	}
	if s.links == 0 {
		// sensor is detached from all controllers, remove it from manager
		s.UnregisterToManager()
	}
}

func (s *sensor) RegisterToManager() {
	// sensor is just activated, initialize it
	s.Init() // FIXME
	s.state = false
	s.eventmgr.RegisterSensor(s)
}

func (s *sensor) Replace_EventManager(logicmgr *LogicManager) {
	if s.links != 0 { // true if we're used currently
		s.eventmgr.RemoveSensor(s)
		s.eventmgr = logicmgr.FindEventManager(s.eventmgr.Type())
		s.eventmgr.RegisterSensor(s)
	} else {
		s.eventmgr = logicmgr.FindEventManager(s.eventmgr.Type())
	}
}

func (s *sensor) LinkToController(controller Controller) {
	s.linkedcontrollers = append(s.linkedcontrollers, controller)
}
func (s *sensor) UnlinkController(controller Controller) {
	for i, contit := range s.linkedcontrollers {
		if contit == controller {
			s.linkedcontrollers = append(s.linkedcontrollers[i:], s.linkedcontrollers[:i+1]...)
			return
		}
	}
	log.Printf("Missing link from sensor %s:%s to controller %s:%s\n",
		s.Parent().Name(), s.Name(), controller.Parent().Name(), controller.Name())
}
func (s *sensor) UnlinkAllControllers() {
	for _, contit := range s.linkedcontrollers {
		contit.UnlinkSensor(s)
	}
	s.linkedcontrollers = nil
}

func (s *sensor) UnregisterToManager() {
	s.eventmgr.RemoveSensor(s)
	s.links = 0
}

func (s *sensor) ActivateControllers(logicmgr *LogicManager) {
	for _, contr := range s.linkedcontrollers {
		if contr.IsActive() {
			logicmgr.AddTriggeredController(contr, s)
		}
	}
}

func (s *sensor) Activate(logicmgr *LogicManager) {
	// calculate if a __triggering__ is wanted
	// don't evaluate a sensor that is not connected to any controller
	if s.links != 0 && !s.suspended {
		result := s.Evaluate()
		// store the state for the rest of the logic system
		s.prev_state = s.state
		s.state = s.IsPositiveTrigger()

		if result {
			// the sensor triggered this frame
			if s.state || !s.tap {
				s.ActivateControllers(logicmgr)
				// reset these counters so that pulse are synchronized with transition
				s.pos_ticks = 0
				s.neg_ticks = 0
			} else {
				result = false
			}
		} else {
			// First, the pulsing behavior, if pulse mode is
			// active. It seems something goes wrong if pulse mode is
			// not set :(
			if s.pos_pulsemode {
				s.pos_ticks++
				if s.pos_ticks > s.pulse_frequency {
					if s.state {
						s.ActivateControllers(logicmgr)
						result = true
					}
					s.pos_ticks = 0
				}
			}
			// negative pulse doesn't make sense in tap mode, skip
			if s.neg_pulsemode && !s.tap {
				s.neg_ticks++
				if s.neg_ticks > s.pulse_frequency {
					if !s.state {
						s.ActivateControllers(logicmgr)
						result = true
					}
					s.neg_ticks = 0
				}
			}
		}

		if s.tap {
			// in tap mode: we send always a negative pulse immediately after a positive pulse
			if !result {
				// the sensor did not trigger on this frame
				if s.prev_state {
					// but it triggered on previous frame => send a negative pulse
					s.ActivateControllers(logicmgr)
					result = true
				}
				// in any case, absence of trigger means sensor off
				s.state = false
			}
		}

		if !result && s.level {
			// This level sensor is connected to at least one controller that was just made
			// active but it did not generate an event yet, do it now to those controllers only
			for _, contr := range s.linkedcontrollers {
				if contr.IsJustActivated() {
					logicmgr.AddTriggeredController(contr, s)
				}
			}
		}
	}
}
