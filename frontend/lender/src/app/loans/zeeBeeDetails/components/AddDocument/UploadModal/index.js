import React, { use, useState } from 'react';
import NewModal from '@/components/NewModal';
import { Button, Form, Input, message, Select } from 'antd';
import CommonUpload from '@/components/CommonUpload';
import { CloseOutlined, FileTextOutlined, InboxOutlined, UploadOutlined } from '@ant-design/icons';
import styles from './index.module.scss';
import { createRequest } from '@/lib/fetch';
import eventBus from '@/lib/bus';
const { TextArea } = Input;
const UploadModal = ({ open, setOpen, isChangeType = true, record = {} }) => {
	const [form] = Form.useForm();
	const { request } = createRequest();
	const [loading, setLoading] = useState(false);
	const [fileList, setFileList] = useState([]);
	const handleOK = async () => {
		try {
			const values = await form.validateFields();
			console.log(values, fileList);

			// setLoading(true);
			// await request('/loan_document/send_requests', {
			// 	method: 'POST',
			// 	body: {
			// 		...values,
			// 	},
			// });
			// message.open({
			// 	type: 'success',
			// 	content: 'Request sent successfully',
			// });
			// eventBus.emit('getLoanDetail');
			// setLoading(false);
			setOpen(false);
		} catch (error) {
			setLoading(false);
		}
	};
	return (
		<NewModal
			title={'Upload new document'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Save"
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				className="h-[400px] overflow-hidden overflow-y-auto"
				initialValues={{
					id: record.id,
				}}
			>
				{isChangeType ? (
					<Form.Item
						label="Document type(Required)"
						name="file_type"
						rules={[
							{
								required: true,
								message: 'Please enter document type!',
							},
						]}
					>
						<Input placeholder="Enter document type" style={{ width: 300 }} />
					</Form.Item>
				) : (
					<div className="pb-4">
						<label>Document type</label>
						<div className="font-medium text-black-dark">{record.file_type}</div>
					</div>
				)}
				{/* <Form.Item label="Comment (Optional）" name="comment">
					<TextArea placeholder="Add a comment" autoSize={{ minRows: 3, maxRows: 5 }} />
				</Form.Item> */}
				<CommonUpload
					immediateUpload={false}
					accept=".pdf,.png,.jpg,.jpeg"
					maxCount={10}
					multiple={true}
					onGetUploadList={setFileList}
					className={[styles.documentUpload].join(' ')}
					limitSize={10}
					disabled={false}
					itemRender={(node, file, fileList, { remove }) => {
						return (
							<div className="flex items-center px-4 py-[5px] gap-2 border border-base-border rounded-md">
								<FileTextOutlined />
								<div>{file.name}</div>
								<CloseOutlined
									onClick={remove}
									style={{ color: 'var(--base-label)' }}
								/>
							</div>
						);
					}}
				>
					<div className=" bg-gray-light py-4 w-full flex flex-col items-center border border-dashed border-base-border rounded-sm text-gray-light cursor-pointer hover:bg-[#E6F4FF]">
						<InboxOutlined style={{ fontSize: '48px', color: 'var(--base-primary)' }} />
						<div className="text-base pt-5 pb-1 text-foreground">
							Click or drag file to this area to upload
						</div>
					</div>
				</CommonUpload>
			</Form>
		</NewModal>
	);
};

export default UploadModal;
export const UploadModalBtn = ({ children, isChangeType, isButton = true, record, ...props }) => {
	const [open, setOpen] = useState(false);
	return (
		<>
			{isButton ? (
				<Button
					type="default"
					size="small"
					onClick={(e) => {
						e.stopPropagation();
						e.preventDefault();
						setOpen(true);
					}}
					{...props}
				>
					{children || 'Send request'}
				</Button>
			) : (
				<div>{children || 'Send request'}</div>
			)}
			{open && (
				<UploadModal
					open={open}
					setOpen={setOpen}
					isChangeType={isChangeType}
					record={record}
				/>
			)}
		</>
	);
};
