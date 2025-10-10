import DocumentStatus from './DocumentStatus';
import MsgSvg from '@/assets/images/message.svg';
import { SendRequestBtn } from './SendRequestModal';
import { request } from 'http';
import { UploadModalBtn } from './AddDocument/UploadModal';
import TableOperation from '@/components/TableOperation';
const getUploadCommonColumns = (typeTitle, isHaveComments = true) => {
	const arr = [
		{
			title: typeTitle || 'Pending borrower upload',
			dataIndex: 'type',
			key: 'type',
			ellipsis: true,
			minWidth: 200,
		},
		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			ellipsis: true,
			width: 150,
			render: (v) => <DocumentStatus status={v} />,
		},

		{
			title: 'Upload Date',
			dataIndex: 'upload_date',
			key: 'upload_date',
			ellipsis: true,
			width: 150,
		},
	];
	if (isHaveComments) {
		arr.splice(2, 0, {
			title: 'Message',
			dataIndex: 'Message',
			key: 'Message',
			ellipsis: true,
			width: 150,
			render: (v) => <MsgSvg />,
		});
	}

	return arr;
};
export const borrowerUploadColumns = (isHide) => {
	const getOpt = () => {
		return {
			title: '',
			dataIndex: 'operation',
			width: 150,
			align: 'right',
			render: (_, record) => (
				<div className="flex items-center gap-2">
					<SendRequestBtn
						maxCount={1}
						list={[
							{
								document_type: record.type,
							},
						]}
					>
						Send request
					</SendRequestBtn>
					<TableOperation
						list={[
							{
								key: 'Upload',
								label: (
									<UploadModalBtn isChangeType={false} isButton={false}>
										Upload document
									</UploadModalBtn>
								),
							},
							{
								key: 'edit',
								label: 'Edit comments',
							},
							{
								key: 'delete',
								label: 'Delete request',
							},
						]}
						onSelect={(key) => {}}
					/>
				</div>
			),
		};
	};
	return [...getUploadCommonColumns('Pending borrower upload'), ...(isHide ? [] : [getOpt()])];
};

export const lenderUploadColumns = (isHide) => {
	// const status = {
	// 	title: 'Status',
	// 	dataIndex: 'pending_status',
	// 	key: 'pending_status',
	// 	ellipsis: true,
	// 	width: 150,
	// 	align: 'center',
	// 	fixed: 'right',
	// 	render: (v, record) => (
	// 		<LoanPendingStatus data={record} successCallback={successCallback} />
	// 	),
	// };
	const getOpt = () => {
		return {
			title: '',
			dataIndex: 'operation',
			width: 150,
			align: 'right',
			render: (_, record) => (
				<div className="flex items-center gap-2">
					<UploadModalBtn isChangeType={false}>Upload</UploadModalBtn>
					<TableOperation
						list={[
							{
								key: 'edit',
								label: 'Edit comments',
							},
							{
								key: 'delete',
								label: 'Delete request',
							},
						]}
						onSelect={(key) => {}}
					/>
				</div>
			),
		};
	};
	return [...getUploadCommonColumns('Pending internal upload'), ...(isHide ? [] : [getOpt()])];
};

export const lenderClosingColumns = () => {
	const getOpt = () => {
		return {
			title: '',
			dataIndex: 'operation',
			width: 120,
			align: 'right',
			fixed: 'right',
			render: (_, record) => (
				<div>
					<UploadModalBtn isChangeType={false}>Upload</UploadModalBtn>
				</div>
			),
		};
	};
	return [...getUploadCommonColumns('Pending internal upload', false), getOpt()];
};
