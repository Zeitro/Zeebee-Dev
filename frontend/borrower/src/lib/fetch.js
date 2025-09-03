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
				showMessage(response.statusText);
				return Promise.reject(`HTTP error! status: ${response.status}`);
			}
			const contentType = response.headers.get('Content-Type');
			if (contentType?.indexOf('application/json') === -1) {
				return response;
			}
			const data = await response.json();

			if (data.code === 401) {
				console.log('401------login');

				showMessage('Please log in again!', 'warning');
				window.location.href = '/login';
				return Promise.reject('Please log in again!');
			}

			if (data.code && data.code !== 200) {
				showMessage(data.message);
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

	// 上传需要得知进度时，只能用XMLHttpRequest，fetch现不支持
	function uploadWithProgress(url, options) {
		return new Promise((resolve, reject) => {
			const { onProgress, body, method = 'POST', isExternalLink = false } = options;
			const xhr = new XMLHttpRequest();

			const finalURL = isExternalLink
				? url
				: process.env.NEXT_PUBLIC_ENV !== 'local'
					? `${baseURL}${url}`
					: `${str}${url}`;

			xhr.open(method, finalURL, true);

			// ✅ 3. 设置认证 Token Header
			const token = localStorageUtil.get('userStore')?.state?.access?.access_token;
			if (token && !isExternalLink) {
				xhr.setRequestHeader('Authorization', `Bearer ${token}`);
			}
			if (isExternalLink) {
				xhr.setRequestHeader('Content-Type', 'application/pdf');
				xhr.setRequestHeader('x-amz-acl', 'public-read');
				xhr.setRequestHeader('x-amz-storage-class', 'STANDARD');
			}
			// ❗️ 注意：此处不能设置 'Content-Type'，必须由浏览器自动处理 FormData

			// 监听上传进度
			xhr.upload.onprogress = (event) => {
				if (event.lengthComputable) {
					const percent = Math.round((event.loaded / event.total) * 100);
					onProgress(percent);
				}
			};

			// 监听请求完成
			xhr.onload = () => {
				if (xhr.status === 200) {
					try {
						// ✅ 4. 移植您的响应拦截逻辑
						const data = JSON.parse(xhr.responseText);
						const contentType = data.headers.get('Content-Type');

						if (contentType?.indexOf('application/json') === -1) {
							resolve(data);
						}
						if (data.code === 401) {
							showMessage('Please log in again!');
							window.location.href = '/login';
							return reject('Please log in again!');
						}

						// 以及通用业务接口返回 { code: 200, data: { url: '...' } }
						if (data.code && data.code !== 200 && data.success !== true) {
							showMessage(data.message);
							return reject(data.message || `fail:${data.code}`);
						}

						resolve(data.data);
					} catch (e) {
						reject(new Error('Failed to parse server response'));
					}
				} else {
					showMessage(xhr.statusText || `HTTP error! status: ${xhr.status}`);
					reject(`HTTP error! status: ${xhr.status}`);
				}
			};

			// 监听网络等其他错误
			xhr.onerror = () => {
				showMessage.error('Network anomaly');
				reject(new Error('Upload request failed due to a network error'));
			};

			// 发起请求
			xhr.send(body);
		});
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

	return { request, cancel, cancelAll, uploadWithProgress };
}
