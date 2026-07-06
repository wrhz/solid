import Go from "./wasm_exec.js"

class GoModule {
    private go: Go = new Go()

    constructor (module: string) {
        this.go.argv = [module]
        
        WebAssembly.instantiateStreaming(fetch("/resource/wasm/" + module + ".wasm"), this.go.importObject)
            .then((result: WebAssembly.WebAssemblyInstantiatedSource) => {
                this.go.run(result.instance)
            });
    }
}

async function load(module: string): Promise<GoModule> {
    return new GoModule(module)
}

export {
    load
}