import { createRequest } from '@repo/shared-methods/utils';
import { message } from 'antd';
import { useState } from 'react';
export const useBorrower = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_applicant/${id}`, {
				method: 'DELETE',
			});
			handleSearch();
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

export const useBorrowingEntity = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/borrower_entity/${id}`, {
				method: 'DELETE',
			});
			handleSearch();
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

export const useProperty = ({ handleSearch }) => {
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setLoading(true);
			await request(`/property/${id}`, { method: 'DELETE' });
			handleSearch();
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
