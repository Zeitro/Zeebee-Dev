import { Form, message } from 'antd';
import { eventBus, createRequest } from '@repo/shared-methods/utils';
import { useState } from 'react';
const useData = ({ id, loan_id, apiType, category, setOpen, onSuccessCallback = () => {} }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	// 【 loan 】- 贷款内创建文档类型
	const handleCreateLoanDocumentType = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			await request('/loan_document', {
				method: 'POST',
				body: {
					...values,
					loan_id,
					category,
				},
			});
			message.open({
				type: 'success',
				content: 'Add Document request successfully',
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
		} catch {
			setLoading(false);
		}
	};

	// 【 transaction 】- 交易内创建文档类型
	const handleCreateTransactionDocumentType = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			await request(`/borrower_property_transaction/${id}/document`, {
				method: 'POST',
				body: {
					names: values.file_types || [],
					parent_id: category,
				},
			});
			message.open({
				type: 'success',
				content: 'Add Document type successfully',
			});
			onSuccessCallback?.();
			setLoading(false);
			setOpen(false);
		} catch {
			setLoading(false);
		}
	};

	// 根据apiType选择创建文档类型
	const switchCreateDocumentType = () => {
		switch (apiType) {
			case 'loan':
				return handleCreateLoanDocumentType();
			case 'transaction':
				return handleCreateTransactionDocumentType();
			default:
				return null;
		}
	};

	return { form, loading, handleCreateLoanDocumentType, switchCreateDocumentType };
};

export default useData;
