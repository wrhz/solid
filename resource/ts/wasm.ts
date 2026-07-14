import { load } from "@/go-wasm.js"

(async function () {
    const add = await load("add");

    console.log(add.add(1, 2));

    console.log(add.message);

    const struct = new add.Struct();
    
    struct.sayInfo();

    add.Struct.sayHello();

    const child = new struct.Child();
})()