package gamelogic

import (
	"log"
)

type Actuator interface {
	LogicBrick

	Activate(ActuatorList)
	Deactivate()
	AddEvent(bool)

	LinkToController(controller Controller)
	UnlinkController(cont Controller)
	UnlinkAllControllers()

	ClrLink()
	IncLink()
	DecLink()
	//IsNoLink() const { return !m_links; }
	//bool IsType(KX_ACTUATOR_TYPE type) { return m_type == type; }
}

/**
 * Use of SG_DList : None
 * Use of SG_QList : element of activated actuator list of their owner
 *                   Head: SCA_IObject::m_activeActuators
 */
type actuator struct {
	LogicBrick
	//friend class SCA_LogicManager;
	_type uint
	links int // number of active links to controllers
	// when 0, the actuator is automatically stopped
	//std::vector<CValue*> m_events;
	posevent bool
	negevent bool

	//std::vector<class SCA_IController*>		m_linkedcontrollers;
	linkedcontrollers ControllerList
}

// protected
func (a *actuator) RemoveAllEvents() {
	a.posevent = false
	a.negevent = false
}

// This class also inherits the default copy constructors
//KX_ACTUATOR_TYPE {
/*const (
	KX_ACT_OBJECT = iota
	KX_ACT_IPO
	KX_ACT_CAMERA
	KX_ACT_SOUND
	KX_ACT_PROPERTY
	KX_ACT_ADD_OBJECT
	KX_ACT_END_OBJECT
	KX_ACT_DYNAMIC
	KX_ACT_REPLACE_MESH
	KX_ACT_TRACKTO
	KX_ACT_CONSTRAINT
	KX_ACT_SCENE
	KX_ACT_RANDOM
	KX_ACT_MESSAGE
	KX_ACT_ACTION
	KX_ACT_CD
	KX_ACT_GAME
	KX_ACT_VISIBILITY
	KX_ACT_2DFILTER
	KX_ACT_PARENT
	KX_ACT_SHAPEACTION
	KX_ACT_STATE
	KX_ACT_ARMATURE
	KX_ACT_STEERING
	KX_ACT_MOUSE
)*/

/*

	SCA_IActuator(SCA_IObject* gameobj, KX_ACTUATOR_TYPE type);

	// UnlinkObject(...)
	// Certain actuator use gameobject pointers (like TractTo actuator)
	// This function can be called when an object is removed to make
	// sure that the actuator will not use it anymore.

	virtual bool UnlinkObject(SCA_IObject* clientobj) { return false; }

	// Update(...)
	// Update the actuator based upon the events received since
	// the last call to Update, the current time and deltatime the
	// time elapsed in this frame ?
	// It is the responsibility of concrete Actuators to clear
	// their event's. This is usually done in the Update() method via
	// a call to RemoveAllEvents()


	virtual bool Update(double curtime, bool frame);
	virtual bool Update();

	// Add an event to an actuator.
*/
//void AddEvent(CValue* event)
func (a *actuator) AddEvent(event bool) {
	if event {
		a.posevent = true
	} else {
		a.negevent = true
	}
}

//virtual void ProcessReplica();

// Return true if all the current events
// are negative. The definition of negative event is
// not immediately clear. But usually refers to key-up events
// or events where no action is required.
func (a *actuator) IsNegativeEvent() bool { return !a.posevent && a.negevent }

// remove this actuator from the list of active actuators
//virtual void Deactivate();
//virtual void Activate(SG_DList& head);
/*
	void	LinkToController(SCA_IController* controller);
	void	UnlinkController(class SCA_IController* cont);
	void	UnlinkAllControllers();
*/

func (a *actuator) ClrLink() { a.links = 0 }
func (a *actuator) IncLink() { a.links++ }

//void DecLink();
func (a *actuator) IsNoLink() bool         { return a.links == 0 }
func (a *actuator) IsType(_type uint) bool { return a._type == _type }

/*
#ifdef WITH_CXX_GUARDEDALLOC
	MEM_CXX_CLASS_ALLOC_FUNCS("GE:SCA_IActuator")
#endif
*/
//};

/////////////////////

func NewActuator(gameobj Object, _type uint) Actuator {
	return &actuator{
		LogicBrick: NewLogicBrick(gameobj),
		_type:      _type,
		links:      0,
		posevent:   false,
		negevent:   false,
	}
}

func (a *actuator) UpdateT(curtime float64, frame bool) bool {
	if frame {
		return a.Update()
	}
	return true
}

func (a *actuator) Update() bool {
	log.Panic("Actuators should override an Update method.")
	return false
}

// TODO
func (a *actuator) Activate(ActuatorList) {}
func (a *actuator) Deactivate()           {}

/*TODO
func (a *actuator)Activate(SG_DList& head) {
	if (QEmpty())
	{
		InsertActiveQList(m_gameobj->m_activeActuators);
		head.AddBack(&m_gameobj->m_activeActuators);
	}
}
*/

/* TODO

// this function is only used to deactivate actuators outside the logic loop
// e.g. when an object is deleted.
void SCA_IActuator::Deactivate()
{
	if (QDelink())
	{
		// the actuator was in the active list
		if (m_gameobj->m_activeActuators.QEmpty())
			// the owner object has no more active actuators, remove it from the global list
			m_gameobj->m_activeActuators.Delink();
	}
}*/

/* FIXME
func (a *actuator) ProcessReplica() {
	a.LogicBrick.ProcessReplica()
	a.RemoveAllEvents()
	a.linkedcontrollers = nil
}
*/

func (a *actuator) DecLink() {
	a.links--
	if a.links < 0 {
		log.Printf("Warning: actuator %s has negative m_links: %d\n", a.Name(), a.links)
		a.links = 0
	}
}

func (a *actuator) LinkToController(controller Controller) {
	a.linkedcontrollers = append(a.linkedcontrollers, controller)
}
func (a *actuator) UnlinkController(controller Controller) {
	for i, contit := range a.linkedcontrollers {
		if contit == controller {
			a.linkedcontrollers = append(a.linkedcontrollers[i:], a.linkedcontrollers[:i+1]...)
			return
		}
	}
	log.Printf("Missing link from actuator %s:%s to controller %s:%s\n",
		a.Parent().Name(), a.Name(), controller.Parent().Name(), controller.Name())
}
func (a *actuator) UnlinkAllControllers() {
	for _, contit := range a.linkedcontrollers {
		contit.UnlinkActuator(a)
	}
	a.linkedcontrollers = nil
}
