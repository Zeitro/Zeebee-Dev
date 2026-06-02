'use client';
import { createRequest } from '@repo/shared-methods/utils';
import { message } from 'antd';
import { useState } from 'react';
import _ from 'lodash';
export const useBorrower = ({ handleSearch, handleClickItem, rowData }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_applicant/${id}`, {
				method: 'DELETE',
			});
			const list = await handleSearch();
			const newList = _.cloneDeep(list);
			if (newList?.length && rowData?.id === id) {
				handleClickItem?.(list[0], newList);
			}
			message.destroy();
			message.open({
				type: 'success',
				content: 'Borrower deleted successfully!',
			});
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	return { handleDelete, loading };
};

export const useBorrowingEntity = ({ handleSearch, handleClickItem, rowData }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_entity/${id}`, {
				method: 'DELETE',
			});
			const list = await handleSearch();
			const newList = _.cloneDeep(list);
			if (newList?.length && rowData?.id === id) {
				handleClickItem?.(list[0], newList);
			}
			message.destroy();
			message.open({
				type: 'success',
				content: 'Borrowing entity deleted successfully!',
			});
			setLoading(false);
		} catch (error) {
			console.log(error);
			setLoading(false);
		}
	};

	return { handleDelete, loading };
};

export const useProperty = ({ handleSearch, handleClickItem, rowData }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/property/${id}`, { method: 'DELETE' });
			const list = await handleSearch();
			const newList = _.cloneDeep(list);
			if (newList?.length && rowData?.id === id) {
				handleClickItem?.(list[0], newList);
			}
			message.destroy();
			message.open({
				type: 'success',
				content: 'Property deleted successfully!',
			});
			setLoading(false);
		} catch {
			setLoading(false);
		}
	};

	return { handleDelete, loading };
};

// export const useTransaction = ({ handleSearch, handleClickItem, rowData }) => {
// 	const { request } = createRequest();
// 	const [loading, setLoading] = useState(false);
// 	const handleDelete = async (id) => {
// 		try {
// 			setLoading(true);
// 			await request(`/transaction/${id}`, { method: 'DELETE' });
// 			const list = await handleSearch();
// 			const newList = _.cloneDeep(list);
// 			if (newList?.length && rowData?.id === id) {
// 				handleClickItem?.(list[0], newList);
// 			}
// 			message.destroy();
// 			message.open({
// 				type: 'success',
// 				content: 'Transaction deleted successfully!',
// 			});
// 			setLoading(false);
// 		} catch {
// 			setLoading(false);
// 		}
// 	};

// 	return { handleDelete, loading };
// };
