
(() => {
    const JS_MB = '%s';
    const JS_IPC = '%s';
    const JS_JS2GO = '%s';
    const JS_GO2JS = '%s';
    const JS_REGISTER_HANDLER = '%s';

    // MB

    window.top[JS_MB] = window.top[JS_MB] || {}
    const mb = window.top[JS_MB];
    mb.newMsg = newMsg;
    mb.replyWaiting = mb.replyWaiting || {};
    mb.handlers = mb.handlers || {};


    // IPC
    window.top[JS_IPC] = window.top[JS_IPC] || {}
    const ipc = window.top[JS_IPC];
    ipc.invoke = invoke;
    ipc.sent = sent;
    ipc.handle = handle;

    // GO 调用 (JS预留函数)
    window.top[JS_GO2JS] = (msgTxt) => {
        const msg = JSON.parse(msgTxt);
        if (msg.replyId) {
            handleReply(msg)
            return
        }

        handleChannel(msg)
        return
    };

    // JS调用 (GO预埋点)
    const toGO = (msg) => window.top[JS_JS2GO](JSON.stringify(msg))
    const registerHandlerToGo = window.top[JS_REGISTER_HANDLER]

    // 注册 callJsFunc (仅 JS 端)
    ipc.handle('callJsFunc', async function (fn, ...args) {
        const func = window.top[fn]
        if (!func) {
            return new Error(`JS function ${fn} not found!`)
        }
        return await func(...args)
    }, true)

    function randStr(len = 8) {
        let characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let randomString = '';
        for (let i = 0; i < len; i++) {
            let randomIndex = Math.floor(Math.random() * characters.length);
            randomString += characters.charAt(randomIndex);
        }
        return randomString;
    }
    function newMsg({ id = '', replyId = '', channel = '', args = [], result = null, error = '' }) {
        return { id, replyId, channel, args, result, error }
    }

    function withTimeout(promise, ms = 10000) {
        let timer;
        const timeout = new Promise((_, reject) => {
            timer = setTimeout(() => reject(new Error('等待IPC Handler返回处理结果超时。')), ms);
        });

        return Promise.race([promise, timeout]).finally(() => clearTimeout(timer))
    }

    // 返回值
    function handleReply(msg) {
        if (!msg.replyId) return;
        const p = mb.replyWaiting[msg.replyId]
        if (!p) return;
        if (msg.error) {
            p.reject(new Error(msg.error))
            return;
        }
        p.resolve(msg.result)
    }

    // 执行handler。（GO 调用此函数，用于执行对应的handler)
    async function handleChannel(msg) {
        const { id, channel, args = [] } = msg || {};
        if (!channel) return;
        const handler = mb.handlers[channel];
        if (!handler) return;
        if (!id) return; // ! 如果 ID 为空，则无须回复
        try {
            const res = await handler(...args); // 支持 promise
            toGO(newMsg({ replyId: id, channel, args, result: res })) // 返回结果
        } catch (err) {
            toGO(newMsg({ replyId: id, channel, args, error: err })) // 返回结果
        }
    }

    // invoke 调用, 有返回值
    function invoke(channel, ...args) {
        const msg = newMsg({ id: randStr(), channel, args });
        return withTimeout(new Promise((resolve, reject) => {
            mb.replyWaiting[msg.id] = { resolve, reject }
            toGO(msg)
        })).finally(() => {
            delete mb.replyWaiting[msg.id]
        })
    }

    // sent 调用，没有返回值
    function sent(channel, ...args) {
        const msg = newMsg({ channel, args });
        toGO(msg)
    }

    // 声明handler
    function handle(channel, handler, onlyInJS = false) {
        window.top[JS_MB] = window.top[JS_MB] || {}
        window.top[JS_MB]['handlers'] = window.top[JS_MB]['handlers'] || {};
        window.top[JS_MB]['handlers'][channel] = handler; // 存入cache
        if (!onlyInJS) registerHandlerToGo(channel)
    }

})();
