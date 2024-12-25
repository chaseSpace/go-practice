// 使用原型实现继承（面向对象的特征之一）
function test_impl_object_oriented_with_prototype() {
    // 创建cat对象，用于行为say、jump
    var cat = {
        name: "cat",
        say() {
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
}

test_impl_object_oriented_with_prototype()