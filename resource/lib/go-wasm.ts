import Go from "./wasm_exec.js"

const go = new Go();

async function load(module: string): Promise<any> {
    if ((window as any).goExports == undefined) {
        (window as any).goExports = {};
    }

    go.argv = [module];
    
    const result = await WebAssembly.instantiateStreaming(fetch("/resource/wasm/" + module + ".wasm"), go.importObject)
    
    go.run(result.instance);
    
    console.log("The " + module + " initialized");

    return (window as any).goExports[module]
}

export {
    load
}