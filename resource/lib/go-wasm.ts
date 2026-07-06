import Go from "./wasm_exec.js"

class GoModule {
    private go: Go = new Go()
    private module: string

    constructor (module: string) {
        if ((window as any).go_exports == undefined) {
            (window as any).go_exports = new Map<string, any>
        }

        this.go.argv = [module]
        this.module = module
        
        WebAssembly.instantiateStreaming(fetch("/resource/wasm/" + module + ".wasm"), this.go.importObject)
            .then((result: WebAssembly.WebAssemblyInstantiatedSource) => {
                this.go.run(result.instance);
                console.log("The " + module + " initialized");
            });
    }

    async call(name: string, ...args: any): Promise<any> {
        while (!(window as any).go_exports.has(this.module)) {
            await new Promise((resolve) => setTimeout(resolve, 10));
        }

        return (window as any).go_exports.get(this.module).get(name)(...args)
    }
}

async function load(module: string): Promise<GoModule> {
    return new GoModule(module)
}

export {
    load,
    GoModule
}