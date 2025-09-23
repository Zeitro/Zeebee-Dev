import React, { useEffect, useState } from 'react';
import NewModal from '@/components/NewModal';
import { DownloadOutlined, InboxOutlined, InfoCircleOutlined } from '@ant-design/icons';
import CommonDownloadButton from '@/components/CommonDownload';
import CommonUpload from '@/components/CommonUpload';
import { createRequest } from '@/lib/fetch';
import { Modal } from 'antd';
import { initConfigFn } from '@/components/NewModal/BtnPopUpModal';
import eventBus from '@/lib/bus';
const ImportLoanModal = ({ setOpen }) => {
	const [confirmLoading, setConfirmLoading] = useState(false);
	const [fileList, setFileList] = useState([]);
	const [percent, setPercent] = useState(0);
	const [reUpload, setReUpload] = useState(null);
	const { uploadWithProgress } = createRequest();
	const [modal, contextHolder] = Modal.useModal();
	const handleUpload = async () => {
		try {
			setConfirmLoading(true);
			const formData = new FormData();
			formData.append('file', fileList[0]);
			formData.append('stage', 'Pending');
			const res = await uploadWithProgress('/loan/import', {
				method: 'POST',
				body: formData,
				onProgress: (v) => setPercent(v),
				isHideTips: true,
			});
			modal.success({
				...initConfigFn,
				title: 'Import successful',
				content: (
					<div className="text-sm pt-6">
						{res?.record_count} loan data have been created.
					</div>
				),
			});
			eventBus.emit('searchPendingLoans');
			setConfirmLoading(false);
			setOpen(false);
		} catch (error) {
			modal.confirm({
				...initConfigFn,
				title: 'Import failed',
				okText: 'Re-upload',
				onCancel: () => setOpen(false),
				onOk: () => {
					setReUpload(new Date().getTime());
					setFileList([]);
					setPercent(0);
				},
				content: (
					<>
						<span className="text-sm">
							Wrong data format! please modify the file and reupload.
						</span>
						{error.message && (
							<p className="text-sm p--6 text-base-error">{error.message}</p>
						)}
					</>
				),
			});
			setConfirmLoading(false);
		}
	};

	useEffect(() => {
		if (contextHolder) {
			eventBus.emit('appendPortalComponent', contextHolder);
		}
	}, [contextHolder]);
	return (
		<>
			<NewModal
				width={572}
				title="Import loans"
				okText={percent ? `${percent}% Importing...` : `Import loans`}
				confirmLoading={confirmLoading}
				cancelButtonProps={{ style: { display: 'none' } }}
				okButtonProps={{ disabled: !fileList?.length }}
				closable={!confirmLoading}
				onOk={handleUpload}
				onCancel={() => setOpen(false)}
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
					<CommonUpload
						immediateUpload={false}
						accept=".xlsx,.xls"
						onGetUploadList={setFileList}
						// disabled={fileList?.length}
						key={reUpload}
					>
						<div className=" bg-gray-light py-4 w-[395px] flex flex-col items-center border border-base-border rounded-sm cursor-pointer">
							<InboxOutlined
								style={{ fontSize: '48px', color: 'var(--base-primary)' }}
							/>
							<div className="text-base pt-5 pb-1 text-foreground">
								Click or drag file to this area to upload
							</div>
							<p className="text-base-label">Only .xlsx/.xls files are accepted</p>
						</div>
					</CommonUpload>
				</div>
			</NewModal>
		</>
	);
};

export default ImportLoanModal;
