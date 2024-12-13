const assert = require('assert');

function test_js_undefined_null() {
    // 两种特殊类型：undefined 和 null
    // Undefined类型，表示未定义，只有一个值：undefined
    // Null类型 只有一个值：null

    // 什么时候会得到 undefined 值？
    // 1. 任何变量在赋值前是 Undefined 类型，值为 undefined
    // 2. 手动声明: var a = undefined;
    var a1 = undefined;
    var a1_1;
    assert.ok(typeof a1 === 'undefined' && a1 === undefined && a1 === a1_1)
    var a2 = void 0;
    assert.ok(a1 === a2)

    // Null类型，表示的是：“定义了但是为空”，一般叫空值，与 undefined 不同，null 是 JavaScript 关键字
    // 什么时候会得到 null 值？
    // 1. 函数的返回值是 null
    // 1. 手动声明: var a = null;
    var n1 = null;
    assert.ok(n1 === null)
    var n2 = function () {
    };
    assert.equal(n2(), null)
}

function test_js_number() {
    var a = Infinity;
    var b = -Infinity;

    // 存在 +0 和 -0 的特殊值
    // - 任何数除以0得到 Infinity
    // - 任何数除以-0得到 -Infinity
    assert.equal(+0, -0, "+0 should equal to -0")
    assert.equal(1 / 0, Infinity, "1 / 0 should equal to infinity")
    assert.equal(1 / -0, -Infinity)

    // 计算后的浮点数总是不等的
    // - 浮点数运算的精度问题导致等式左右的结果并不是严格相等，而是相差了个微小的值
    // - 所以实际上，这里错误的不是结论，而是比较的方法，正确的比较方法是使用 JavaScript 提供的最小精度值：
    assert.equal(1.1, 1.1)
    assert.notEqual(0.1 + 0.2, 0.3)
    assert.ok(Math.abs(0.1 + 0.2 - 0.3) < Number.EPSILON) // 检查等式左右两边差的绝对值是否小于最小精度，才是正确的比较浮点数的方法
}

function test_builder_func() {
    // Number、String 和 Boolean，三个构造器是两用的，当跟 new 搭配时，它们产生对象；
    var s1 = new String('abc');
    assert.ok(typeof s1 === 'object' && s1 !== 'abc')

    var n1 = new Number(1);
    assert.ok(typeof n1 === 'object')

    var b1 = new Boolean(true);
    assert.ok(typeof b1 === 'object')

    // 不搭配new时，它们表示强制类型转换
    var s2 = String(1);
    assert.ok(typeof s2 === 'string' && s2 === '1')
}

function test_js_fuzziness() {
    // JavaScript 语言设计上试图模糊对象和基本类型之间的关系
    // 我们可以把对象的方法在基本类型上使用
    // 点 是一个装箱操作！
    String.prototype.hello = () => console.log("hello test_js_fuzziness");
    "".hello()
}

function test_type_conversion_string2num() {
    // 类型转换 —— Number 和 String 之间的转换

    // 1. string形式的数值面量有几种进制：十进制、二进制、八进制和十六进制（都支持转number）
    // 十进制：10
    // 二进制：0b10
    // 八进制：0o10
    // 十六进制：0xff
    // string还支持正负号科学计数法，E/e表示：1.23e-4， -1.23e-4

    // 2. 常用来str2number的函数是 parseInt 和 parseFloat、Number（推荐）
    // - parseInt 不支持科学计数法！！！
    // - parseInt(string, [radix]) 只传第一个参数时只支持 16 进制前缀“0x”，而且会忽略非数字字符
    // - 在一些古老的浏览器环境中，parseInt 还支持 0 开头的数字作为 8 进制前缀，这是很多错误的来源（无法演示）
    assert.equal(parseInt("0x10"), 16)
    assert.equal(parseInt("0x10x"), 16) // 忽略非数字字符（以及后续字符）
    assert.equal(parseInt("1e2"), 1) // 不支持科学计数法
    // - 在任何环境下，都建议传入 parseInt 的第二个参数，而 parseFloat 则直接把原字符串作为十进制来解析
    assert.equal(parseInt("0100", 10), 100)
    assert.equal(parseInt(" 10", 10), 10) // 忽略前导空格
    assert.equal(parseInt("", 10), NaN) // 空值转换为NaN

    // - 使用parseFloat(string) 解析浮点数字符串，支持科学计数法
    assert.equal(parseFloat("0b10"), 0) // 不支持2进制
    assert.equal(parseFloat("0x10"), 0) // 不支持十六进制
    assert.equal(parseFloat("0o10"), 0) // 不支持8进制
    assert.equal(parseFloat("10"), 10)
    assert.equal(parseFloat("1.23e2"), 123)
    assert.equal(parseFloat("-1.23e-2"), -0.0123)
    assert.equal(parseFloat(""), NaN)

    // - Number(string) 支持进制、科学计数法、忽略前导空格，常用！
    assert.equal(Number("0b10"), 2)
    assert.equal(Number("0x10"), 16)
    assert.equal(Number("0o10"), 8)
    assert.equal(Number(" 10"), 10)
    assert.equal(Number(""), 0) // 不同于 parseInt/parseFloat的行为，Number空值转换为0
    assert.equal(Number("1.23e2"), 123)
    assert.equal(Number("-1.23e-2"), -0.0123)
}

function test_type_conversion_num2string() {
    // toString([radix])
    assert.equal((123).toString(), "123")
    assert.equal((123).toString(2), "1111011")
    assert.equal((123).toString(8), "173")
    assert.equal((123).toString(16), "7b")
    assert.equal((1e222).toString(), "1e+222") // 遇到大数时会自动变成科学计数法
    assert.equal((1e-222).toString(), "1e-222") // 同理
    assert.equal((123).toExponential(), "1.23e+2")  // 科学计数
    assert.equal((123).toFixed(2), "123.00") // 保留小数点后2位
    // assert.equal((123).toPrecision(0), "1e+2") // 崩溃，参数必须在 [1,100] 内
    assert.equal((123).toPrecision(1), "1e+2") // 指定小数点前后的总位数
    assert.equal((123).toPrecision(2), "1.2e+2") // 指定小数点前后的总位数
    assert.equal((123).toPrecision(3), "123") // 指定小数点前后的总位数
    assert.equal((123).toPrecision(4), "123.0") // 指定小数点前后的总位数
}

function test_symbol_boxing() {
    // Symbol 函数无法使用 new 来调用，但我们仍可以利用装箱机制来得到一个 Symbol 对象
    // 方式：我们可以利用一个函数的 call 方法来强迫产生装箱

    var symbolObject = (function () {
        return this;
    }).call(Symbol("a"));

    assert.equal(typeof symbolObject, 'object')
    assert.equal(symbolObject.constructor, Symbol)
    assert.ok(symbolObject instanceof Object)

    // 获取任意实例的构造函数名
    assert.equal(Object.getPrototypeOf(symbolObject).constructor.name, 'Symbol')
    assert.equal(Object.getPrototypeOf(Number(1)).constructor.name, 'Number')
    assert.equal(Object.getPrototypeOf(Boolean(true)).constructor.name, 'Boolean')
    assert.equal(Object.getPrototypeOf(String("")).constructor.name, 'String')

    // 每一类装箱对象都有一个原型名称（也可以叫类名），不可更改
    // call方法会产生装箱行为，有性能开销！
    assert.equal(Object.prototype.toString.call(symbolObject), '[object Symbol]')
    assert.equal(Object.prototype.toString.call(Number(1)), "[object Number]")
}

function test_object_unboxing() {
    // 拆箱
    // 在 JavaScript 标准中，规定了 ToPrimitive 函数，它是对象类型到基本类型的转换（即，拆箱转换）
}

function main() {
    test_js_undefined_null()
    test_js_number()
    test_builder_func()
    test_js_fuzziness()
    test_type_conversion_string2num()
    test_type_conversion_num2string()
    test_symbol_boxing()
}

main()