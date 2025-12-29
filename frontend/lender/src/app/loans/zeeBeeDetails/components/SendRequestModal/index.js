import React, { use, useMemo, useState } from 'react';
import NewModal from '@/components/NewModal';
import { Button, Dropdown, Form, Input, message, Select } from 'antd';
import CommonUpload from '@/components/CommonUpload';
import { CloseOutlined, FileTextOutlined, InboxOutlined, UploadOutlined } from '@ant-design/icons';
const { TextArea } = Input;
import { createRequest } from '@/lib/fetch';
import eventBus from '@/lib/bus';
const SendRequestModal = ({ open, setOpen, maxCount, list = [] }) => {
	const { request } = createRequest();
	const [form] = Form.useForm();
	const requestList = Form.useWatch('documents', form);
	const [loading, setLoading] = useState(false);
	const remainderList = useMemo(() => {
		const arr = list
			.filter((item) => {
				if (!requestList?.find?.((i) => i.file_type === item.file_type)) {
					return item;
				}
			})
			.map((i) => ({
				key: i.id,
				label: i.file_type,
			}));
		return arr;
	}, [requestList, list]);
	const handleOK = async () => {
		try {
			const values = await form.validateFields();
			setLoading(true);
			await request('/loan_document/send_requests', {
				method: 'POST',
				body: {
					...values,
				},
			});
			message.open({
				type: 'success',
				content: 'Request sent successfully',
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<NewModal
			title={'Send request'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Send request"
			confirmLoading={loading}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				className="h-[400px] overflow-hidden overflow-y-auto"
				initialValues={{
					documents: list,
				}}
			>
				<Form.List name="documents">
					{(fields, { add, remove }) => {
						return (
							<>
								{fields?.map((field) => {
									const item = requestList?.[field.name] || {};
									return (
										<div
											key={field.key}
											className="not-last:border-b border-base-border mb-4"
										>
											<div className="flex items-start justify-between">
												<div className="pb-4">
													<label>Document type</label>
													<div className="font-medium text-black-dark">
														{item.file_type}
													</div>
												</div>
												{fields.length > 1 && (
													<Button
														type="primary"
														ghost
														size="middle"
														onClick={() => remove(field.name)}
													>
														Delete
													</Button>
												)}
											</div>
											<Form.Item
												label="Message to borrower （Optional）"
												name={[field.name, 'comment']}
											>
												<TextArea
													placeholder="Add message to borrower"
													autoSize={{ minRows: 3, maxRows: 5 }}
												/>
											</Form.Item>
										</div>
									);
								})}
								{(maxCount > 1 || !maxCount) && remainderList.length > 0 && (
									<Dropdown
										menu={{
											items: remainderList,
											onClick: ({ key }) => {
												const item = list.find((i) => i.id === Number(key));
												add({ id: Number(key), file_type: item.file_type });
											},
										}}
										placement="topLeft"
										arrow={{ pointAtCenter: true }}
									>
										<Button
											size="middle"
											onClick={() => add()}
											type="primary"
											ghost
										>
											Add more requests
										</Button>
									</Dropdown>
								)}
							</>
						);
					}}
				</Form.List>
			</Form>
		</NewModal>
	);
};

export default SendRequestModal;

export const SendRequestBtn = ({ children, maxCount, list, ...props }) => {
	const [open, setOpen] = useState(false);
	return (
		<div onClick={(e) => e.stopPropagation()}>
			<Button
				type="default"
				size="small"
				onClick={(e) => {
					e.stopPropagation();
					setOpen(true);
				}}
				{...props}
			>
				{children || 'Send request'}
			</Button>
			{open && (
				<SendRequestModal open={open} setOpen={setOpen} maxCount={maxCount} list={list} />
			)}
		</div>
	);
};
