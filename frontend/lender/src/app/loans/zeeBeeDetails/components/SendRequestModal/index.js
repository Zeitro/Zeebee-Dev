import React, { use, useMemo, useState } from 'react';
import NewModal from '@/components/NewModal';
import { Button, Dropdown, Form, Input, Select } from 'antd';
import CommonUpload from '@/components/CommonUpload';
import { CloseOutlined, FileTextOutlined, InboxOutlined, UploadOutlined } from '@ant-design/icons';
const { TextArea } = Input;
const SendRequestModal = ({ open, setOpen, maxCount, list = [{ document_type: 111 }] }) => {
	const [form] = Form.useForm();
	const requestList = Form.useWatch('request', form);
	const remainderList = useMemo(() => {
		const arr = list
			.filter((item) => {
				console.log(requestList);

				if (!requestList?.find?.((i) => i.document_type === item.document_type)) {
					return item;
				}
			})
			.map((i) => ({
				key: i.document_type,
				label: i.document_type,
			}));
		return arr;
	}, [requestList, list]);
	const handleOK = async () => {
		const values = await form.validateFields();
	};
	return (
		<NewModal
			title={'Add document request'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Send request"
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				className="h-[400px] overflow-hidden overflow-y-auto"
				initialValues={{
					request: [...list],
				}}
			>
				<Form.List name="request">
					{(fields, { add, remove }) => {
						return (
							<>
								{fields?.map((field, index) => (
									<div
										key={field.key}
										className="not-last:border-b border-base-border mb-4"
									>
										<div className="flex items-start justify-between">
											<div className="pb-4">
												<label>Document type</label>
												<div className="font-medium text-black-dark">
													1111
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
											name="message"
										>
											<TextArea
												placeholder="Add message to borrower"
												autoSize={{ minRows: 3, maxRows: 5 }}
											/>
										</Form.Item>
									</div>
								))}
								{(maxCount > 1 || !maxCount) && remainderList.length > 0 && (
									<Dropdown
										menu={{
											items: remainderList,
											onSelect: ({ key }) =>
												add({
													document_type: key,
												}),
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
		<>
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
		</>
	);
};
