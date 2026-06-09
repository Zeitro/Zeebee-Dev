import React, { useRef, useState } from 'react';
import { Button, Form, Input, message } from 'antd';
import { NewModal, UploadFileButton } from '@repo/shared-methods';
import { CloseOutlined, FileTextOutlined, InboxOutlined } from '@ant-design/icons';
import styles from './index.module.scss';
import { createRequest, eventBus } from '@repo/shared-methods/utils';

const UploadModal = ({ open, setOpen, isChangeType = true, record = {}, loan_id }) => {
	const accept = '.pdf,.png,.jpg,.jpeg,.xlsx,.xls,.docx,.doc,.pptx,.ppt';
	const [form] = Form.useForm();
	const { request } = createRequest();
	const uploadRef = useRef(null);
	const [loading, setLoading] = useState(false);
	const [fileList, setFileList] = useState([]);
	const handleOK = async () => {
		try {
			const values = await form.validateFields();
			let id = record?.id;
			setLoading(true);

			if (isChangeType) {
				// 创建文档类型
				const { loan_documents } = await request('/loan_document', {
					method: 'POST',
					body: {
						file_types: [values.file_type],
						loan_id,
					},
				});

				id = loan_documents[0].id;
			}
			const file_urls = await uploadRef.current.triggerUpload(fileList);
			await request(`/loan_document/${id}/files`, {
				method: 'POST',
				body: {
					file_urls,
				},
			});
			message.open({
				type: 'success',
				content: 'Upload successfully',
			});
			eventBus.emit('getLoanDetail');
			setLoading(false);
			setOpen(false);
		} catch (error) {
			console.log(error);

			setLoading(false);
		}
	};

	return (
		<NewModal
			title={isChangeType ? 'Upload new document' : 'Upload document'}
			width={650}
			open={open}
			onCancel={() => setOpen(false)}
			onOk={handleOK}
			okText="Save"
			confirmLoading={loading}
			okButtonProps={{
				disabled: !fileList.length,
			}}
		>
			<Form
				form={form}
				layout="vertical"
				autoComplete="off"
				size="large"
				className="h-[500px] overflow-hidden overflow-y-auto"
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
				<UploadFileButton
					ref={uploadRef}
					immediateUpload={false}
					accept={accept}
					maxCount={10}
					multiple={true}
					onGetUploadList={setFileList}
					className={[styles.documentUpload].join(' ')}
					disabled={false}
					itemRender={(node, file, fileList, { remove }) => {
						return (
							<div className="flex items-center px-4 py-[5px] mb-6 gap-2 border border-base-border rounded-md w-max">
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
					<div className=" bg-gray-light py-20 px-10 w-full flex flex-col items-center border border-dashed border-base-border rounded-sm text-gray-light cursor-pointer hover:bg-[#E6F4FF]">
						<InboxOutlined style={{ fontSize: '48px', color: 'var(--base-primary)' }} />
						<div className="text-base pt-5 pb-1 text-foreground">
							Click or drag file to this area to upload
						</div>
						<div className="text-base-label text-sm pt-4">
							<span className="font-medium">Upload requirements</span>
							<br />
							1. Accepted formats: PDF, PNG, JPG, JPEG, Excel (.xlsx, .xls), Word
							(.docx, .doc), and PowerPoint (.pptx, .ppt).
							<br />
							2. Maximum file size: 20 MB.
						</div>
					</div>
				</UploadFileButton>
			</Form>
		</NewModal>
	);
};

export default UploadModal;
export const UploadModalBtn = ({
	children,
	isChangeType,
	isButton = true,
	record,
	isStopPropagation,
	...props
}) => {
	const [open, setOpen] = useState(false);
	return (
		<div onClick={(e) => isStopPropagation && e.stopPropagation()}>
			{isButton ? (
				<Button
					type="default"
					size="small"
					onClick={() => {
						setOpen(true);
					}}
					{...props}
				>
					{children || 'Send request'}
				</Button>
			) : (
				<div
					onClick={() => {
						setOpen(true);
					}}
				>
					{children || 'Send request'}
				</div>
			)}
			{open && (
				<UploadModal
					open={open}
					setOpen={setOpen}
					isChangeType={isChangeType}
					record={record}
				/>
			)}
		</div>
	);
};
