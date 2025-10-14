import DocumentStatus from './DocumentStatus';
import { UploadModalBtn } from './UploadModal';
import { formatDate } from '@/lib/date';
import { Tooltip } from 'antd';
const getUploadCommonColumns = (typeTitle, isHaveComments = true) => {
	const getStatus = (v) => {
		switch (v) {
			case 'Pending Upload':
				return '--';
			case 'Requested':
				return 'Waiting for upload';
			default:
				return v;
		}
	};
	const arr = [
		{
			title: typeTitle || 'Requested document',
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
			width: 150,
			render: (v) => <DocumentStatus status={v} alias={getStatus(v)} />,
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
					<span className="font-normal text-foreground">{v}</span>
				</Tooltip>
			),
		});
	}

	return arr;
};

export const documentsColumns = () => {
	const getOpt = () => {
		return {
			title: '',
			dataIndex: 'operation',
			width: 120,
			align: 'right',
			fixed: 'right',
			render: (_, record) => (
				<UploadModalBtn
					isChangeType={false}
					record={record}
					disabled={record.status === 'Pending upload'}
				>
					Upload document
				</UploadModalBtn>
			),
		};
	};
	return [...getUploadCommonColumns('Requested document', true), getOpt()];
};
