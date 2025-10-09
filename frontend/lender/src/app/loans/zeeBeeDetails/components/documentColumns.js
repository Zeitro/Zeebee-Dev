import DocumentStatus from './DocumentStatus';
import MsgSvg from '@/assets/images/message.svg';
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
			title: 'Comments',
			dataIndex: 'comments',
			key: 'comments',
			ellipsis: true,
			width: 150,
			render: (v) => <MsgSvg />,
		});
	}
	return arr;
};
export const borrowerUploadColumns = () => {
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
			width: 120,
			align: 'left',
			fixed: 'right',
			render: (_, record) => 'send request',
		};
	};
	return [...getUploadCommonColumns('Pending borrower upload'), getOpt()];
};

export const lenderUploadColumns = () => {
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
			width: 120,
			align: 'left',
			fixed: 'right',
			render: (_, record) => 'upload',
		};
	};
	return [...getUploadCommonColumns('Pending internal upload'), getOpt()];
};

export const lenderClosingColumns = () => {
	const getOpt = () => {
		return {
			title: '',
			dataIndex: 'operation',
			width: 120,
			align: 'left',
			fixed: 'right',
			render: (_, record) => 'upload documents',
		};
	};
	return [...getUploadCommonColumns('Pending internal upload', false), getOpt()];
};
