# core

基于Golang实现的EC组件框架，可以用作各种分布式系统代码框架，例如游戏服务器、APP服务器、爬虫服务器等等。

## 主要功能
### EC组件系统
* EC组件系统是OOP的一种实现，主张用组合代替继承，一般由 Entity <-> Component 组成。Entity代表实体对象，用于管理Component，Component代表方法和数据的集合，用于实现逻辑功能。
* 主要优点是在程序运行期间，能动态添加和删除组件，修改对象的能力，例如游戏中 兔子 -> 狼，狼 -> 斑马，通过挂载不同Component就可以很方便的实现切换。
* 因为EC树只是一种Entity的管理方式，并非EC系统必须包含的，所以本层不提供EC树。

### 运行时
* 结合EC组件系统，设计了一种运行时架构，让不同的Entity工作在指定的线程中，并且工作在不同线程上的Entity可以跨线程相互通信。
* 线程之间通信使用 闭包 + 队列 的方式，即将闭包压入被调线程的调用队列，这样可以在运行时中保持调用顺序，实现线程安全。

### 事件系统
* 类似C#的事件机制，注意非线程安全。

## 主要对象与函数
### Context
* 上下文，提供存储变量、跨线程报告异常、停止线程几项功能，贯穿所有代码。

### APP
* 应用，从Context继承，同时增加了Entity管理、开始停止几项功能，贯穿所有代码。

### Entity
* 实体，提供Component管理功能。

### Component
* 组件，提供一组生命周期回调函数：
` [Init] -> [Awake] -> [EntityInit] -> [Start] -> [Update] -> [LateUpdate] -> [EntityShut] -> [Halt] -> [Shut]`

### Runtime
* 线程运行时，从Context继承，用于给Entity提供多线程运行环境，提供GC能力，贯穿所有代码。

### Frame
* 结合Runtime，可以调整Runtime的运行方式。

### SafeStack
* 跨线程安全调用栈，可用于提供不会死锁的跨线程调用，会阻塞当前线程。

### UnsafeCall
* 跨线程不安全调用，可以选择是否阻塞当前线程，阻塞当前线程可能会造成死锁，例如：`线程A -> 线程B -> 线程A`。

### 事件：
* Hook 
	* 钩子，事件的接收端，可以同时绑定多个EventSource。

* EventSource 
	* 事件源，事件的发送端，可以被不同的Hook绑定。

### 事件函数：
* BindEvent 
	* 绑定Hook与EventSource。

* UnbindEvent 
	* 解绑定Hook与EventSource。

* UnbindAllEventSource 
	* Hook解绑定所有已绑定的EventSource。

* UnbindAllHook 
	* EventSource解绑定所有已绑定的Hook。
