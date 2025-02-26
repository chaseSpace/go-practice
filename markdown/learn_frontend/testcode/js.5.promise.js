function test_promise() {
    var r = new Promise(function (resolve, reject) {
        console.log("a"); // 同步执行
        resolve() // 异步
    });
    r.then(() => console.log("c")); // promise.then是(resolve, reject)的回调函数，所以是最后执行
    console.log("b") // 先于 c 输出
}

// setTimeout、setInterval是浏览器提供的宏任务函数，优先级低于Promise产生的微任务
// 即使Promise.then延时1s添加异步任务，也是优先于 setTimeout
function test_setTimeout_priority_lower_promise() {
    setTimeout(() => {
        console.log("setTimeout")
    }, 0)
    new Promise(function (resolve, reject) {
        resolve()
    }).then(() => {
        var begin = Date.now();
        while (Date.now() - begin < 1000) ; // sleep 1s
        console.log("Promise")
    })
}

function test_wrap_sleep() {
    // 实践中将 setTimeout 包装成支持异步的 sleep 函数
    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    sleep(1000).then(() => {
        console.log("then ...")
    })
}


// async/await 是 ES2016 新加入的特性，它提供了用 for、if 等代码结构来编写异步的方式。
// 它的运行时基础是 Promise
async function test_async_await() {
    // 返回 Promise 对象的函数也是一个异步函数，等效于 async 函数
    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    async function foo() {
        console.log("foo-a")
        await sleep(2000)
        console.log("foo-b")
    }

    await foo()
}

test_promise()
test_setTimeout_priority_lower_promise()
test_wrap_sleep()
test_async_await()