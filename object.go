package gamelogic

type SensorList []Sensor
type ControllerList []Controller
type ActuatorList []Actuator
type ObjectList []Object

type Object interface {
	// friends: StateActuator, Actuator, Controller

	Name() string

	Controllers() ControllerList
	Sensors() SensorList
	Actuators() ActuatorList

	ActiveActuators() ActuatorList
	ActivateControllers() ControllerList

	AddSensor(Sensor)
	AddController(Controller)
	AddActuator(Actuator)
	// Reserve ---

	RegisterActuator(Actuator)
	UnregisterActuator(Actuator)
	RegisterObject(Object)
	UnregisterObject(Object)

	/**
	 * UnlinkObject(...)
	 * this object is informed that one of the object to which it holds a reference is deleted
	 * returns true if there was indeed a reference.
	 */

	//virtual bool UnlinkObject(SCA_IObject* clientobj) { return false; }

	FindSensor(string) Sensor
	FindActuator(string) Actuator
	FindController(string) Controller

	//void SetCurrentTime(float currentTime) {}

	//virtual void ReParentLogic();

	/**
	 * Set whether or not to ignore activity culling requests
	 */
	/*void SetIgnoreActivityCulling(bool b)
	{
		m_ignore_activity_culling = b;
	}*/

	/**
	 * Set whether or not this object wants to ignore activity culling
	 * requests
	 */
	/*bool GetIgnoreActivityCulling()
	{
		return m_ignore_activity_culling;
	}*/

	// Suspend all progress.
	Suspend()

	// Resume progress
	Resume()

	// Set init state
	//void SetInitState(unsigned int initState) { m_initState = initState; }

	// initialize the state when object is created
	//void ResetState(void) { SetState(m_initState); }

	// Set the object state
	SetState(state uint)

	// Get the object state
	//unsigned int GetState(void)	{ return m_state; }

	//	const class MT_Point3&	ConvertPythonPylist(PyObject *pylist);

	/*
		virtual int GetGameObjectType() {return -1;}

		typedef enum ObjectTypes {
			OBJ_ARMATURE=0,
			OBJ_CAMERA=1,
			OBJ_LIGHT=2,
		} ObjectTypes;
	*/

}

type object struct {
	sensors     SensorList
	controllers ControllerList
	actuators   ActuatorList

	registeredActuators ActuatorList // actuators that use a pointer to this object
	registeredObjects   ObjectList   // objects that hold reference to this object

	/*
		SG_Dlist: element of objects with active actuators
		          Head: SCA_LogicManager::m_activeActuators
		SG_QList: Head of active actuators list on this object
		          Elements: SCA_IActuator
	*/
	activeActuators ActuatorList

	/*
		SG_Dlist: element of list os lists with active controllers
		          Head: SCA_LogicManager::m_activeControllers
		SG_QList: Head of active controller list on this object
		          Elements: SCA_IController
	*/
	activeControllers ControllerList

	/*
	   SG_Dlist: element of list of lists of active controllers
	             Head: SCA_LogicManager::m_activeControllers
	   SG_QList: Head of active bookmarked controller list globally
	             Elements: SCA_IController with bookmark option
	*/

	//static activeBookmarkedControllers

	// static class MT_Point3 m_sDummy;

	// Ignore activity culling requests?
	ignore_activity_culling bool

	initState  uint
	state      uint
	firstState interface{}
	suspended  bool

	name string
}

func NewObject() Object {
	return &object{
		initState:  0,
		state:      0,
		firstState: nil,
		suspended:  false,
	}
}

func (obj *object) Name() string        { return obj.name }
func (obj *object) SetName(name string) { obj.name = name }

func (obj *object) Controllers() ControllerList         { return obj.controllers }
func (obj *object) Sensors() SensorList                 { return obj.sensors }
func (obj *object) Actuators() ActuatorList             { return obj.actuators }
func (obj *object) ActiveActuators() ActuatorList       { return obj.activeActuators }
func (obj *object) ActivateControllers() ControllerList { return obj.activeControllers }

func (obj *object) AddSensor(act Sensor) {
	obj.sensors = append(obj.sensors, act)
}
func (obj *object) AddController(act Controller) {
	obj.controllers = append(obj.controllers, act)
}
func (obj *object) AddActuator(act Actuator) {
	obj.actuators = append(obj.actuators, act)
}

func (obj *object) RegisterActuator(act Actuator) {
	obj.registeredActuators = append(obj.registeredActuators, act)
}
func (obj *object) RegisterObject(act Object) {
	obj.registeredObjects = append(obj.registeredObjects, act)
}

func (obj *object) UnregisterActuator(act Actuator) {
	for i, ita := range obj.registeredActuators {
		if ita == act {
			obj.registeredActuators = append(obj.registeredActuators[:i], obj.registeredActuators[i+1:]...)
		}
	}
}
func (obj *object) UnregisterObject(act Object) {
	for i, ita := range obj.registeredObjects {
		if ita == act {
			obj.registeredObjects = append(obj.registeredObjects[:i], obj.registeredObjects[i+1:]...)
		}
	}
}

func (obj *object) ReParentLogic() {
	// TODO refactor it

	//oldactuators := obj.Actuators()
	for act, ita := range obj.actuators {
		newactuator := ita.Replica()
		newactuator.ReParent(obj)
		// actuators are initially not connected to any controller
		newactuator.SetActive(false)
		newactuator.ClrLink() // XXX wat this?
		obj.actuators[act] = newactuator
	}

	//oldcontrollers := obj.Controllers()
	for con, itc := range obj.controllers {
		newcontroller := itc.Replica()
		newcontroller.ReParent(obj)
		newcontroller.SetActive(false)
		obj.controllers[con] = newcontroller
	}

	// convert sensors last so that actuators are already available for Actuator sensor
	//oldsensors := obj.Sensors()
	for sen, its := range obj.sensors {
		newsensor := its.Replica()
		newsensor.ReParent(obj)
		newsensor.SetActive(false)
		// sensors are initially not connected to any controller
		newsensor.ClrLink() // XXX wat this?
		obj.sensors[sen] = newsensor
	}

	// a new object cannot be client of any actuator
	obj.registeredActuators = ActuatorList{}
	obj.registeredObjects = ObjectList{}

	//obj.actuators = oldactuators
	//obj.controllers = oldcontrollers
	//obj.sensors = oldsensors
}

func (obj *object) FindSensor(name string) Sensor {
	for _, val := range obj.sensors {
		if val.Name() == name {
			return val
		}
	}
	return nil
}
func (obj *object) FindController(name string) Controller {
	for _, val := range obj.controllers {
		if val.Name() == name {
			return val
		}
	}
	return nil
}
func (obj *object) FindActuator(name string) Actuator {
	for _, val := range obj.actuators {
		if val.Name() == name {
			return val
		}
	}
	return nil
}

func (obj *object) Suspend() {
	if (!obj.ignore_activity_culling) && (!obj.suspended) {
		obj.suspended = true
		/* flag suspend for all sensors */
		for _, sensor := range obj.sensors {
			sensor.Suspend()
		}
	}
}

func (obj *object) Resume() {
	if obj.suspended {
		obj.suspended = false
		/* flag suspend for all sensors */
		for _, sensor := range obj.sensors {
			sensor.Resume()
		}
	}
}

func (obj *object) SetState(state uint) {
	// we will update the state in two steps:
	// 1) set the new state bits that are 1
	// 2) clr the new state bits that are 0
	// This to ensure continuity if a sensor is attached to two states
	// that are switching state: no need to deactive and reactive the sensor

	tmpstate := obj.state | state
	if tmpstate != obj.state {
		// update the status of the controllers
		for _, ctrl := range obj.controllers {
			ctrl.ApplyState(tmpstate)
		}
	}
	obj.state = state
	if obj.state != tmpstate {
		for _, ctrl := range obj.controllers {
			ctrl.ApplyState(obj.state)
		}
	}
}
