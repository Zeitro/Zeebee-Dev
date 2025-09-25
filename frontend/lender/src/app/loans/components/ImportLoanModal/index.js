import React, { useState } from 'react';
import NewModal from '@/components/NewModal';
import { DownloadOutlined, InboxOutlined, InfoCircleOutlined } from '@ant-design/icons';
import CommonDownloadButton from '@/components/CommonDownload';
import CommonUpload from '@/components/CommonUpload';
const ImportLoanModal = ({ setOpen }) => {
	const [confirmLoading, setConfirmLoading] = useState(false);

	const handleUpload = async () => {
		try {
			setConfirmLoading(true);
			setConfirmLoading(false);
			setOpen(false);
		} catch (error) {
			setConfirmLoading(false);
		}
	};
	return (
		<NewModal
			width={572}
			title="Import loans"
			onCancel={() => setOpen(false)}
			okText={`Import loans`}
			onOk={handleUpload}
			confirmLoading={confirmLoading}
		>
			<div className="flex flex-col items-center gap-6">
				<div className="pt-[26px] pb-[10px]">
					<InfoCircleOutlined
						style={{ fontSize: '16px', color: 'var(--base-primary)' }}
						className="text-base-primary"
					/>
					<span className="ml-2">
						Please use the template provided to ensure correct formatting.
					</span>
				</div>
				<CommonDownloadButton
					icon={<DownloadOutlined style={{ fontSize: '16px' }} />}
					className="w-max"
					fileUrl="/3rdPartyTemplate.xlsx"
				>
					Download template
				</CommonDownloadButton>
				<CommonUpload>
					<div className=" bg-gray-light py-4 mb-4 w-[395px] flex flex-col items-center border border-base-border rounded-sm text-gray-light">
						<InboxOutlined style={{ fontSize: '48px', color: 'var(--base-primary)' }} />
						<div className="text-base pt-5 pb-1 text-foreground">
							Click or drag file to this area to upload
						</div>
						<p className="text-base-label">Only .csv files are accepted</p>
					</div>
				</CommonUpload>
			</div>
		</NewModal>
	);
};

export default ImportLoanModal;
