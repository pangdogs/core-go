# core

基于Golang实现的EC组件框架，可以用作各种分布式系统代码组织框架，例如游戏服务器、APP服务器、爬虫服务器等等。

## 主要功能
### EC组件系统
* EC组件系统是OOP的一种实现，主张用组合代替继承，一般由 Entity <-> Component 组成。Entity代表实体对象，用于管理Component，Component代表方法和数据的集合，用于实现逻辑功能。
* 主要优点是在程序运行期间，能动态添加和删除组件，修改对象的能力，例如游戏中 兔子 -> 狼，狼 -> 斑马，通过挂载不同Component就可以很方便的实现切换。
* 因为EC树只是一种Entity的管理方式，并非EC系统必须包含的，所以本层不提供EC树。

### 运行时
* 结合EC组件系统，设计了一种运行时架构，使Entity工作在指定的线程中，工作在不同线程上的Entity可以安全的相互通信。
* 运行时之间安全通信使用 闭包 + 队列 的方式实现，即将闭包压入被调运行时的调用队列，被调运行时不断从队列中取出闭包执行，这样实现安全的线程间通信。

### 事件系统
* 类似C#的事件机制，注意非线程安全，不能用于实现跨线程通知。

## 主要对象与函数
### ServiceContext
* 服务上下文，提供存储变量、全局Entity管理、全局异常报告、全局Cancel等几项功能，所有方法线程安全。

### Service
* 服务，提供开始、停止运行功能。

### RuntimeContext
* 运行时上下文，提供存储变量、Entity管理、异常报告、Cancel、跨运行时安全调用、获取帧数据等几项功能，除了Entity管理与获取帧数据外，其他方法均线程安全。

### Runtime
* 运行时，提供开始、停止运行功能，用于驱动Entity与Component生命周期运转。

### Frame
* 帧，结合Runtime，可以调整Runtime的运行方式。

### SafeCall
* 跨运行时安全调用，递归调用会失败并超时，例如：`线程A -> 线程B -> 线程A`。

### Entity
* 实体，提供Component管理功能。
* 生命周期：`[Init] -> [InitFin] -> [Update] -> [LateUpdate] -> [Shut] -> [ShutFin] `

### Component
* 组件，用于拓展编写逻辑。
* 生命周期：`[Awake] -> [Start] -> [Update] -> [LateUpdate] -> [Shut]`

### 事件：
* Event
	* 事件，提供主动通知，即生产者。

* Hook 
	* 钩子，消费者订阅事件产生的句柄，用于解除订阅。

* BindEvent，BindEventWithPriority
	* 绑定事件回调函数。

* UnbindEvent
	* 解绑定事件回调函数，比使用Hook解除订阅性能差，且在同个回调函数绑定多次事件的情况下，只能从最后依次解除，无法指定解除哪一个。

* EventRecursion 
	* 定义递归发送事件时的行为，即在事件回调中再次发送事件如何处理。
```
  	EventRecursion_Allow    允许递归
	EventRecursion_Disallow 不允许递归
	EventRecursion_Discard  事件广播时，跳过已进入的回调函数
	EventRecursion_Deep     只会在递归最深层广播事件，并且跳过已进入的回调
```

* eventcode包
	* 用于生成发送事件代码，使用`go:generate`功能加在事件定义代码头部，即可生成代码。
	* 通常使用`//go:generate go run github.com/pangdogs/core/eventcode -decl $GOFILE -package $GOPACKAGE`
