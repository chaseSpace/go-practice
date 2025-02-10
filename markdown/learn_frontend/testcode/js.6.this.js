// this 对于普通函数来说就是函数【执行时而不是定义时】的作用域
function testShowThis() {
    function showThis() {
        console.log(this);
    }

    var o = {
        showThis: showThis
    }
    showThis(); // global环境
    o.showThis(); // o本身
}

// this 对于箭头函数来说固定获取global作用域
function testShowArrowThis() {
    const showThis = () => {
        console.log(this);
    }

    var o = {
        showThis: showThis
    }
    showThis(); // global环境
    o.showThis(); // global环境
}

// 类方法的情况不一样：undefined
function testClsShowThis() {
    class C {
        showThis() {
            console.log(this);
        }
    }
    var o = new C();
    var showThis = o.showThis;

    showThis(); // undefined
    o.showThis(); // o
}

function testSetThis(){
    // Function.prototype.call 和 Function.prototype.apply 可以指定函数调用时传入的 this 值
    function foo(a, b, c){
        console.log(this);
        console.log(a, b, c);
    }
    // 这里 call 和 apply 作用是一样的，只是传参方式有区别
    foo.call({}, 1, 2, 3);
    foo.apply({}, [1, 2, 3]);
    foo.bind({},1,2,3)(); // bind 返回一个函数，需要调用

    // 但是，箭头函数和cls方法不接受this指派，语法仍能使用，但不生效
    const arrowFoo = (a, b, c) => {
        console.log(this);
        console.log(a, b, c);
    }
    arrowFoo.call({}, 1, 2, 3);
}


/*
* 说结论：生成器函数、异步生成器函数和异步普通函数跟普通函数行为是一致的，异步箭头函数与箭头函数行为是一致的
* */

// testShowThis()
// testShowArrowThis()
// testClsShowThis()
testSetThis()