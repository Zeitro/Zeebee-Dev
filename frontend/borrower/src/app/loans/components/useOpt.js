'use client';
import React, { useState } from 'react';
import { createRequest } from '@/lib/fetch';
import { message, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { assemblyAddress } from '@/lib/format';
export default function useOpt(tableFilterRef) {
	const [modal, contextHolder] = Modal.useModal();
	const { request } = createRequest();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const handleDelete = async (id) => {
		try {
			setConfirmLoading(true);
			await request(`/loan/${id}`, { method: 'DELETE' });
			tableFilterRef?.current?.handleSearch();
			setConfirmLoading(false);
			message.open({
				type: 'success',
				content: 'Delete successful!',
			});
		} catch (error) {
			setConfirmLoading(false);
		}
	};

	const showDeleteModal = ({ property, id }) => {
		modal.confirm({
			title: 'Are you sure you want to delete this loan?',
			icon: <ExclamationCircleOutlined className="!text-base-warning" />,
			content: (
				<div>
					You are about to delete the loan associated with:
					<p className="font-bold pb-8">{assemblyAddress(property?.address, true)}</p>
					Once deleted , this loan will no longer appear in your portfolio and cannot be
					recovered.
				</div>
			),
			onOk: () => handleDelete(id),
			okText: 'Delete loan',
			okButtonProps: {
				danger: true,
			},
			confirmLoading,
			centered: true,
		});
	};
	return { showDeleteModal, contextHolder };
}
