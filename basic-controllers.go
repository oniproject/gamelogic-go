package gamelogic

// TODO ExpressionController

// возможно c.Parent() это неправильно

type cXOR struct {
	Controller
}

func XOR(gameobj Object) Controller {
	return &cXOR{
	//TODO
	}
}

/*
func (c *cXOR) Replica() interface{} {
	replica := XOR(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}
*/
func (c *cXOR) Trigger(logicmgr *LogicManager) {
	result := false
	for _, sensor := range c.LinkedSensors() {
		if sensor.State() {
			if result == true {
				result = false
				break
			}
			result = true
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}

type cXNOR struct {
	Controller
}

func XNOR(gameobj Object) Controller {
	return &cXNOR{
	//TODO
	}
}

/*
func (c *cXNOR) Replica() interface{} {
	replica := XNOR(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}
*/
func (c *cXNOR) Trigger(logicmgr *LogicManager) {
	result := true
	for _, sensor := range c.LinkedSensors() {
		if sensor.State() {
			if result == false {
				result = true
				break
			}
			result = false
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}

type cAND struct {
	Controller
}

func AND(gameobj Object) Controller {
	return &cAND{
	//TODO
	}
}

/*
func (c *cAND) Replica() interface{} {
	replica := AND(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}*/
func (c *cAND) Trigger(logicmgr *LogicManager) {
	result := true
	for _, sensor := range c.LinkedSensors() {
		if !sensor.State() {
			result = false
			break
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}

type cNAND struct {
	Controller
}

func NAND(gameobj Object) Controller {
	return &cNAND{
	//TODO
	}
}

/*
func (c *cNAND) Replica() interface{} {
	replica := NAND(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}*/
func (c *cNAND) Trigger(logicmgr *LogicManager) {
	result := false
	for _, sensor := range c.LinkedSensors() {
		if !sensor.State() {
			result = true
			break
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}

type cOR struct {
	Controller
}

func OR(gameobj Object) Controller {
	return &cOR{
	//TODO
	}
}

/*
func (c *cOR) Replica() interface{} {
	replica := OR(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}*/
func (c *cOR) Trigger(logicmgr *LogicManager) {
	result := false
	for _, sensor := range c.LinkedSensors() {
		if sensor.State() {
			result = true
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}

type cNOR struct {
	Controller
}

func NOR(gameobj Object) Controller {
	return &cNOR{
	//TODO
	}
}

/*
func (c *cNOR) Replica() interface{} {
	replica := OR(c.Parent())
	// this will copy properties and so on...
	replica.ProcessReplica()
	return replica
}*/
func (c *cNOR) Trigger(logicmgr *LogicManager) {
	result := true
	for _, sensor := range c.LinkedSensors() {
		if sensor.State() {
			result = false
		}
	}
	for _, actuator := range c.LinkedActuators() {
		logicmgr.AddActiveActuator(actuator, result)
	}
}
