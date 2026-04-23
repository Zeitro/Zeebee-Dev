'use client';
import React, { useState } from 'react';
import { assemblyAddress, createRequest } from '@repo/shared-methods/utils';
import { message, Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { useGlobalStore } from '@/store/global';

export default function useOpt(tableFilterRef) {
	const [modal, contextHolder] = Modal.useModal();
	const [confirmLoading, setConfirmLoading] = useState(false);
	const { request } = createRequest();
	const router = useRouter();
	const { changeSidebarCollapsed } = useGlobalStore();

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
		} catch {
			setConfirmLoading(false);
		}
	};

	const showDeleteModal = ({ loan_property, id }) => {
		modal.confirm({
			title: 'Are you sure you want to delete this loan?',
			icon: <ExclamationCircleOutlined className="!text-base-warning" />,
			content: (
				<div>
					You are about to delete the loan associated with:
					<p className="font-bold pb-8">
						{assemblyAddress(loan_property?.address, true)}
					</p>
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

	const handleSelectedRow = (record, stage) => {
		let url;
		if (stage === 'Draft') {
			url =
				record?.source === 'Borrower'
					? `/loans/${record?.id}`
					: `/loans/createZeeBeeLoan?id=${record?.id}`;
		} else if (record?.source === 'Borrower') {
			url = `/loans/${record?.id}`;
		} else {
			url = `/loans/zeeBeeDetails?id=${record?.id}`;
		}
		router.push(url);
		setTimeout(() => changeSidebarCollapsed(true), 300);
	};
	return { showDeleteModal, contextHolder, handleSelectedRow };
}
