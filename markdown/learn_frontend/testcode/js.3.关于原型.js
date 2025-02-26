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

function test_how_keyword_new_works() {
    // new 操作符是一种特殊的调用方式，用于创建对象的实例
    // 它与类（class）概念相似，但 JavaScript 在 ES6 引入 class 关键字之前，并没有真正的类。
    // new 操作符提供了一种模拟面向对象编程（OOP）的方式
    // *** ES中提供的class在原理上仍是使用的原型 ***

    // 具体说明，new 操作符会做以下几件事情：
    // 0. new 接收一个构造函数和一组调用参数
    // 1. 创建一个新对象
    // 2. 将构造函数的 prototype 属性设置到新对象
    // 3. 将参数传入构造函数，执行构造函数
    // 4. 返回对象。（如果构造函数返回一个对象，则直接返回此对象；若构造函数返回空或undefined，则返回第一步中创建的对象）

    // 不使用class，可通过函数模拟new 操作符工作过程，下面是创建一个Person对象实例的例子
    function Person(name, age) {
        this.name = name;
        this.age = age;
        this.say = function () {
            console.log(`My name is ${this.name}, I'm ${this.age} years old.`);
        };
    }

    var p = new Person('John', 30);
    p.say();

    // 若要实现继承（共享属性/方法），可直接修改 prototype 属性
    function Bird(name) {
        this.name = name;
    }

    Bird.prototype.age = 5;
    Bird.prototype.fly = function () {
        console.log(`${this.name} is flying with ${this.age}.`);
    }
    var b = new Bird('Eagle');
    b.fly();
    var c = new Bird('Parrot');
    c.fly()  // 所有的Bird实例都会共享这个Bird构造器的prototype属性，包括age属性和fly方法
}

function test_oop_with_class() {
    // 从ES6开始，JavaScript提供了class语法，用于定义类。从此成为官方推荐的 OOP 编程方式
    // 类的写法实际上也是由原型运行时来承载的，逻辑上 JavaScript 认为每个类是有共同原型的一组对象
    // *类中定义的方法和属性则会被写在原型对象之上
    class Person {
        // 类定义中的函数和字段属性，都仅供实例可访问
        constructor(name, age) {
            // 数据型属性建议定义在构造器内
            this.name = name;
            this.age = age;
            this.#private_field = 2;
        }

        instance_attr_number = 1;

        // 静态属性/方法只能被类访问，实例不可访问！
        static static_string = "string";

        // 静态方法
        static static_fn() {
            console.log("static_fn", this.name); // 静态方法不可以访问实例属性
        }

        say() {
            console.log(`My name is ${this.name}, I'm ${this.age} years old.`);
        }

        #private_field = 1;
    }

    // 在类定义外面创建的方法属于静态方法，不能通过this 访问实例属性
    Person.jump = function () {
        console.log("jump");
    };
    Person.static_attr_number = 1;
    // 类只能访问类属性、类方法；不能访问实例属性、实例方法
    console.assert(Person.static_string === "string")
    console.assert(Person.static_attr_number === 1)
    console.assert(Person.jump !== undefined)
    console.assert(Person.instance_attr_number === undefined)
    console.assert(Person.say === undefined)

    // 同样的，实例可访问实例属性，无法访问类属性/方法
    inst = new Person("john", 1);
    console.assert(inst.static_string === undefined)
    console.assert(inst.name === "john")
    console.assert(inst.say !== undefined)
    console.assert(inst.jump === undefined)
    console.assert(inst.static_attr_number === undefined)

    // // 注意：static/类属性、方法可以动态删除！
    delete Person.static_attr_number && delete Person.static_string
    console.assert(Person.static_attr_number === undefined)
    console.assert(Person.static_string === undefined)
    console.assert(inst.static_string === undefined)

    // 私有字段不可在类外访问
    console.assert(Person.#private_field === undefined, "私有字段不可在类外访问1")
    console.assert(inst.#private_field === undefined, "私有字段不可在类外访问2")
}

function test_inheritance_with_cls() {
    class Animal {
        constructor(name) {
            this.name = name;
        }

        eat() {
            console.log(`${this.name} is eating.`);
        }
    }

    // 继承自Animal
    class Bird extends Animal {
        constructor(name, age) {
            super(name); // 使用super()调用父类构造器
            this.age = age;
        }

        fly() {// 可以使用父类的属性
            console.log(`${this.name} is flying with ${this.age}.`);
        }
    }

    console.assert(new Bird("swan").name === "swan")
    new Bird("swan", 2).fly()
}

test_impl_object_oriented_with_prototype()
test_print_builtin_type_name()
test_how_keyword_new_works()
test_oop_with_class()
test_inheritance_with_cls()