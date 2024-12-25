// 对象的数据属性
function test_obj_data_property() {
    var o = {a: 1};
    o.b = 2;
    //a 和 b 皆为数据属性
    console.log(Object.getOwnPropertyDescriptor(o, "a")) // {value: 1, writable: true, enumerable: true, configurable:true}
    console.log(Object.getOwnPropertyDescriptor(o, "b")) // {value: 2, writable: true, enumerable: true, configurable: true}

    // defineProperty 修改特性
    Object.defineProperty(o, "b", {value: 2, writable: false, enumerable: false, configurable: true});
    console.log("修改b的特征后：", Object.getOwnPropertyDescriptor(o, "b"))
    o.b = 3; // b 无法被修改，但不会报错
    console.assert(o.b === 2, "b 无法被修改，不会报错");
}

// 对象的访问属性
function test_obj_access_property() {
    // 在创建对象时，也可以使用 get 和 set 关键字来创建访问器属性
    var o = {
        _a: 1,
        get b() {
            return this._a + 1
        },
        set b(v) {
            this._a = v * 2
        }
    }
    /* 访问器属性的特征:
    {
      get: [Function: get b],
      set: [Function: set b],
      enumerable: true,
      configurable: true
    }
    */
    console.log("访问器属性的特征：", Object.getOwnPropertyDescriptor(o, "b"))
    console.assert(o.b === 2, "访问器属性 b 的值是 _a 的值加 1")
    o.b = 3;
    console.assert(o.b === 7, "访问器属性 b 的值是新值 *2 再加 1")
}

test_obj_data_property()
test_obj_access_property()
