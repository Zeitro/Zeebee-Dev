import { Button, Divider } from 'antd';
import React, { useState } from 'react';
import { PlusOutlined, SendOutlined, UploadOutlined } from '@ant-design/icons';
import { NewModal } from '@repo/shared-methods';
import AddRequest from './AddRequest';
import UploadModal from './UploadModal';

export const AddDocumentModal = ({
	open,
	setOpen,
	setOpenUpload,
	setOpenAdd,
	isHideAddRequest,
}) => {
	return (
		<NewModal
			title={null}
			width={460}
			open={open}
			onCancel={() => setOpen(false)}
			footer={null}
		>
			<div className="py-[60px] px-12 flex flex-col justify-center">
				{!isHideAddRequest && (
					<>
						<div className="pb-4">
							Create a new document request for the borrowers or processors.
						</div>
						<Button
							type="primary"
							size="large"
							icon={<SendOutlined />}
							onClick={() => {
								setOpenAdd(true);
								setOpen(false);
							}}
						>
							Add Document request
						</Button>
						<Divider className="py-8" plain>
							or
						</Divider>
					</>
				)}
				<div className="pb-4">
					Directly upload a document that hasn't been added as a document request item.
				</div>
				<Button
					type="default"
					size="large"
					icon={<UploadOutlined />}
					onClick={() => {
						setOpenUpload(true);
						setOpen(false);
					}}
				>
					Upload new document
				</Button>
			</div>
		</NewModal>
	);
};

export const AddDocument = ({ loan_id, category, isHideAddRequest, ...props }) => {
	const [open, setOpen] = useState(false);
	const [openAdd, setOpenAdd] = useState(false);
	const [openUpload, setOpenUpload] = useState(false);
	return (
		<>
			<Button
				type="primary"
				icon={<PlusOutlined />}
				size="large"
				onClick={() => setOpen(true)}
				{...props}
			>
				Add Document
			</Button>
			<AddDocumentModal
				open={open}
				setOpen={setOpen}
				setOpenAdd={setOpenAdd}
				setOpenUpload={setOpenUpload}
				isHideAddRequest={isHideAddRequest}
			/>
			{openAdd && (
				<AddRequest
					open={openAdd}
					setOpen={setOpenAdd}
					loan_id={loan_id}
					category={category}
				/>
			)}
			{openUpload && (
				<UploadModal
					open={openUpload}
					setOpen={setOpenUpload}
					loan_id={loan_id}
					category={category}
				/>
			)}
		</>
	);
};
