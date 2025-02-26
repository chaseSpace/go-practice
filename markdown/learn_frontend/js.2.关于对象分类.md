## 原生对象

我们把 JavaScript 中，能够通过语言本身的构造器创建的对象称作原生对象。在 JavaScript 标准中，提供了 30 多个构造器。

| 基本类型    | 基础功能和数据结构 | 错误类型           | 二进制操作             | 带类型的数组            |
|---------|-----------|----------------|-------------------|-------------------|
| Boolean | Array     | Error          | ArrayBuffer       | Float32Array      |
| String  | Date      | EvalError      | SharedArrayBuffer | Float64Array      |
| Number  | RegExp    | RangeError     | DataView          | Int8Array         |
| Symbol  | Promise   | ReferenceError |                   | Int16Array        |
| Object  | Proxy     | SyntaxError    |                   | Int32Array        |
|         | Map       | TypeError      |                   | UInt8Array        |
|         | WeakMap   | URIError       |                   | UInt16Array       |
|         | Set       |                |                   | UInt32Array       |
|         | WeakSet   |                |                   | Uint8ClampedArray |
|         | Function  |                |                   |                   |

通过这些构造器，我们可以用 new 运算创建新的对象，所以我们把这些对象称作原生对象。几乎所有这些构造器的能力都是无法用纯
JavaScript 代码实现的，它们也无法用 class/extend 语法来继承。

## 函数对象与构造器对象

### 1. 函数对象

在JS中，具有 [[call]] 私有字段的对象都可以叫做函数。

JavaScript
用对象模拟函数的设计代替了一般编程语言中的函数，它们可以像其它语言的函数一样被调用、传参。任何宿主只要提供了“具有 [[call]]
私有字段的对象”，就可以被 JavaScript 函数调用语法支持。
