import { load } from "@/go-wasm.js"

(async function () {
    const add = await load("add");

    console.log(await add.call("add", 1, 2));
} ());