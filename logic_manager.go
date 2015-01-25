package gamelogic

/*
const ( //EVENT_MANAGER_TYPE {
	KEYBOARD_EVENTMGR = iota
	MOUSE_EVENTMGR
	ALWAYS_EVENTMGR
	TOUCH_EVENTMGR
	PROPERTY_EVENTMGR
	TIME_EVENTMGR
	RANDOM_EVENTMGR
	RAY_EVENTMGR
	NETWORK_EVENTMGR
	JOY_EVENTMGR
	ACTUATOR_EVENTMGR
	BASIC_EVENTMGR
)*/

type Event interface {
	Release()
	AddRef() Event
}

type EventManager interface {
	//class SCA_LogicManager* m_logicmgr; /* all event manager subclasses use this (other then TimeEventManager) */

	// use a set to speed-up insertion/removal
	//std::set <class SCA_ISensor*>				m_sensors;
	//SG_DList		m_sensors;

	//EVENT_MANAGER_TYPE		m_mgrtype;

	//SCA_EventManager(SCA_LogicManager* logicmgr, EVENT_MANAGER_TYPE mgrtype);

	RemoveSensor(sensor Sensor)
	NextFrame(curtime, fixedtime float64)
	//NextFrame()
	UpdateFrame()
	EndFrame()
	RegisterSensor(sensor Sensor)
	Type() uint
	//SG_DList &GetSensors() { return m_sensors; }

	Replace_LogicManager(logicmgr *LogicManager) //{ m_logicmgr= logicmgr; }

	/*void SCA_EventManager::RegisterSensor(class SCA_ISensor* sensor) {
	  	m_sensors.AddBack(sensor);
	  }
	  void SCA_EventManager::RemoveSensor(class SCA_ISensor* sensor) {
	  	sensor->Delink();
	  }
	  void SCA_EventManager::NextFrame(double curtime, double fixedtime) {
	  	NextFrame();
	  }
	  void SCA_EventManager::NextFrame() {
	  	assert(false); // && "Event managers should override a NextFrame method");
	  }
	  int SCA_EventManager::GetType() {
	  	return (int) m_mgrtype;
	  }
	*/
}

//using namespace std;
//typedef std::list<class SCA_IController*> controllerlist;
//typedef std::map<class SCA_ISensor*,controllerlist > sensormap_t;

/**
 * This manager handles sensor, controllers and actuators.
 * logic executes each frame the following way:
 * find triggering sensors
 * build list of controllers that are triggered by these triggering sensors
 * process all triggered controllers
 * during this phase actuators can be added to the active actuator list
 * process all active actuators
 * clear triggering sensors
 * clear triggered controllers
 * (actuators may be active during a longer timeframe)
 */

type BlendObj interface{}
type GameObject interface{}

type LogicManager struct {
	eventmanagers []EventManager

	// SG_DList: Head of objects having activated actuators
	//           element: SCA_IObject::m_activeActuators
	//SG_DList							m_activeActuators;
	activeActuators ActuatorList
	// SG_DList: Head of objects having activated controllers
	//           element: SCA_IObject::m_activeControllers
	//SG_DList							m_triggeredControllerSet;
	//TODO triggeredControllerSet WTF

	// need to find better way for this
	// also known as FactoryManager...
	mapStringToGameObjects map[string]interface{}
	mapStringToMeshes      map[string]interface{}
	mapStringToActions     map[string]interface{}

	map_gamemeshname_to_blendobj map[string]BlendObj
	map_blendobj_to_gameobj      map[BlendObj]GameObject
}

/*
public:
	SCA_LogicManager();
	virtual ~SCA_LogicManager();

	//void	SetKeyboardManager(SCA_KeyboardManager* keyboardmgr) { m_keyboardmgr=keyboardmgr;}
	void	RegisterEventManager(SCA_EventManager* eventmgr);
	void	RegisterToSensor(SCA_IController* controller,
							 class SCA_ISensor* sensor);
	void	RegisterToActuator(SCA_IController* controller,
							   class SCA_IActuator* actuator);

	void	BeginFrame(double curtime, double fixedtime);
	void	UpdateFrame(double curtime, bool frame);
	void	EndFrame();
*/
func (l *LogicManager) AddActiveActuator(actua Actuator, event bool) {
	actua.SetActive(true)
	actua.Activate(l.activeActuators)
	actua.AddEvent(event)
}

/*

	void	AddTriggeredController(SCA_IController* controller, SCA_ISensor* sensor);
	SCA_EventManager*	FindEventManager(int eventmgrtype);
	vector<class SCA_EventManager*>	GetEventManagers() { return m_eventmanagers; }

	void	RemoveGameObject(const STR_String& gameobjname);

	// remove Logic Bricks from the running logicmanager
	void	RemoveSensor(SCA_ISensor* sensor);
	void	RemoveController(SCA_IController* controller);
	void	RemoveActuator(SCA_IActuator* actuator);


	// for the scripting... needs a FactoryManager later (if we would have time... ;)
	void	RegisterMeshName(const STR_String& meshname,void* mesh);
	void	UnregisterMeshName(const STR_String& meshname,void* mesh);
	CTR_Map<STR_HashedString,void*>&	GetMeshMap() { return m_mapStringToMeshes; }
	CTR_Map<STR_HashedString,void*>&	GetActionMap() { return m_mapStringToActions; }

	void	RegisterActionName(const STR_String& actname,void* action);

	void*	GetActionByName (const STR_String& actname);
	void*	GetMeshByName(const STR_String& meshname);

	void	RegisterGameObjectName(const STR_String& gameobjname,CValue* gameobj);
	class CValue*	GetGameObjectByName(const STR_String& gameobjname);

	void	RegisterGameMeshName(const STR_String& gamemeshname, void* blendobj);
	void*	FindBlendObjByGameMeshName(const STR_String& gamemeshname);

	void	RegisterGameObj(void* blendobj, CValue* gameobj);
	void	UnregisterGameObj(void* blendobj, CValue* gameobj);
	CValue*	FindGameObjByBlendObj(void* blendobj);


#ifdef WITH_CXX_GUARDEDALLOC
	MEM_CXX_CLASS_ALLOC_FUNCS("GE:SCA_LogicManager")
#endif
};
*/

/*
#if 0
// this kind of fixes bug 398 but breakes games, so better leave it out for now.
// a removed object's gameobject (and logicbricks and stuff) didn't get released
// because it was still in the m_mapStringToGameObjects map.
func (l *LogicManager)RemoveGameObject(const STR_String& gameobjname)
{
	int numgameobj = m_mapStringToGameObjects.size();
	for (int i = 0; i < numgameobj; i++)
	{
		CValue** gameobjptr = m_mapStringToGameObjects.at(i);
		assert(gameobjptr);

		if (gameobjptr)
		{
			if ((*gameobjptr).GetName() == gameobjname)
				(*gameobjptr).Release();
		}
	}

	m_mapStringToGameObjects.remove(gameobjname);
}
#endif
*/

func (l *LogicManager) RegisterEventManager(eventmgr EventManager) {
	l.eventmanagers = append(l.eventmanagers, eventmgr)
}

func (l *LogicManager) RegisterGameObjectName(name string, gameobj Object) {
	l.mapStringToGameObjects[name] = gameobj
}
func (l *LogicManager) RegisterGameMeshName(name string, blendobj interface{}) {
	l.map_gamemeshname_to_blendobj[name] = blendobj
}

func (l *LogicManager) RegisterGameObj(blendobj interface{}, gameobj interface{}) {
	l.map_blendobj_to_gameobj[blendobj] = gameobj
}

func (l *LogicManager) UnregisterGameObj(blendobj interface{}, gameobj interface{}) {
	obp, ok := l.map_blendobj_to_gameobj[blendobj]
	if ok && obp == gameobj {
		delete(l.map_blendobj_to_gameobj, blendobj)
	}
}

func (l *LogicManager) GetGameObjectByName(name string) interface{} {
	return l.mapStringToGameObjects[name]
}

func (l *LogicManager) FindGameObjByBlendObj(blendobj interface{}) interface{} {
	return l.map_blendobj_to_gameobj[blendobj]
}
func (l *LogicManager) FindBlendObjByGameMeshName(name string) interface{} {
	return l.map_gamemeshname_to_blendobj[name]
}

func (l *LogicManager) RemoveSensor(sensor Sensor) {
	sensor.UnlinkAllControllers()
	sensor.UnregisterToManager()
}
func (l *LogicManager) RemoveController(controller Controller) {
	controller.UnlinkAllSensors()
	controller.UnlinkAllActuators()
	controller.Deactivate()
}
func (l *LogicManager) RemoveActuator(actuator Actuator) {
	actuator.UnlinkAllControllers()
	actuator.Deactivate()
	actuator.SetActive(false)
}

func (l *LogicManager) RegisterToSensor(controller Controller, sensor Sensor) {
	sensor.LinkToController(controller)
	controller.LinkToSensor(sensor)
}

func (l *LogicManager) RegisterToActuator(controller Controller, actua Actuator) {
	actua.LinkToController(controller)
	controller.LinkToActuator(actua)
}

func (l *LogicManager) BeginFrame(curtime float64, fixedtime float64) {
	for _, ie := range l.eventmanagers {
		ie.NextFrame(curtime, fixedtime)
	}

	/* TODO
	for (SG_QList* obj = (SG_QList*)m_triggeredControllerSet.Remove();
		obj != NULL;
		obj = (SG_QList*)m_triggeredControllerSet.Remove())
	{
		for (SCA_IController* contr = (SCA_IController*)obj.QRemove();
			contr != NULL;
			contr = (SCA_IController*)obj.QRemove())
		{
			contr.Trigger(this);
			contr.ClrJustActivated();
		}
	}
	*/
}

/* TODO

func (l *LogicManager)UpdateFrame(double curtime, bool frame)
{
	for (vector<SCA_EventManager*>::const_iterator ie=m_eventmanagers.begin(); !(ie==m_eventmanagers.end()); ie++)
		(*ie).UpdateFrame();

	SG_DList::iterator<SG_QList> io(m_activeActuators);
	for (io.begin(); !io.end(); )
	{
		SG_QList* ahead = *io;
		// increment now so that we can remove the current element
		++io;
		SG_QList::iterator<SCA_IActuator> ia(*ahead);
		for (ia.begin(); !ia.end();  )
		{
			SCA_IActuator* actua = *ia;
			// increment first to allow removal of inactive actuators.
			++ia;
			if (!actua.Update(curtime, frame))
			{
				// this actuator is not active anymore, remove
				actua.QDelink();
				actua.SetActive(false);
			} else if (actua.IsNoLink())
			{
				// This actuator has no more links but it still active
				// make sure it will get a negative event on next frame to stop it
				// Do this check after Update() rather than before to make sure
				// that all the actuators that are activated at same time than a state
				// actuator have a chance to execute.
				bool event = false;
				actua.RemoveAllEvents();
				actua.AddEvent(event);
			}
		}
		if (ahead.QEmpty())
		{
			// no more active controller, remove from main list
			ahead.Delink();
		}
	}
}
*/

func (l *LogicManager) GetActionByName(name string) interface{} {
	return l.mapStringToActions[name]
}
func (l *LogicManager) GetMeshByName(name string) interface{} {
	return l.mapStringToMeshes[name]
}

func (l *LogicManager) RegisterMeshName(name string, mesh interface{}) {
	l.mapStringToMeshes[name] = mesh
}
func (l *LogicManager) UnregisterMeshName(name string) {
	delete(l.mapStringToMeshes, name)
}

func (l *LogicManager) RegisterActionName(name string, action interface{}) {
	l.mapStringToActions[name] = action
}

func (l *LogicManager) EndFrame() {
	for _, ie := range l.eventmanagers {
		ie.EndFrame()
	}
}

func (l *LogicManager) AddTriggeredController(controller Controller, sensor Sensor) {
	// TODO controller.Activate(m_triggeredControllerSet);
}

func (l *LogicManager) FindEventManager(eventmgrtype uint) EventManager {
	// find an eventmanager of a certain type
	for _, emgr := range l.eventmanagers {
		if emgr.Type() == eventmgrtype {
			return emgr
		}
	}
	return nil
}
