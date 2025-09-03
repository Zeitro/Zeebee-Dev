// 公共提示
function showMessage() {
    // todo
}

export function createRequest() {
    const controllers = new Map();

    async function request(url, options = {}, key = url) {
        // 创建控制器用于取消请求
        const controller = new AbortController();
        controllers.set(key, controller);

        const fetchOptions = {
            ...options,
            signal: controller.signal,
            headers: {
                'Content-Type': 'application/json',
                ...(options.headers || {}),
            },
        };

        try {
            const response = await fetch(url, fetchOptions);

            // 响应拦截：统一处理响应
            if (!response.ok) {
                showMessage(`请求失败：${response.status}`);
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();

            
            if (data.code && data.code !== 0) {
                showMessage(data.message || '请求出错');
                throw new Error(data.message || '请求出错');
            }

            return data;
        } catch (error) {
            if (error.name !== 'AbortError') {
                showMessage(error.message || '网络异常');
            }
            throw error;
        } finally {
            controllers.delete(key);
        }
    }

    // 取消单个请求
    function cancel(key) {
        const controller = controllers.get(key);
        if (controller) {
            controller.abort();
            controllers.delete(key);
        }
    }

    // 取消全部请求
    const cancelAll = () => {
        controllers.forEach((controller) => {
            controller.abort();
        });
        controllers.clear();
    }

    return { request, cancel, cancelAll };
}
