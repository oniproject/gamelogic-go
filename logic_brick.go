package gamelogic

type logic_brick struct {
	gameobj                Object
	Execute_Priority       int
	Execute_Ueber_Priority int

	bActive  bool
	eventval Event // TODO
	text     string
	name     string
	/* protected
	//unsigned long		m_drawcolor;
	void RegisterEvent(CValue* eventval);
	void RemoveEvent();
	CValue* GetEvent();
	*/
}

type LogicBrick interface {
	ExecutePriority() int
	UeberExecutePriority() int
	SetExecutePriority(execute_Priority int)
	SetUeberExecutePriority(execute_Priority int)

	Parent() Object

	/*
		virtual void	ReParent(SCA_IObject* parent);
		virtual void	Relink(CTR_Map<CTR_HashedPtr, void*> *obj_map);
		virtual void Delete() { Release(); }
	*/

	ReParent(parent Object)
	// TODO Relink(CTR_Map<CTR_HashedPtr, void*> *obj_map);
	//virtual void Delete() { Release(); }

	// act as a BoolValue (with value IsPositiveTrigger)
	//virtual CValue*	Calc(VALUE_OPERATOR op, CValue *val);
	//virtual CValue*	CalcFinal(VALUE_DATA_TYPE dtype, VALUE_OPERATOR op, CValue *val);

	Text() string
	Number() float64
	Name() string
	SetName(string)

	IsActive() bool

	SetActive(active bool)

	// insert in a QList at position corresponding to m_Execute_Priority
	/*void			    InsertActiveQList(SG_QList& head)
	{
		SG_QList::iterator<SCA_ILogicBrick> it(head);
		for (it.begin(); !it.end() && m_Execute_Priority > (*it)->m_Execute_Priority; ++it);
		it.add_back(this);
	}*/

	// insert in a QList at position corresponding to m_Execute_Priority
	// inside a longer list that contains elements of other objects.
	// Sorting is done only between the elements of the same object.
	// head is the head of the combined list
	// current points to the first element of the object in the list, NULL if none yet
	/*
		void			    InsertSelfActiveQList(SG_QList& head, SG_QList** current)
		{
			if (!*current)
			{
				// first element can be put anywhere
				head.QAddBack(this);
				*current = this;
				return;
			}
			// note: we assume current points actually to one o our element, skip the tests
			SG_QList::iterator<SCA_ILogicBrick> it(head,*current);
			if (m_Execute_Priority <= (*it)->m_Execute_Priority)
			{
				// this element comes before the first
				*current = this;
			}
			else {
				for (++it; !it.end() && (*it)->m_gameobj == m_gameobj &&  m_Execute_Priority > (*it)->m_Execute_Priority; ++it);
			}
			it.add_back(this);
		}
	*/

	LessComparedTo(other LogicBrick) bool

	// runtime variable, set when Triggering the python controller
	//static class SCA_LogicManager*	m_sCurrentLogicManager;

	// for moving logic bricks between scenes
	//Replace_IScene(SCA_IScene *val) {}
	//Replace_NetworkScene(NG_NetworkScene *val) {}

}

func NewLogicBrick(gameobj Object) LogicBrick {
	return &logic_brick{gameobj: gameobj, text: "KX_LogicBrick"}
}

func (l *logic_brick) Parent() Object        { return l.gameobj }
func (l *logic_brick) IsActive() bool        { return l.bActive }
func (l *logic_brick) SetActive(active bool) { l.bActive = active }

//class NG_NetworkScene;
//class SCA_IScene;
/*
class SCA_ILogicBrick : public CValue
{
	Py_Header
protected:
	SCA_IObject*		m_gameobj;
	int					m_Execute_Priority;
	int					m_Execute_Ueber_Priority;

	bool				m_bActive;
	CValue*				m_eventval;
	STR_String			m_text;
	STR_String			m_name;
	//unsigned long		m_drawcolor;
	void RegisterEvent(CValue* eventval);
	void RemoveEvent();
	CValue* GetEvent();

public:
	SCA_ILogicBrick(SCA_IObject* gameobj);
	virtual ~SCA_ILogicBrick();

	void SetExecutePriority(int execute_Priority);
	void SetUeberExecutePriority(int execute_Priority);

	SCA_IObject*	GetParent() { return m_gameobj; }

	virtual void	ReParent(SCA_IObject* parent);
	virtual void	Relink(CTR_Map<CTR_HashedPtr, void*> *obj_map);
	virtual void Delete() { Release(); }

	// act as a BoolValue (with value IsPositiveTrigger)
	virtual CValue*	Calc(VALUE_OPERATOR op, CValue *val);
	virtual CValue*	CalcFinal(VALUE_DATA_TYPE dtype, VALUE_OPERATOR op, CValue *val);

	virtual const STR_String &	GetText();
	virtual double		GetNumber();
	virtual STR_String&	GetName();
	virtual void		SetName(const char *);

	bool				IsActive()
	{
		return m_bActive;
	}

	void				SetActive(bool active)
	{
		m_bActive=active;
	}

	// insert in a QList at position corresponding to m_Execute_Priority
	void			    InsertActiveQList(SG_QList& head)
	{
		SG_QList::iterator<SCA_ILogicBrick> it(head);
		for (it.begin(); !it.end() && m_Execute_Priority > (*it)->m_Execute_Priority; ++it);
		it.add_back(this);
	}

	// insert in a QList at position corresponding to m_Execute_Priority
	// inside a longer list that contains elements of other objects.
	// Sorting is done only between the elements of the same object.
	// head is the head of the combined list
	// current points to the first element of the object in the list, NULL if none yet
	void			    InsertSelfActiveQList(SG_QList& head, SG_QList** current)
	{
		if (!*current)
		{
			// first element can be put anywhere
			head.QAddBack(this);
			*current = this;
			return;
		}
		// note: we assume current points actually to one o our element, skip the tests
		SG_QList::iterator<SCA_ILogicBrick> it(head,*current);
		if (m_Execute_Priority <= (*it)->m_Execute_Priority)
		{
			// this element comes before the first
			*current = this;
		}
		else {
			for (++it; !it.end() && (*it)->m_gameobj == m_gameobj &&  m_Execute_Priority > (*it)->m_Execute_Priority; ++it);
		}
		it.add_back(this);
	}

	virtual	bool		LessComparedTo(SCA_ILogicBrick* other);

	// runtime variable, set when Triggering the python controller
	static class SCA_LogicManager*	m_sCurrentLogicManager;


	// for moving logic bricks between scenes
	virtual void		Replace_IScene(SCA_IScene *val) {}
	virtual void		Replace_NetworkScene(NG_NetworkScene *val) {}


};

#



SCA_LogicManager* SCA_ILogicBrick::m_sCurrentLogicManager = NULL;
*/

/*
SCA_ILogicBrick::~SCA_ILogicBrick()
{
	RemoveEvent();
}
*/

func (l *logic_brick) ExecutePriority() int {
	return l.Execute_Priority
}
func (l *logic_brick) UeberExecutePriority() int {
	return l.Execute_Ueber_Priority
}
func (l *logic_brick) SetExecutePriority(execute_Priority int) {
	l.Execute_Priority = execute_Priority
}
func (l *logic_brick) SetUeberExecutePriority(execute_Priority int) {
	l.Execute_Ueber_Priority = execute_Priority
}

func (l *logic_brick) ReParent(parent Object) { l.gameobj = parent }

/*
func(l *logic_brick) Relink(CTR_Map<CTR_HashedPtr, void*> *obj_map)
{
	// nothing to do
}

CValue* SCA_ILogicBrick::Calc(VALUE_OPERATOR op, CValue *val)
{
	CValue* temp = new CBoolValue(false,"");
	CValue* result = temp->Calc(op,val);
	temp->Release();

	return result;
}



CValue* SCA_ILogicBrick::CalcFinal(VALUE_DATA_TYPE dtype,
								   VALUE_OPERATOR op,
								   CValue *val)
{
	// same as bool implementation, so...
	CValue* temp = new CBoolValue(false,"");
	CValue* result = temp->CalcFinal(dtype,op,val);
	temp->Release();

	return result;
}

*/

func (l *logic_brick) Text() string {
	if l.name != "" {
		return l.name
	}
	return l.text
}
func (l *logic_brick) Number() float64     { return -1 }
func (l *logic_brick) Name() string        { return l.name }
func (l *logic_brick) SetName(name string) { l.name = name }

func (l *logic_brick) LessComparedTo(other LogicBrick) bool {
	return (l.UeberExecutePriority() < other.UeberExecutePriority()) ||
		((l.UeberExecutePriority() == other.UeberExecutePriority()) &&
			(l.ExecutePriority() < other.ExecutePriority()))
}

func (l *logic_brick) RegisterEvent(eventval Event) {
	if l.eventval != nil {
		l.eventval.Release()
	}

	l.eventval = eventval.AddRef()
}

func (l *logic_brick) RemoveEvent() {
	if l.eventval != nil {
		l.eventval.Release()
		l.eventval = nil
	}
}

func (l *logic_brick) Event() Event {
	if l.eventval != nil {
		return l.eventval.AddRef()
	}

	return nil
}
