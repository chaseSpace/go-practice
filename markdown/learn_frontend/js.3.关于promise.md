## 理解Promise

### 理解事件循环

1. 执行栈：首先执行的是同步代码，即执行栈（execution context stack）中的代码。
2. 微任务队列：当执行栈清空后，事件循环会检查微任务队列，如果有微任务，就会立即执行所有微任务，直到队列为空。
3. 宏任务队列：微任务执行完毕后，事件循环会检查宏任务队列，执行一个宏任务（如果存在）。注意，每次事件循环迭代中只会执行一个宏任务。
4. 渲染：在宏任务执行后，浏览器可能会进行渲染（如果有必要）。
5. 重复迭代：事件循环会不断重复上述过程，直到执行栈、微任务队列和宏任务队列都为空。

### 微任务和宏任务

在 JavaScript 的事件循环（Event Loop）中，任务被分为两种类型：微任务（Microtasks）和宏任务（Macrotasks）。
我们把**宿主发起的任务**称为**宏观任务**，把 JavaScript 引擎发起的任务称为微观任务。

#### 微任务（Microtasks）

微任务是最高优先级的任务，它们在当前执行栈清空后立即执行。微任务通常由 JavaScript 引擎内部创建，用于处理那些需要尽快完成的操作，以避免延迟。微任务的典型来源包括：

* Promise 的 then、catch 和 finally 回调。
* Promise 的 resolve 和 reject 操作。
* Object.observe（如果存在的话，用于对象观察）。
* MutationObserver 的回调，用于监听 DOM 变动。
* IntersectionObserver 的回调，用于监听元素与视口的交叉状态。

#### 宏任务（Macrotasks）

宏任务是由宿主环境（如浏览器或 Node.js）发起的任务，常见的宏任务包括：

- setTimeout 和 setInterval（浏览器环境）
- setImmediate（Node.js 环境）
- requestAnimationFrame（浏览器环境）
- I/O 操作（如文件读写、网络请求，Node.js 环境）
- 用户交互事件（如点击、键盘输入等，浏览器环境）