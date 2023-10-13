require.config({paths: {vs: (window.basePath || '') + "static/monaco-editor/min/vs"}});
// vs/editor/editor.main.xxx.js  有多个文件需要导入
require(
    [
        "vs/editor/editor.main",
    ]
    , () => {
        window.monaco = monaco;
        let onMonacoList = window.onMonacoList || [];
        onMonacoList.forEach(one => {
            one()
        })
    })

window.onMonacoList = []
window.onMonacoLoad = (one) => {
    if (window.monaco) {
        one()
    } else {
        window.onMonacoList.push(one)
    }

}