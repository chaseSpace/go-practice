const assert = require('assert');

function test_host_object() {
    // nodejs 中不提供 window 对象
    assert.equal(typeof window, 'undefined')
    console.assert(typeof globalThis, 'object')
}

function test_func_and_constructor() {
    // 1. 函数对象的定义是：具有 [[call]] 私有字段的对象
    // [[call]] 私有字段必须是一个引擎中定义的函数，我们无法在用户态中定义它。

    // 2. 构造器对象的定义是：具有 [[construct]] 私有字段的对象

    // 大部分内置对象 都实现了 [[call]] 和 [[construct]] 私有字段。只是它们的行为不同
    a = new String;
    console.log(a) // [String: '']  构造器语法是创建一个对象
    a = String(123)
    console.assert(a === '123') // 函数语法是做类型转换

    d = new Date;
    console.log(d) // 2024-12-30T07:15:55.947Z Date构造器返回一个Date对象
    console.assert(typeof d === 'object')
    d = Date()
    console.assert(typeof d === 'string')

    // 对 function 语法，两种语法的行为基本一致，但构造器仍然返回一个对象
    function f() {
        console.log('111')
        return 'f'
    }

    // 都会打印 111
    console.assert(f() === 'f')
    console.log(typeof new f === 'object')
}

test_host_object()
test_func_and_constructor()