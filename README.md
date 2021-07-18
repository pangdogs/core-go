# core

基于Golang实现的EC组件框架，可以用作各种分布式系统底层基石，例如游戏服务器、APP服务器、P2P系统、数字货币系统。

## 主要功能
* EC组件系统
	* EC组件系统是对OOP代码设计方法的扩展，EC组件系统一般由 Entity <-> Component 组成，Entity代表实体对象，可以存储公共数据，Component代表方法和数据的集合，可以抽象设计功能。不同于传统OOP，EC组件代码设计思路是利用组合代替继承，优点是可以在程序运行期间，可以通过动态修改组件来修改对象能力。

* 线程池
	* 结合EC组件系统，设计了一种线程池架构，让不同的Entity工作在指定的线程中，并且工作在不同线程上的Entity可以跨线程相互通信。

* 事件系统
	* 提供事件绑定与通知机制，结合线程池可以工作在不同线程上。

## 主要对象与函数
* Context
	* 上下文，提供存储变量、跨线程报告异常、停止线程几项功能，贯穿所有代码。

* APP
	* 应用，从Context继承，同时增加了Entity管理、开始停止几项功能，贯穿所有代码。

* Entity
	* 实体，提供Component管理功能。

* Component
	* 组件，提供一组生命周期回调函数：
	[Init] -> [Awake] -> [EntityInit] -> [Start] -> [EntityShut] -> [Halt] -> [Shut]

* Runtime
	* 线程运行时，从Context继承，用于给Entity提供多线程运行环境，贯穿所有代码。

* Frame
	* 结合Runtime，可以调整Runtime的运行方式。

* SafeStack
	* 跨线程安全调用栈，可用于提供不会死锁的跨线程调用，会阻塞当前线程。

* UnsafeCall
	* 跨线程不安全调用，可以选择是否阻塞当前线程，阻塞当前线程可能会造成死锁，例如：线程A -> 线程B -> 线程A。

* 事件：
	* Hook
		钩子，事件的接收端，可以同时绑定多个EventSource。

	* EventSource
		事件源，事件的发送端，可以被不同的Hook绑定。

* 事件函数：
	* BindEvent
	绑定Hook与EventSource。

	* UnbindEvent
	解绑定Hook与EventSource。

	* UnbindAllEventSource
	Hook解绑定所有已绑定的EventSource。

	* UnbindAllHook
	EventSource解绑定所有已绑定的Hook。
