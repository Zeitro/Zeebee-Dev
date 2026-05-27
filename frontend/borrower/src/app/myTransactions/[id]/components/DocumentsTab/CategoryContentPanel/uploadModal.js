import React, { useRef, useState } from 'react';
import { NewModal, UploadFileButton } from '@repo/shared-methods';
import { Button, message } from 'antd';
import { CloseOutlined, FileTextOutlined, InboxOutlined } from '@ant-design/icons';
import { createRequest } from '@repo/shared-methods/utils';
// const { TextArea } = Input;
const UploadModal = ({ transaction_id, open, setOpen, record = {}, onRefresh }) => {
	const { request } = createRequest();
	const uploadRef = useRef(null);
	const [loading, setLoading] = useState(false);
	const [fileList, setFileList] = useState([]);

	const handleOK = async () => {
		try {
			setLoading(true);
			const result = await uploadRef.current.triggerUpload(fileList);
			const files = result?.map((i) => ({
				file_name: i.file_name,
				storage_type: i.storage_type,
				upload_type: i.upload_type,
				file_url: i.url,
			}));
			await request(
				`/borrower_property_transaction/${transaction_id}/document/${record.id}/file`,
				{
					method: 'POST',
					body: {
						files,
					},
				},
			);
			message.open({
				type: 'success',
				content: 'Upload successfully',
			});
			onRefresh?.();
			setLoading(false);
			setOpen(false);
		} catch {
			setLoading(false);
		}
	};

	return (
		<NewModal
			title={'Upload document'}
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
			<div className="h-[400px] overflow-hidden overflow-y-auto">
				<div className="pb-4">
					<label>Document type</label>
					<div className="font-medium text-black-dark">{record.document_type}</div>
				</div>

				{/* <Form.Item label="Comment (Optional）" name="comment">
					<TextArea placeholder="Add a comment" autoSize={{ minRows: 3, maxRows: 5 }} />
				</Form.Item> */}
				<UploadFileButton
					ref={uploadRef}
					isReturnMoreInfo={true}
					immediateUpload={false}
					accept=".pdf,.png,.jpg,.jpeg"
					maxCount={10}
					multiple={true}
					onGetUploadList={setFileList}
					className={[
						'w-full [&_.ant-upload-list]:w-max [&_.ant-upload]:w-full [&_.ant-upload]:pb-6',
					].join(' ')}
					limitSize={10}
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
					<div className=" bg-gray-light py-4 w-full flex flex-col items-center border border-dashed border-base-border rounded-sm text-gray-light cursor-pointer hover:bg-[#E6F4FF]">
						<InboxOutlined style={{ fontSize: '48px', color: 'var(--base-primary)' }} />
						<div className="text-base pt-5 pb-1 text-foreground">
							Click or drag file to this area to upload
						</div>
					</div>
				</UploadFileButton>
			</div>
		</NewModal>
	);
};

export default UploadModal;
