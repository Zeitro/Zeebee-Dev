import { getFileTypeBySuffix } from '../lib/common';
import { createRequest } from '../lib/fetch';

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
	file_name = null,
) => {
	return request(url, {
		method: 'POST',
		body: {
			type,
			file_suffix,
			file_name,
		},
	});
};

// 上传文件到亚马逊
// file 为blob文件，本地上传的文件则传file.originFileObj
export const uploadAwsS3 = async ({
	url,
	file,
	onProgress = (e) => {},
	type = null,
	isReturnMoreInfo = false, // true返回更多信息，如文件url，文件大小，文件类型等，false只返回文件url
}) => {
	const file_suffix = file.name.split('.')?.pop?.()?.toLowerCase?.();
	const fileType = type || getFileTypeBySuffix(file_suffix);
	const { url: uploadUrl, ...moreInfo } = await getUploadAwsUrl(
		url,
		fileType,
		file_suffix,
		file?.name,
	);
	const res = await uploadWithProgress(uploadUrl, {
		method: 'PUT',
		body: file,
		onProgress,
		isExternalLink: true,
	});
	if (!res) return null;
	if (isReturnMoreInfo) {
		return {
			url: uploadUrl.split('?')[0],
			...moreInfo,
		};
	}
	return uploadUrl.split('?')[0];
};
