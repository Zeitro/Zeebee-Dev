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

	const commonJudge = async ({ ok, isExternalLink, fetchRes, xhr }) => {
		if (!ok) {
			return {
				message: xhr?.statusText || 'System server failure!',
			};
		}
		let data;
		let contentType;
		if (fetchRes) {
			contentType = fetchRes?.headers?.get?.('Content-Type');
			if (contentType?.indexOf('application/json') === -1) {
				return {
					res: fetchRes,
				};
			}
			data = await fetchRes?.json?.();
		}
		if (xhr?.responseText) {
			data = JSON.parse(xhr?.responseText || '{}');
			contentType = data?.headers?.get?.('Content-type');
		}

		if (!data && isExternalLink) {
			return true;
		}
		if (contentType?.indexOf('application/json') === -1) {
			return {
				res: fetchRes || data,
			};
		}

		if (data?.code && data?.code === 401) {
			window.location.href = '/login';
			showMessage('Please log in again!', 'warning');
			return false;
		}
		if (data?.code && data?.code !== 200) {
			console.log('err400', data);

			return {
				message: data.message || `fail:${data.code}`,
				code: data.code,
			};
		}
		return {
			res: data.data,
		};
	};

	async function request(url, options = {}, key = url) {
		// 创建控制器用于取消请求
		const controller = new AbortController();
		controllers.set(key, controller);
		const {
			body,
			method = 'GET',
			isExternalLink = false,
			headers,
			isHideTips,
			...restConfig
		} = options;
		const token = localStorageUtil.get('userStore')?.state?.access?.access_token;
		const headersConfig = {
			'Content-Type': 'application/json',
			Authorization: '',
			...(headers || {}),
		};
		if (token && !isExternalLink) {
			headersConfig.Authorization = `Bearer ${token}`;
		}
		const fetchOptions = {
			...restConfig,
			method: method,
			signal: controller.signal,
			headers: headersConfig,
		};
		if (body instanceof Object) {
			fetchOptions.body = JSON.stringify(body);
		}
		let response;
		const finalURL =
			process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}${url}` : `${str}${url}`;
		try {
			response = await fetch(isExternalLink ? url : finalURL, fetchOptions);
		} catch (error) {
			console.log(error, '----error----');
			if (error.name !== 'AbortError' && !isHideTips) {
				showMessage(error.message || 'Network anomaly');
				return Promise.reject(error);
			}
			return;
		} finally {
			controllers.delete(key);
		}
		const data = await commonJudge({ ok: response?.ok, fetchRes: response, isExternalLink });

		if (!data) return Promise.reject({ message: '401' });
		if (data?.message) {
			console.log('0000', message);
			if (!isHideTips) {
				showMessage(data.message);
			}
			return Promise.reject(data);
		}
		return data === true ? data : data.res;
	}

	// 上传需要得知进度时，只能用XMLHttpRequest，fetch现不支持
	function uploadWithProgress(url, options) {
		return new Promise((resolve, reject) => {
			const {
				onProgress,
				body,
				method = 'GET',
				isExternalLink = false,
				isHideTips,
			} = options;
			const xhr = new XMLHttpRequest();

			const finalURL =
				process.env.NEXT_PUBLIC_ENV !== 'local' ? `${baseURL}${url}` : `${str}${url}`;

			xhr.open(method, isExternalLink ? url : finalURL, true);

			// ✅ 3. 设置认证 Token Header
			const token = localStorageUtil.get('userStore')?.state?.access?.access_token;
			if (token && !isExternalLink) {
				xhr.setRequestHeader('Authorization', `Bearer ${token}`);
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
			xhr.onload = async () => {
				const data = await commonJudge({ ok: xhr.status === 200, xhr, isExternalLink });
				if (!data) return reject({ message: '401' });
				if (data.message) {
					if (!isHideTips) {
						showMessage(data.message, 'warning');
					}
					return reject(data);
				}
				resolve(data === true ? data : data.res);
			};

			// 监听网络等其他错误
			xhr.onerror = () => {
				showMessage('Network anomaly');
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
