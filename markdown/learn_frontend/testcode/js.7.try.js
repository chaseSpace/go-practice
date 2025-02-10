function testFinallyReturn() {
    function bar() {
        try {
            return 0;
        } catch (err) {

        } finally { // finally 的内容一定会执行
            console.log("a")
            return 1; // 会覆盖上面的0，实践中不建议 return 在 finally 中
        }
    }

    console.assert(bar() === 1)
}

testFinallyReturn()