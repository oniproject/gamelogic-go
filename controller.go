package gamelogic

import (
	"log"
)

type controller struct {
	LogicBrick
	linkedsensors   SensorList
	linkedactuators ActuatorList
	statemask       uint
	justActivated   bool
	bookmark        bool
}

func (c *controller) IsJustActivated() bool { return c.justActivated }

func (c *controller) ClrJustActivated()         { c.justActivated = false }
func (c *controller) SetBookmark(bookmark bool) { c.bookmark = bookmark }

// TODO
func (c *controller) Activate(ControllerList) {}
func (c *controller) Deactivate()             {}

/* FIXME
func (c *controller) Deactivate() {
	// the controller can only be part of a sensor m_newControllers list
	c.Delink()
}
func (c *controller) Activate( SG_DList& head ) {
	 TODO
	if (QEmpty())
	{
		if (m_bookmark)
		{
			m_gameobj->m_activeBookmarkedControllers.QAddBack(this);
			head.AddFront(&m_gameobj->m_activeBookmarkedControllers);
		}
		else
		{
			InsertActiveQList(m_gameobj->m_activeControllers);
			head.AddBack(&m_gameobj->m_activeControllers);
		}

}
*/

func (c *controller) SetState(state uint) { c.statemask = state }

/**
 * Use of SG_DList element: none
 * Use of SG_QList element: build ordered list of activated controller on the owner object
 *                          Head: SCA_IObject::m_activeControllers
 */

type Controller interface {
	LogicBrick

	//ProcessReplica() // FIXME

	//virtual void Trigger(class SCA_LogicManager* logicmgr)=0;
	LinkToSensor(Sensor)
	LinkToActuator(Actuator)
	LinkedSensors() SensorList
	LinkedActuators() ActuatorList
	UnlinkAllSensors()
	UnlinkAllActuators()
	UnlinkActuator(actua Actuator)
	UnlinkSensor(sensor Sensor)
	SetState(uint)
	ApplyState(uint)
	IsJustActivated() bool
	ClrJustActivated()
	SetBookmark(bookmark bool)

	Deactivate()
	Activate(ControllerList)
}

func NewController(gameobj Object) Controller {
	return &controller{
		LogicBrick:    NewLogicBrick(gameobj),
		statemask:     0,
		justActivated: false,
	}
}

/*
SCA_IController::~SCA_IController()
{
	//UnlinkAllActuators();
}*/

func (c *controller) LinkedSensors() SensorList {
	return c.linkedsensors
}
func (c *controller) LinkedActuators() ActuatorList {
	return c.linkedactuators
}

func (c *controller) UnlinkAllSensors() {
	for _, sensor := range c.linkedsensors {
		if c.IsActive() {
			sensor.DecLink()
		}
		sensor.UnlinkController(c)
	}
	c.linkedsensors = nil
}
func (c *controller) UnlinkAllActuators() {
	for _, actuator := range c.linkedactuators {
		if c.IsActive() {
			actuator.DecLink()
		}
		actuator.UnlinkController(c)
	}
	c.linkedactuators = nil
}

func (c *controller) LinkToActuator(actua Actuator) {
	c.linkedactuators = append(c.linkedactuators, actua)
	if c.IsActive() {
		actua.IncLink()
	}
}
func (c *controller) UnlinkActuator(actua Actuator) {
	for i, actit := range c.linkedactuators {
		if actit == actua {
			if c.IsActive() {
				actit.DecLink()
			}
			c.linkedactuators = append(c.linkedactuators[:i], c.linkedactuators[i+1:]...)
			return
		}
	}

	log.Printf("Missing link from controller %s:%s to actuator %s:%s\n",
		c.Parent().Name(), c.Name(), actua.Parent().Name(), actua.Name())
}

func (c *controller) LinkToSensor(sensor Sensor) {
	c.linkedsensors = append(c.linkedsensors, sensor)
	if c.IsActive() {
		sensor.IncLink()
	}
}
func (c *controller) UnlinkSensor(sensor Sensor) {
	for i, sens := range c.linkedsensors {
		if sens == sensor {
			if c.IsActive() {
				sensor.DecLink()
			}
			c.linkedsensors = append(c.linkedsensors[:i], c.linkedsensors[i+1:]...)
			return
		}
	}

	log.Printf("Missing link from controller %s:%s to sensor %s:%s\n",
		c.Parent().Name(), c.Name(), sensor.Parent().Name(), sensor.Name())
}

func (c *controller) ApplyState(state uint) {
	if c.statemask&state != 0 {
		if !c.IsActive() {
			// reactive the controller, all the links to actuator are valid again
			for _, actit := range c.linkedactuators {
				actit.IncLink()
			}
			for _, sensor := range c.linkedsensors {
				sensor.IncLink()
			}
			c.SetActive(true)
			c.justActivated = true
		}
	} else if c.IsActive() {
		for _, actit := range c.linkedactuators {
			actit.DecLink()
		}
		for _, sensor := range c.linkedsensors {
			sensor.DecLink()
		}
		c.SetActive(false)
		c.justActivated = false
	}
}
