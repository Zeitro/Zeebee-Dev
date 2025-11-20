import { getFileTypeBySuffix, createRequest } from '@repo/shared-methods/utils';
const { request, uploadWithProgress } = createRequest();
// 公共的下拉字段
export const getSelectOptions = (name) => {
	const params = new URLSearchParams({ name });
	return request(`/select_options?${params}`, {
		method: 'GET',
	});
};
// lender下拉
export const getLenders = () => {
	return request('/loan/lender_companies', {
		method: 'GET',
	});
};
// borrower下拉
export const getBorrowers = () => {
	return request('/applicant/borrower', {
		method: 'POST',
	});
};

// 获取地址state列表
export const getStateList = () => {
	return request('/us_states', {
		method: 'GET',
	});
};

// 获取地址city列表通过state
export const getCityList = (state) => {
	const params = new URLSearchParams({ state });
	return request(`/us_cities?${params}`, {
		method: 'GET',
	});
};

// 获取上传文件到亚马逊的外部url
export const getUploadAwsUrl = (
	url = '/storage/presign_url',
	type = 'pdf',
	file_suffix = 'pdf',
) => {
	return request(url, {
		method: 'POST',
		body: {
			type,
			file_suffix,
		},
	});
};

// 上传文件到亚马逊
export const uploadAwsS3 = async ({ url, file, onProgress = (e) => e, type = null }) => {
	const file_suffix = file.name.split('.')?.pop?.()?.toLowerCase?.();
	const fileType = type || getFileTypeBySuffix(file_suffix);
	const { url: uploadUrl } = await getUploadAwsUrl(url, fileType, file_suffix);
	const res = await uploadWithProgress(uploadUrl, {
		method: 'PUT',
		body: file,
		onProgress,
		isExternalLink: true,
	});
	if (res) {
		return uploadUrl.split('?')[0];
	}
};
