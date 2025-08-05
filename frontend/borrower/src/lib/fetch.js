import { getLocalStorage } from '@/lib/local';
import { message } from 'antd';
// 公共提示
function showMessage(msg) {
	message.destroy();
	message.error(msg);
	// todo
}
const str = '/api/v1';
const baseURL = process.env.NEXT_PUBLIC_API_URL + str;

export function createRequest() {
	const controllers = new Map();

	async function request(url, options = {}, key = url) {
		// 创建控制器用于取消请求
		const controller = new AbortController();
		controllers.set(key, controller);

		const token = getLocalStorage('access')?.access_token;
		const headers = {
			'Content-Type': 'application/json',
			Authorization: '',
			...(options.headers || {}),
		};
		if (token) {
			headers.Authorization = `Bearer ${token}`;
		}
		const fetchOptions = {
			...options,
			method: options.method,
			signal: controller.signal,
			headers,
		};
		if (options.body instanceof Object) {
			fetchOptions.body = JSON.stringify(options.body);
		}

		// if (options.method === 'POST') {
		// 	fetchOptions.body = JSON.stringify(options.body);
		// }

		try {
			const response = await fetch(
				process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}${url}` : `${str}${url}`,
				fetchOptions,
			);

			// 响应拦截：统一处理响应
			if (!response.ok) {
				showMessage(`fail:${response.status}`);
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();

			if (data.code && data.code !== 200) {
				showMessage(`fail:${data.code},message:${data.message}` || `fail:${data.code}`);
				throw new Error(data.message || `fail:${data.code}`);
			}

			return data.data;
		} catch (error) {
			if (error.name !== 'AbortError') {
				showMessage(error.message || 'Network anomaly');
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
	};

	return { request, cancel, cancelAll };
}
