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
* APP
* Entity
* Component
* Runtime
* Frame
* Hook
* EventSource
* BindEvent
* UnbindEvent
* UnbindAllEventSource
* UnbindAllHook
* SafeStack
* UnsafeCall
