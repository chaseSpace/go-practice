// 使用原型实现继承（面向对象的特征之一）
function test_impl_object_oriented_with_prototype() {
    // 创建cat对象（对比类），拥有行为say、jump
    var cat = {
        name: "cat", // 属性
        say() { // 行为
            console.log("meow~");
        },
        jump() {
            console.log("jump");
        }
    }

    // 通过Object.create创建cat2，实现复用cat对象，即继承
    var cat2 = Object.create(cat, {
        name: {
            value: "cat2",
            enumerable: true,
            writable: true,
            configurable: true
        }
    })
    console.assert(cat.name === "cat", "cat.name === 'cat'")
    console.assert(cat2.name === "cat2", "cat2.name === 'cat2'")

    // 不同于cpp、java，js足够“灵活”，可以在runtime中修改对象属性、添加/修改方法
    cat2.name = "cat2_new"
    cat2.jump = function () {
        console.log("jump_new");
    }
    console.assert(cat2.name === "cat2_new", "cat2.name === 'cat2_new'")
}

// 输出内置类型名称
function test_print_builtin_type_name() {
    // JS中的基础类型都是一个独特的对象（类），都有自己的唯一名字
    // 在 ES3 和之前的版本，只能使用 Object.prototype.toString 访问到内置类型名称
    var o = new Object;
    var n = new Number;
    var s = new String;
    var b = new Boolean;
    var d = new Date;
    var arg = function () {
        return arguments
    }();
    var r = new RegExp;
    var f = new Function;
    var arr = new Array;
    var e = new Error;

    /*[
      '[object Object]',
      '[object Number]',
      '[object String]',
      '[object Boolean]',
      '[object Date]',
      '[object Arguments]',
      '[object RegExp]',
      '[object Function]',
      '[object Array]',
      '[object Error]'
    ]
    */
    console.log("打印内置类型名称", [o, n, s, b, d, arg, r, f, arr, e].map(v => Object.prototype.toString.call(v)));

    // 在 ES5 开始，[[class]] 私有属性被 Symbol.toStringTag 代替
    // - 这里使用 Symbol.toStringTag 修改对象的唯一名字
    let myObj = {
        value: 123,
        [Symbol.toStringTag]: 'MyObject'
    };
    console.assert(Object.prototype.toString.call(myObj) === "[object MyObject]")
    console.assert(myObj + "" === "[object MyObject]")

    // - 这里使用 Symbol.toStringTag 修改Number的唯一名字
    let myNum = new Number; // 得到对象
    myNum[Symbol.toStringTag] = 'MyNumber';
    console.assert(Object.prototype.toString.call(myNum) === "[object MyNumber]")
    console.assert(myNum + "" === "0") // 注意！！！Number等基本类型与str相加时，进行的是底层隐式转换，不会调用 Number的toString()
}

test_impl_object_oriented_with_prototype()
test_print_builtin_type_name()