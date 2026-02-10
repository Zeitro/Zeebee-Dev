import DocumentStatus from './DocumentStatus';
import { SendRequestBtn } from './SendRequestModal';
import { UploadModalBtn } from './AddDocument/UploadModal';
import { TableOperation, DeleteBtn } from '@repo/shared-methods';
import { formatDate } from '@repo/shared-methods/utils';
import { Tooltip } from 'antd';
// import { DeleteBtn } from '@/components/NewModal/BtnPopUpModal';
import { EditCommentModalBtn } from './EditCommentModal';
const getUploadCommonColumns = (typeTitle, isHaveComments = true, statusWidth) => {
	const arr = [
		{
			title: typeTitle || 'Pending borrower upload',
			dataIndex: 'file_type',
			key: 'file_type',
			ellipsis: true,
			minWidth: 200,
		},
		{
			title: 'Status',
			dataIndex: 'status',
			key: 'status',
			ellipsis: true,
			width: statusWidth || (isHaveComments ? 150 : 366),
			render: (v) => <DocumentStatus status={v} />,
		},

		{
			title: 'Upload Date',
			dataIndex: 'last_uploaded_at',
			key: 'last_uploaded_at',
			ellipsis: true,
			width: 100,
			render: (v) => formatDate(v),
		},
	];
	if (isHaveComments) {
		arr.splice(2, 0, {
			title: 'Comment',
			dataIndex: 'comment',
			key: 'comment',
			ellipsis: true,
			width: 200,
			render: (v) => (
				<Tooltip title={v}>
					<span className="font-normal text-foreground">{v || '--'}</span>
				</Tooltip>
			),
		});
	}

	return arr;
};
export const borrowerUploadColumns = ({
	isHide,
	handleClickOpt,
	deleteLoading,
	disabledOperate,
}) => {
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
								file_type: record.file_type,
								id: record.id,
							},
						]}
						disabled={disabledOperate || record.status !== 'Pending request'}
					>
						Send request
					</SendRequestBtn>
					<TableOperation
						sliceCount={0}
						disabled={disabledOperate}
						list={[
							{
								key: 'Upload',
								label: (
									<UploadModalBtn
										isChangeType={false}
										isButton={false}
										record={record}
									>
										Upload document
									</UploadModalBtn>
								),
							},
							{
								key: 'edit',
								label: (
									<EditCommentModalBtn record={record}>
										Edit comment
									</EditCommentModalBtn>
								),
							},
							{
								key: 'delete',
								label: (
									<DeleteBtn
										isButton={false}
										content={record?.file_type}
										onOk={() => handleClickOpt('delete', record)}
										confirmLoading={deleteLoading}
									>
										Delete request
									</DeleteBtn>
								),
							},
						]}
						onSelect={handleClickOpt}
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
				<UploadModalBtn isChangeType={false} record={record} isStopPropagation={true}>
					Upload
				</UploadModalBtn>
			),
		};
	};
	return [
		...getUploadCommonColumns('Pending internal upload', false),
		...(isHide ? [] : [getOpt()]),
	];
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
				<UploadModalBtn isChangeType={false} record={record} isStopPropagation={true}>
					Upload
				</UploadModalBtn>
			),
		};
	};
	return [...getUploadCommonColumns('Pending internal upload', false, 150), getOpt()];
};
