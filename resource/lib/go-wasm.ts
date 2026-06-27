const wasi = {
  fd_write: () => {},
  environ_sizes_get: () => {},
};

WebAssembly.instantiateStreaming(fetch('main.wasm'), { wasi_snapshot_preview1: wasi })
  .then(result => {
    const exports = result.instance.exports;
    console.log(exports.add(3, 4));
});