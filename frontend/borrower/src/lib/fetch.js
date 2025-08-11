import { localStorageUtil } from '@/lib/local';
import { message } from 'antd';
// 公共提示
function showMessage(msg, type = 'error') {
	message.destroy();
	message.open({
		content: msg,
		type,
	});
}
const str = '/api/v1';
const baseURL = process.env.NEXT_PUBLIC_API_URL + str;

export function createRequest() {
	const controllers = new Map();

	async function request(url, options = {}, key = url) {
		// 创建控制器用于取消请求
		const controller = new AbortController();
		controllers.set(key, controller);

		const token = localStorageUtil.get('userStore')?.state?.access?.access_token;
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
				showMessage(`fail:${response.status},${response.statusText}`);
				console.log(response);

				return Promise.reject(`HTTP error! status: ${response.status}`);
			}

			const data = await response.json();

			if (data.code === 401) {
				showMessage('Please log in again!', 'warning');
				window.location.href = '/login';
				return Promise.reject('Please log in again!');
			}

			if (data.code && data.code !== 200) {
				showMessage(`fail:${data.code},message:${data.message}`);
				return Promise.reject(data.message || `fail:${data.code}`);
			}

			return data.data;
		} catch (error) {
			if (error.name !== 'AbortError') {
				showMessage(error.message || 'Network anomaly');
				return Promise.reject(error);
			}
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
